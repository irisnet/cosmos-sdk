package context

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/common"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tendermintLiteProxy "github.com/tendermint/tendermint/lite/proxy"
	"github.com/tendermint/iavl"
	"github.com/cosmos/cosmos-sdk/store"
	"strings"
	"github.com/tendermint/tendermint/crypto"
	"encoding/json"
)

// Broadcast the transaction bytes to Tendermint
func (ctx CoreContext) BroadcastTx(tx []byte) (*ctypes.ResultBroadcastTxCommit, error) {

	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return res, err
	}

	if res.CheckTx.Code != uint32(0) {
		return res, errors.Errorf("checkTx failed: (%d) %s",
			res.CheckTx.Code,
			res.CheckTx.Log)
	}
	if res.DeliverTx.Code != uint32(0) {
		return res, errors.Errorf("deliverTx failed: (%d) %s",
			res.DeliverTx.Code,
			res.DeliverTx.Log)
	}
	return res, err
}

// Broadcast the transaction bytes to Tendermint
func (ctx CoreContext) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {

	node, err := ctx.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BroadcastTxAsync(tx)
	if err != nil {
		return res, err
	}

	return res, err
}

// Query information about the connected node
func (ctx CoreContext) Query(path string) (res []byte, err error) {
	return ctx.query(path, nil)
}

// QueryStore from Tendermint with the provided key and storename
func (ctx CoreContext) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	return ctx.queryStore(key, storeName, "key")
}

// Query from Tendermint with the provided storename and subspace
func (ctx CoreContext) QuerySubspace(cdc *wire.Codec, subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, err := ctx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}
	cdc.MustUnmarshalBinary(resRaw, &res)
	return
}

// Query from Tendermint with the provided storename and path
func (ctx CoreContext) query(path string, key common.HexBytes) (res []byte, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, err
	}

	opts := rpcclient.ABCIQueryOptions{
		Height:  ctx.Height,
		Trusted: ctx.TrustNode,
	}
	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}
	resp := result.Response
	if resp.Code != uint32(0) {
		return res, errors.Errorf("query failed: (%d) %s", resp.Code, resp.Log)
	}

	// Data from trusted node or subspace doesn't need verification
	if ctx.TrustNode || !isQueryStoreWithProof(path) {
		return resp.Value,nil
	}

	if ctx.Cert == nil {
		return resp.Value,errors.Errorf("missing valid certifier to verify data from untrusted node")
	}

	// AppHash for height H is in header H+1
	commit, err := tendermintLiteProxy.GetCertifiedCommit(resp.Height+1, node, ctx.Cert)
	if err != nil {
		return nil, err
	}

	var rangeProof iavl.RangeProof
	cdc := wire.NewCodec()
	err = cdc.UnmarshalBinary(resp.Proof, &rangeProof)
	if err != nil {
		return res, errors.Wrap(err, "failed to unmarshalBinary rangeProof")
	}

	var multiStoreCommitInfo store.MultiStoreCommitInfo
	err = cdc.UnmarshalBinary(rangeProof.Appendix, &multiStoreCommitInfo)
	if err != nil {
		return res, errors.Wrap(err, "failed to unmarshalBinary Appendix in rangeProof")
	}
	// Validate the substore commit hash against trusted appHash
	substoreCommitHash, err :=  store.VerifyMultiStoreCommitInfo(multiStoreCommitInfo.StoreName, multiStoreCommitInfo.CommitIDList, commit.Header.AppHash)
	if err != nil {
		return  nil, errors.Wrap(err, "failed in verifying the proof against appHash")
	}
	err = store.VerifyRangeProof(resp.Key, resp.Value, substoreCommitHash, &rangeProof)
	if err != nil {
		return  nil, errors.Wrap(err, "failed in the range proof verification")
	}

	return resp.Value, nil
}

// Query from Tendermint with the provided storename and path
func (ctx CoreContext) queryStore(key cmn.HexBytes, storeName, endPath string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}

// Get the from address from the name flag
func (ctx CoreContext) GetFromAddress() (from sdk.AccAddress, err error) {

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := ctx.FromAddressName
	if name == "" {
		return nil, errors.Errorf("must provide a from address name")
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("no key for: %s", name)
	}

	return sdk.AccAddress(info.GetPubKey().Address()), nil
}

// build the transaction from the msg
func (ctx CoreContext) BuildTransaction(accnum, sequence, gas int64, msg sdk.Msg) ([]byte, error) {
	chainID := ctx.ChainID
	if chainID == "" {
		return nil, errors.Errorf("chain ID required but not specified")
	}
	memo := ctx.Memo

	signMsg := auth.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: int64(accnum),
		Sequence:      int64(sequence),
		Msgs:          []sdk.Msg{msg},
		Memo:          memo,
		Fee:           auth.NewStdFee(gas, sdk.Coin{}),
	}

	return signMsg.Bytes(),nil
}

// build the transaction from the msg
func (ctx CoreContext) BroadcastTransaction(txData []byte, signatures [][]byte, publicKeys [][]byte) (*ctypes.ResultBroadcastTxCommit, error) {
	var stdSignDoc auth.StdSignDoc//, transactionSigs []auth.StdSignature
	if err := json.Unmarshal(txData,&stdSignDoc); err != nil {
		return nil, err
	}

	var msgs []sdk.Msg
	for _, msgByte := range stdSignDoc.Msgs {
		var stdMsg sdk.Msg
		if err := json.Unmarshal(msgByte,stdMsg); err != nil {
			return nil, err
		}
		msgs = append(msgs, stdMsg)
	}

	var fee auth.StdFee
	if err := json.Unmarshal(stdSignDoc.Fee,&fee); err != nil {
		return nil, err
	}

	if len(signatures) != len(publicKeys) {
		return nil, errors.New("signatures length doesn't equal to publicKeys length")
	}

	var stdSignatures []auth.StdSignature
	for index,signature := range signatures {

		public,err := crypto.PubKeyFromBytes(publicKeys[index])
		if err != nil {
			return nil, err
		}

		sig,err := crypto.SignatureFromBytes(signature)
		if err != nil {
			return nil, err
		}

		stdSignatures = append(stdSignatures,auth.StdSignature{
			PubKey:        public,
			Signature:     sig,
			AccountNumber: stdSignDoc.AccountNumber,
			Sequence:      stdSignDoc.Sequence,
		})
	}
	// marshal bytes
	tx := auth.NewStdTx(msgs, fee, stdSignatures, stdSignDoc.Memo)

	txBytes,err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}

	return ctx.BroadcastTx(txBytes)
}

// sign and build the transaction from the msg
func (ctx CoreContext) SignAndBuild(name, passphrase string, msgs []sdk.Msg, cdc *wire.Codec) ([]byte, error) {

	// build the Sign Messsage from the Standard Message
	chainID := ctx.ChainID
	if chainID == "" {
		return nil, errors.Errorf("chain ID required but not specified")
	}
	accnum := ctx.AccountNumber
	sequence := ctx.Sequence
	memo := ctx.Memo

	fee := sdk.Coin{}
	if ctx.Fee != "" {
		parsedFee, err := sdk.ParseCoin(ctx.Fee)
		if err != nil {
			return nil, err
		}
		fee = parsedFee
	}

	signMsg := auth.StdSignMsg{
		ChainID:       chainID,
		AccountNumber: accnum,
		Sequence:      sequence,
		Msgs:          msgs,
		Memo:          memo,
		Fee:           auth.NewStdFee(ctx.Gas, fee), // TODO run simulate to estimate gas?
	}

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	// sign and build
	bz := signMsg.Bytes()

	sig, pubkey, err := keybase.Sign(name, passphrase, bz)
	if err != nil {
		return nil, err
	}
	sigs := []auth.StdSignature{{
		PubKey:        pubkey,
		Signature:     sig,
		AccountNumber: accnum,
		Sequence:      sequence,
	}}

	// marshal bytes
	tx := auth.NewStdTx(signMsg.Msgs, signMsg.Fee, sigs, memo)

	return cdc.MarshalBinary(tx)
}

// sign and build the transaction from the msg
func (ctx CoreContext) ensureSignBuild(name string, msgs []sdk.Msg, cdc *wire.Codec) (tyBytes []byte, err error) {
	ctx, err = EnsureAccountNumber(ctx)
	if err != nil {
		return nil, err
	}
	// default to next sequence number if none provided
	ctx, err = EnsureSequence(ctx)
	if err != nil {
		return nil, err
	}

	var txBytes []byte

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, err
	}
	var passphrase string
	// Only need a passphrase for locally-stored keys
	if info.GetType() == "local" {
		passphrase, err = ctx.GetPassphraseFromStdin(name)
		if err != nil {
			return nil, fmt.Errorf("Error fetching passphrase: %v", err)
		}
	}
	txBytes, err = ctx.SignAndBuild(name, passphrase, msgs, cdc)
	if err != nil {
		return nil, fmt.Errorf("Error signing transaction: %v", err)
	}

	return txBytes, err
}

// sign and build the transaction from the msg
func (ctx CoreContext) EnsureSignBuildBroadcast(name string, msgs []sdk.Msg, cdc *wire.Codec) (err error) {

	txBytes, err := ctx.ensureSignBuild(name, msgs, cdc)
	if err != nil {
		return err
	}

	if ctx.Async {
		res, err := ctx.BroadcastTxAsync(txBytes)
		if err != nil {
			return err
		}
		if ctx.JSON {
			type toJSON struct {
				TxHash string
			}
			valueToJSON := toJSON{res.Hash.String()}
			JSON, err := cdc.MarshalJSON(valueToJSON)
			if err != nil {
				return err
			}
			fmt.Println(string(JSON))
		} else {
			fmt.Println("Async tx sent. tx hash: ", res.Hash.String())
		}
		return nil
	}
	res, err := ctx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}
	if ctx.JSON {
		// Since JSON is intended for automated scripts, always include response in JSON mode
		type toJSON struct {
			Height   int64
			TxHash   string
			Response string
		}
		valueToJSON := toJSON{res.Height, res.Hash.String(), fmt.Sprintf("%+v", res.DeliverTx)}
		JSON, err := cdc.MarshalJSON(valueToJSON)
		if err != nil {
			return err
		}
		fmt.Println(string(JSON))
		return nil
	}
	if ctx.PrintResponse {
		fmt.Printf("Committed at block %d. Hash: %s Response:%+v \n", res.Height, res.Hash.String(), res.DeliverTx)
	} else {
		fmt.Printf("Committed at block %d. Hash: %s \n", res.Height, res.Hash.String())
	}
	return nil
}

// get the next sequence for the account address
func (ctx CoreContext) GetAccountNumber(address []byte) (int64, error) {
	if ctx.Decoder == nil {
		return 0, errors.New("accountDecoder required but not provided")
	}

	res, err := ctx.QueryStore(auth.AddressStoreKey(address), ctx.AccountStore)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		fmt.Printf("No account found.  Returning 0.\n")
		return 0, err
	}

	account, err := ctx.Decoder(res)
	if err != nil {
		panic(err)
	}

	return account.GetAccountNumber(), nil
}

// get the next sequence for the account address
func (ctx CoreContext) NextSequence(address []byte) (int64, error) {
	if ctx.Decoder == nil {
		return 0, errors.New("accountDecoder required but not provided")
	}

	res, err := ctx.QueryStore(auth.AddressStoreKey(address), ctx.AccountStore)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		fmt.Printf("No account found, defaulting to sequence 0\n")
		return 0, err
	}

	account, err := ctx.Decoder(res)
	if err != nil {
		panic(err)
	}

	return account.GetSequence(), nil
}

// get passphrase from std input
func (ctx CoreContext) GetPassphraseFromStdin(name string) (pass string, err error) {
	buf := client.BufferStdin()
	prompt := fmt.Sprintf("Password to sign with '%s':", name)
	return client.GetPassword(prompt, buf)
}

// GetNode prepares a simple rpc.Client
func (ctx CoreContext) GetNode() (rpcclient.Client, error) {
	if ctx.ClientMgr != nil {
		return ctx.ClientMgr.getClient(), nil
	}
	if ctx.Client == nil {
		return nil, errors.New("must define node URI")
	}
	return ctx.Client, nil
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType can be app or store
// if subpath equals to store or key, then return true
func isQueryStoreWithProof(path string) (bool) {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	paths := strings.SplitN(path[1:], "/", 3)
	if len(paths) != 3 {
		return false
	}
	// WARNING This should be consistent with query method in iavlstore.go
	if paths[2] == "store" || paths[2] == "key" {
		return true
	}
	return false
}