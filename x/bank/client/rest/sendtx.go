package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/client"
		"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
	"encoding/base64"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(queryCtx context.QueryContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", SendRequestHandlerFn(cdc, kb, queryCtx)).Methods("POST")
}

type sendBody struct {
	// fees is not used currently
	// Fees             sdk.Coin  `json="fees"`
	Amount           sdk.Coins `json:"amount"`
	LocalAccountName string    `json:"name"`
	Password         string    `json:"password"`
	ChainID          string    `json:"chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              int64     `json:"gas"`
}

type transferBody struct {
	ChainID         string  `json:"chain_id"`
	FromAddress		string	`json:"from_address"`
	ToAddress		string	`json:"to_address"`
	Amount			int64 	`json:"amount"`
	Denomination 	string 	`json:"denomination"`
	AccountNumber	int64	`json:"account_number"`
	Sequence		int64	`json:"sequence"`
	EnsureAccAndSeq bool 	`json:"ensure_account_sequence"`
	Gas				int64	`json:"gas"`
}

type signedBody struct {
	TransactionData	[]byte		`json:"transaction_data"`
	Signatures		[][]byte	`json:"signature_list"`
	PublicKeys		[][]byte	`json:"public_key_list"`
}

var msgCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgCdc)
}

// SendRequestHandlerFn - http request handler to send coins to a address
func SendRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, queryCtx context.QueryContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// collect data
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var m sendBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = msgCdc.UnmarshalJSON(body, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := client.BuildMsg(sdk.AccAddress(info.GetPubKey().Address()), to, m.Amount)
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			Gas:           m.Gas,
			ChainID:       m.ChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
		}

		txBytes, err := txCtx.BuildAndSign(m.LocalAccountName, m.Password, []sdk.Msg{msg})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := queryCtx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := wire.MarshalJSONIndent(cdc, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func RegisterLCDRoutes(routerGroup *gin.RouterGroup, ctx context.QueryContext, cdc *wire.Codec, kb keys.Keybase) {
	routerGroup.POST("/accounts/:address/send", SendRequestFn(cdc, ctx, kb))
	routerGroup.POST("/create_transfer", CreateTransferTransactionFn(cdc, ctx))
	routerGroup.POST("/signed_transfer", BroadcastSignedTransferTransactionFn(cdc, ctx))
}

func CreateTransferTransactionFn(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var transferBody transferBody
		if err := gtx.BindJSON(&transferBody); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		amount := sdk.NewCoin(transferBody.Denomination,transferBody.Amount)

		var amounts sdk.Coins
		amounts = append(amounts,amount)

		fromAddress,err := sdk.AccAddressFromBech32(transferBody.FromAddress)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}
		toAddress,err := sdk.AccAddressFromBech32(transferBody.ToAddress)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}
		// build message
		msg := client.BuildMsg(fromAddress, toAddress, amounts)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		accountNumber := transferBody.AccountNumber
		sequence := transferBody.Sequence
		gas := transferBody.Gas

		if transferBody.EnsureAccAndSeq {
			if ctx.AccDecoder == nil {
				ctx = ctx.WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			}
			accountNumber,err = ctx.GetAccountNumber(fromAddress)
			if err != nil {
				httputil.NewError(gtx, http.StatusInternalServerError, err)
				return
			}
			sequence,err = ctx.GetAccountSequence(fromAddress)
			if err != nil {
				httputil.NewError(gtx, http.StatusInternalServerError, err)
				return
			}
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			Gas:           gas,
			ChainID:       transferBody.ChainID,
			AccountNumber: accountNumber,
			Sequence:      sequence,
		}

		txByteForSign, err := txCtx.BuildTransactionForSign(msg)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		base64TxData := make([]byte, base64.StdEncoding.EncodedLen(len(txByteForSign)))
		base64.StdEncoding.Encode(base64TxData,txByteForSign)

		httputil.Response(gtx,string(base64TxData))
	}
}

func BroadcastSignedTransferTransactionFn(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var signedTransaction signedBody
		if err := gtx.BindJSON(&signedTransaction); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		var txData []byte
		if _,err := base64.StdEncoding.Decode(txData, signedTransaction.TransactionData); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		var signatures [][]byte
		for _,data := range signedTransaction.Signatures {
			var base64DecodedData []byte
			if _,err := base64.StdEncoding.Decode(base64DecodedData, data); err != nil {
				httputil.NewError(gtx, http.StatusBadRequest, err)
				return
			}
			signatures = append(signatures, base64DecodedData)
		}

		var publicKeys [][]byte
		for _,data := range signedTransaction.PublicKeys {
			var base64DecodedData []byte
			if _,err := base64.StdEncoding.Decode(base64DecodedData, data); err != nil {
				httputil.NewError(gtx, http.StatusBadRequest, err)
				return
			}
			publicKeys = append(publicKeys, base64DecodedData)
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
		}

		txDataForBroadcast, err := txCtx.BuildTransaction(cdc, txData,signatures,publicKeys)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}
		res, err := ctx.BroadcastTx(txDataForBroadcast)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		httputil.Response(gtx,res)
	}
}

func SendRequestFn(cdc *wire.Codec, ctx context.QueryContext, kb keys.Keybase) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		bech32addr := gtx.Param("address")

		address, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		var m sendBody
		if err := gtx.BindJSON(&m); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			httputil.NewError(gtx, http.StatusUnauthorized, err)
			return
		}

		from := sdk.AccAddress(info.GetPubKey().Address())

		to, err := sdk.AccAddressFromBech32(address.String())
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		// build message
		msg := client.BuildMsg(from, to, m.Amount)
		if err != nil { // XXX rechecking same error ?
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			Gas:           m.Gas,
			ChainID:       m.ChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
		}

		txBytes, err := txCtx.BuildAndSign(m.LocalAccountName, m.Password, []sdk.Msg{msg})
		if err != nil {
			httputil.NewError(gtx, http.StatusUnauthorized, err)
			return
		}

		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		httputil.Response(gtx,res)
	}
}