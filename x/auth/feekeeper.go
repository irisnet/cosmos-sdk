package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"fmt"
	"errors"
)

var (
	collectedFeesKey = []byte("collectedFees")
	NativeFeeTokenKey = "feeToken/native"
	NativeGasPriceThresholdKey  = "feeToken/native/gasPrice/threshold"
	FeeExchangeRatePrefix = "feeToken/derived/exchange/rate/"	//  key = feeToken/derived/exchange/rate/<denomination>, rate = BigInt(value)/10^18
	RatePrecision = int64(1000000000) //10^9
)

// This FeeCollectionKeeper handles collection of fees in the anteHandler
// and setting of MinFees for different fee tokens
type FeeCollectionKeeper struct {

	getter params.Getter

	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewFeeKeeper returns a new FeeKeeper
func NewFeeCollectionKeeper(cdc *wire.Codec, key sdk.StoreKey, getter params.Getter) FeeCollectionKeeper {
	return FeeCollectionKeeper{
		key: key,
		cdc: cdc,
		getter: getter,
	}
}

// Adds to Collected Fee Pool
func (fck FeeCollectionKeeper) GetCollectedFees(ctx sdk.Context) sdk.Coins {
	store := ctx.KVStore(fck.key)
	bz := store.Get(collectedFeesKey)
	if bz == nil {
		return sdk.Coins{}
	}

	feePool := &(sdk.Coins{})
	fck.cdc.MustUnmarshalBinary(bz, feePool)
	return *feePool
}

// Sets to Collected Fee Pool
func (fck FeeCollectionKeeper) setCollectedFees(ctx sdk.Context, coins sdk.Coins) {
	bz := fck.cdc.MustMarshalBinary(coins)
	store := ctx.KVStore(fck.key)
	store.Set(collectedFeesKey, bz)
}

// Adds to Collected Fee Pool
func (fck FeeCollectionKeeper) addCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Plus(coins)
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}

func (fck FeeCollectionKeeper) refundCollectedFees(ctx sdk.Context, coins sdk.Coins) sdk.Coins {
	newCoins := fck.GetCollectedFees(ctx).Minus(coins)
	if !newCoins.IsNotNegative() {
		panic("fee collector contains negative coins")
	}
	fck.setCollectedFees(ctx, newCoins)

	return newCoins
}


// Clears the collected Fee Pool
func (fck FeeCollectionKeeper) ClearCollectedFees(ctx sdk.Context) {
	fck.setCollectedFees(ctx, sdk.Coins{})
}

func (fck FeeCollectionKeeper) FeePreprocess(ctx sdk.Context, coins sdk.Coins, gasLimit int64) sdk.Error {
	if gasLimit <= 0 {
		return sdk.ErrInternal(fmt.Sprintf("gaslimit %d should be larger than 0", gasLimit))
	}
	nativeFeeToken, err := fck.getter.GetString(ctx, NativeFeeTokenKey)
	if err != nil {
		panic(err)
	}
	nativeGasPriceThreshold, err := fck.getter.GetString(ctx, NativeGasPriceThresholdKey)
	if err != nil {
		panic(err)
	}
	threshold, ok := sdk.NewIntFromString(nativeGasPriceThreshold)
	if !ok {
		panic(errors.New("failed to parse gas price from string"))
	}

	equivalentTotalFee := sdk.NewInt(0)
	for _,coin := range coins {
		if coin.Denom != nativeFeeToken {
			exchangeRateKey := FeeExchangeRatePrefix + coin.Denom
			rateString, err := fck.getter.GetString(ctx, exchangeRateKey)
			if err != nil {
				continue
			}
			rate, ok := sdk.NewIntFromString(rateString)
			if !ok {
				panic(errors.New("failed to parse rate from string"))
			}
			equivalentFee := rate.Div(sdk.NewInt(RatePrecision)).Mul(coin.Amount)
			equivalentTotalFee = equivalentTotalFee.Add(equivalentFee)

		} else {
			equivalentTotalFee = equivalentTotalFee.Add(coin.Amount)
		}
	}

	gasPrice := equivalentTotalFee.Div(sdk.NewInt(gasLimit))
	if gasPrice.LT(threshold) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("gas price %s is less than threshold %s", gasPrice.String(), threshold.String()))
	}
	return nil
}

type GenesisState struct {
	FeeTokenNative string `json:"fee_token_native"`
	GasPriceThreshold int64 `json:"gas_price_threshold"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		FeeTokenNative: "steak",
		GasPriceThreshold: 20000000000, //2*10^10
	}
}

func InitGenesis(ctx sdk.Context, setter params.Setter, data GenesisState) {
	setter.SetString(ctx, NativeFeeTokenKey, data.FeeTokenNative)
	setter.SetString(ctx, NativeGasPriceThresholdKey, sdk.NewInt(data.GasPriceThreshold).String())
}