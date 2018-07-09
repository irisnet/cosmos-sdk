package rest

import (
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/client"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
	"encoding/base64"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc("/accounts/{address}/send", SendRequestHandlerFn(cdc, kb, ctx)).Methods("POST")
	r.HandleFunc("/create_transfer", CreateTransferTransaction(cdc, ctx)).Methods("POST")
	r.HandleFunc("/signed_transfer", BroadcastSignedTransferTransaction(cdc, ctx)).Methods("POST")
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
func SendRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// collect data
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		address, err := sdk.GetAccAddressBech32(bech32addr)
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

		to, err := sdk.GetAccAddressHex(address.String())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// build message
		msg := client.BuildMsg(info.GetPubKey().Address(), to, m.Amount)
		if err != nil { // XXX rechecking same error ?
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// add gas to context
		ctx = ctx.WithGas(m.Gas)
		// add chain-id to context
		ctx = ctx.WithChainID(m.ChainID)

		// sign
		ctx = ctx.WithAccountNumber(m.AccountNumber)
		ctx = ctx.WithSequence(m.Sequence)
		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, []sdk.Msg{msg}, cdc)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// send
		res, err := ctx.BroadcastTx(txBytes)
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

func CreateTransferTransaction(cdc *wire.Codec, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var transferBody transferBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = msgCdc.UnmarshalJSON(body, &transferBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		amount := sdk.NewCoin(transferBody.Denomination,transferBody.Amount)

		var amounts sdk.Coins
		amounts = append(amounts,amount)

		fromAddress,err := sdk.GetAccAddressBech32(transferBody.FromAddress)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		toAddress,err := sdk.GetAccAddressBech32(transferBody.ToAddress)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		// build message
		msg := client.BuildMsg(fromAddress, toAddress, amounts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		accountNumber := transferBody.AccountNumber
		sequence := transferBody.Sequence
		gas := transferBody.Gas

		if transferBody.EnsureAccAndSeq {
			if ctx.Decoder == nil {
				ctx = ctx.WithDecoder(authcmd.GetAccountDecoder(cdc))
			}
			accountNumber,err = ctx.GetAccountNumber(fromAddress)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			sequence,err = ctx.NextSequence(fromAddress)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}

		txByteForSign, err := ctx.BuildTransaction(accountNumber, sequence, gas, msg, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		base64TxData := make([]byte, base64.StdEncoding.EncodedLen(len(txByteForSign)))
		base64.StdEncoding.Encode(base64TxData,txByteForSign)

		w.Write(base64TxData)
	}
}

func BroadcastSignedTransferTransaction(cdc *wire.Codec, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var signedTransaction signedBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err = msgCdc.UnmarshalJSON(body, &signedTransaction); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var txData []byte
		if _,err = base64.StdEncoding.Decode(txData, signedTransaction.TransactionData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var signatures [][]byte
		for _,data := range signedTransaction.Signatures {
			var base64DecodedData []byte
			if _,err = base64.StdEncoding.Decode(base64DecodedData, data); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			signatures = append(signatures, base64DecodedData)
		}

		var publicKeys [][]byte
		for _,data := range signedTransaction.PublicKeys {
			var base64DecodedData []byte
			if _,err := base64.StdEncoding.Decode(base64DecodedData, data); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			publicKeys = append(publicKeys, base64DecodedData)
		}

		res, err := ctx.BroadcastTransaction(txData,signatures,publicKeys,cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func RegisterLCDRoutes(routerGroup *gin.RouterGroup, ctx context.CoreContext, cdc *wire.Codec, kb keys.Keybase) {
	routerGroup.POST("/accounts/:address/send", SendAssetWithKeystoreHandlerFn(cdc, ctx, kb))
	routerGroup.POST("/create_transfer", CreateTransferTransactionFn(cdc, ctx))
	routerGroup.POST("/signed_transfer", BroadcastSignedTransferTransactionFn(cdc, ctx))
}

// @Summary Build transaction
// @Description Build transaction
// @Tags ICS20
// @Accept  json
// @Produce  json
// @Param transferBody body transferBody true "create transaction parameters"
// @Success 200 {string} string
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /ICS20/create_transfer [post]
func CreateTransferTransactionFn(cdc *wire.Codec, ctx context.CoreContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var transferBody transferBody
		if err := gtx.ShouldBindJSON(&transferBody); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		amount := sdk.NewCoin(transferBody.Denomination,transferBody.Amount)

		var amounts sdk.Coins
		amounts = append(amounts,amount)

		fromAddress,err := sdk.GetAccAddressBech32(transferBody.FromAddress)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}
		toAddress,err := sdk.GetAccAddressBech32(transferBody.ToAddress)
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
			if ctx.Decoder == nil {
				ctx = ctx.WithDecoder(authcmd.GetAccountDecoder(cdc))
			}
			accountNumber,err = ctx.GetAccountNumber(fromAddress)
			if err != nil {
				httputil.NewError(gtx, http.StatusInternalServerError, err)
				return
			}
			sequence,err = ctx.NextSequence(fromAddress)
			if err != nil {
				httputil.NewError(gtx, http.StatusInternalServerError, err)
				return
			}
		}

		txByteForSign, err := ctx.BuildTransaction(accountNumber, sequence, gas, msg, cdc)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		base64TxData := make([]byte, base64.StdEncoding.EncodedLen(len(txByteForSign)))
		base64.StdEncoding.Encode(base64TxData,txByteForSign)

		gtx.Status(http.StatusOK)
		gtx.Writer.Write(base64TxData)
	}
}

// @Summary Broadcast signed transaction
// @Description Broadcast signed transaction
// @Tags ICS20
// @Accept  json
// @Produce  json
// @Param signedTransaction body signedBody true "signed transaction"
// @Success 200 {object} ResultBroadcastTxCommit
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /ICS20/signed_transfer [post]
func BroadcastSignedTransferTransactionFn(cdc *wire.Codec, ctx context.CoreContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		var signedTransaction signedBody
		if err := gtx.ShouldBindJSON(&signedTransaction); err != nil {
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

		res, err := ctx.BroadcastTransaction(txData, signatures, publicKeys, cdc)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}
		gtx.JSON(http.StatusOK, res)
	}
}

// SendRequestHandlerFn - http request handler to send coins to a address
// @Summary Send coins to a address
// @Description This API require the Cosmos-LCD have keystore module. It will ask keystore module for transaction signature
// @Tags ICS20
// @Accept  json
// @Produce  json
// @Param sendAsset body sendBody true "transfer asset"
// @Success 200 {object} ResultBroadcastTxCommit
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts/{address}/send [post]
func SendAssetWithKeystoreHandlerFn(cdc *wire.Codec, ctx context.CoreContext, kb keys.Keybase) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		bech32addr := gtx.Param("address")

		address, err := sdk.GetAccAddressBech32(bech32addr)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		var m sendBody
		if err := gtx.ShouldBindJSON(&m); err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			httputil.NewError(gtx, http.StatusUnauthorized, err)
			return
		}

		to, err := sdk.GetAccAddressHex(address.String())
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		// build message
		msg := client.BuildMsg(info.GetPubKey().Address(), to, m.Amount)
		if err != nil { // XXX rechecking same error ?
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		// add gas to context
		ctx = ctx.WithGas(m.Gas)
		// add chain-id to context
		ctx = ctx.WithChainID(m.ChainID)

		// sign
		ctx = ctx.WithAccountNumber(m.AccountNumber)
		ctx = ctx.WithSequence(m.Sequence)
		txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, []sdk.Msg{msg}, cdc)
		if err != nil {
			httputil.NewError(gtx, http.StatusUnauthorized, err)
			return
		}

		// send
		res, err := ctx.BroadcastTx(txBytes)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			return
		}

		gtx.JSON(http.StatusOK, res)
	}
}