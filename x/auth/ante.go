package auth

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"errors"
	"runtime/debug"
)

const (
	deductFeesCost    sdk.Gas = 10
	memoCostPerByte   sdk.Gas = 1
	verifyCost                = 100
	maxMemoCharacters         = 100
)

// NewAnteHandler returns an AnteHandler that checks
// and increments sequence numbers, checks signatures & account numbers,
// and deducts fees from the first signer.
func NewAnteHandler(am AccountMapper, fck FeeCollectionKeeper) sdk.AnteHandler {

	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {

		// This AnteHandler requires Txs to be StdTxs
		stdTx, ok := tx.(StdTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		// set the gas meter
		newCtx = ctx.WithGasMeter(sdk.NewGasMeter(stdTx.Fee.Gas))

		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
					res = sdk.ErrOutOfGas(log).Result()
					res.GasWanted = stdTx.Fee.Gas
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()

		err := validateBasic(stdTx)
		if err != nil {
			return newCtx, err.Result(), true
		}

		sigs := stdTx.GetSignatures()
		signerAddrs := stdTx.GetSigners()
		msgs := tx.GetMsgs()

		// charge gas for the memo
		newCtx.GasMeter().ConsumeGas(memoCostPerByte*sdk.Gas(len(stdTx.GetMemo())), "memo")

		// Get the sign bytes (requires all account & sequence numbers and the fee)
		sequences := make([]int64, len(sigs))
		accNums := make([]int64, len(sigs))
		for i := 0; i < len(sigs); i++ {
			sequences[i] = sigs[i].Sequence
			accNums[i] = sigs[i].AccountNumber
		}

		fee := StdFee{
			Gas: stdTx.Fee.Gas,
			Amount: sdk.Coins{fck.GetNativeFeeToken(ctx, stdTx.Fee.Amount)},
		}

		err = fck.FeePreprocess(newCtx, fee.Amount, fee.Gas)
		if err != nil {
			return newCtx, err.Result(), true
		}

		// Check sig and nonce and collect signer accounts.
		var signerAccs = make([]Account, len(signerAddrs))
		var firstAccount Account
		for i := 0; i < len(sigs); i++ {
			signerAddr, sig := signerAddrs[i], sigs[i]

			// check signature, return account with incremented nonce
			signBytes := StdSignBytes(newCtx.ChainID(), accNums[i], sequences[i], fee, msgs, stdTx.GetMemo())
			signerAcc, res := processSig(
				newCtx, am,
				signerAddr, sig, signBytes,
			)
			if !res.IsOK() {
				return newCtx, res, true
			}

			if i == 0 {
				firstAccount = signerAcc
			} else {
				am.SetAccount(newCtx, signerAcc)
			}

			signerAccs[i] = signerAcc
		}

		if firstAccount != nil {
			newCtx.GasMeter().ConsumeGas(deductFeesCost, "deductFees")
			firstAccount, res = deductFees(firstAccount, fee)
			if !res.IsOK() {
				return newCtx, res, true
			}
			// Save the account.
			am.SetAccount(newCtx, firstAccount)
			fck.addCollectedFees(newCtx, fee.Amount)
		}

		// cache the signer accounts in the context
		newCtx = WithSigners(newCtx, signerAccs)

		res.GasWanted = stdTx.Fee.Gas
		//GetNativeFeeToken and FeePreprocess will ensure that there must be just one fee token
		res.FeeAmount = fee.Amount[0].Amount.Int64()
		res.FeeDenom = fee.Amount[0].Denom

		return newCtx, res, false // continue...
	}
}

func NewFeeRefundHandler(am AccountMapper, fck FeeCollectionKeeper) sdk.FeeRefundHandler {
	return func(ctx sdk.Context, tx sdk.Tx, txResult sdk.Result) (result sdk.Result, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("encountered panic error during fee refund, recovered: %v\nstack:\n%v", r, string(debug.Stack())))
			}
		}()

		txAccounts := GetSigners(ctx)
		// If this tx failed in anteHandler, txAccount length will be less than 1
		if len(txAccounts) < 1 {
			result.FeeAmount = txResult.FeeAmount
			return result, nil
		}
		firstAccount := txAccounts[0]
		//If all gas has been consumed, then there is no necessary to run fee refund process
		if txResult.GasWanted <= txResult.GasUsed {
			result.FeeAmount = txResult.FeeAmount
			return result, nil
		}

		stdTx, ok := tx.(StdTx)
		if !ok {
			return sdk.Result{}, errors.New("transaction is not Stdtx")
		}
		// Refund process will also cost gas, but this is compensation for previous fee deduction.
		// It is not reasonable to consume users' gas. So the context gas is reset to transaction gas
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(stdTx.Fee.Gas))

		fee := StdFee{
			Gas: stdTx.Fee.Gas,
			Amount: sdk.Coins{fck.GetNativeFeeToken(ctx, stdTx.Fee.Amount)}, // consume gas
		}

		unusedGas := txResult.GasWanted - txResult.GasUsed
		var refundCoins sdk.Coins
		for _,coin := range fee.Amount {
			newCoin := sdk.Coin{
				Denom:	coin.Denom,
				Amount: coin.Amount.Mul(sdk.NewInt(unusedGas)).Div(sdk.NewInt(txResult.GasWanted)),
			}
			refundCoins = append(refundCoins, newCoin)
		}
		coins := am.GetAccount(ctx, firstAccount.GetAddress()).GetCoins()   // consume gas
		err = firstAccount.SetCoins(coins.Plus(refundCoins))
		if err != nil {
			return sdk.Result{}, err
		}

		am.SetAccount(ctx, firstAccount)                                    // consume gas
		fck.refundCollectedFees(ctx, refundCoins)                           // consume gas
		// There must be just one fee token
		result.FeeAmount = fee.Amount[0].Amount.Mul(sdk.NewInt(txResult.GasUsed)).Div(sdk.NewInt(txResult.GasWanted)).Int64()

		return
	}
}

// Validate the transaction based on things that don't depend on the context
func validateBasic(tx StdTx) (err sdk.Error) {
	// Assert that there are signatures.
	sigs := tx.GetSignatures()
	if len(sigs) == 0 {
		return sdk.ErrUnauthorized("no signers")
	}

	// Assert that number of signatures is correct.
	var signerAddrs = tx.GetSigners()
	if len(sigs) != len(signerAddrs) {
		return sdk.ErrUnauthorized("wrong number of signers")
	}

	memo := tx.GetMemo()
	if len(memo) > maxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf("maximum number of characters is %d but received %d characters",
				maxMemoCharacters, len(memo)))
	}
	return nil
}

// verify the signature and increment the sequence.
// if the account doesn't have a pubkey, set it.
func processSig(
	ctx sdk.Context, am AccountMapper,
	addr sdk.AccAddress, sig StdSignature, signBytes []byte) (
	acc Account, res sdk.Result) {

	// Get the account.
	acc = am.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(addr.String()).Result()
	}

	// Check account number.
	accnum := acc.GetAccountNumber()
	if accnum != sig.AccountNumber {
		return nil, sdk.ErrInvalidSequence(
			fmt.Sprintf("Invalid account number. Got %d, expected %d", sig.AccountNumber, accnum)).Result()
	}

	// Check and increment sequence number.
	seq := acc.GetSequence()
	if seq != sig.Sequence {
		return nil, sdk.ErrInvalidSequence(
			fmt.Sprintf("Invalid sequence. Got %d, expected %d", sig.Sequence, seq)).Result()
	}
	err := acc.SetSequence(seq + 1)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	// If pubkey is not known for account,
	// set it from the StdSignature.
	pubKey := acc.GetPubKey()
	if pubKey == nil {
		pubKey = sig.PubKey
		if pubKey == nil {
			return nil, sdk.ErrInvalidPubKey("PubKey not found").Result()
		}
		if !bytes.Equal(pubKey.Address(), addr) {
			return nil, sdk.ErrInvalidPubKey(
				fmt.Sprintf("PubKey does not match Signer address %v", addr)).Result()
		}
		err = acc.SetPubKey(pubKey)
		if err != nil {
			return nil, sdk.ErrInternal("setting PubKey on signer's account").Result()
		}
	}

	// Check sig.
	ctx.GasMeter().ConsumeGas(verifyCost, "ante verify")
	if !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return
}

// Deduct the fee from the account.
// We could use the CoinKeeper (in addition to the AccountMapper,
// because the CoinKeeper doesn't give us accounts), but it seems easier to do this.
func deductFees(acc Account, fee StdFee) (Account, sdk.Result) {
	coins := acc.GetCoins()
	feeAmount := fee.Amount

	newCoins := coins.Minus(feeAmount)
	if !newCoins.IsNotNegative() {
		errMsg := fmt.Sprintf("%s < %s", coins, feeAmount)
		return nil, sdk.ErrInsufficientFunds(errMsg).Result()
	}
	err := acc.SetCoins(newCoins)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc, sdk.Result{}
}

// BurnFeeHandler burns all fees (decreasing total supply)
func BurnFeeHandler(_ sdk.Context, _ sdk.Tx, _ sdk.Coins) {}
