package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
)

// register REST routes
func RegisterRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, authcmd.GetAccountDecoder(cdc), ctx),
	).Methods("GET")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.GetAccAddressBech32(bech32addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := ctx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query account. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		// print out whole account
		output, err := cdc.MarshalJSON(account)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall query result. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

// register to Cosmos-LCD swagger routes
func RegisterLCDRoutes(routerGroup *gin.RouterGroup, ctx context.CoreContext, cdc *wire.Codec, storeName string) {
	routerGroup.GET("accounts/:address",QueryKeysRequestHandlerFn(storeName,cdc,authcmd.GetAccountDecoder(cdc),ctx))
}

// @Summary Query account information
// @Description Get the detailed information for specific address
// @Tags ICS20
// @Accept  json
// @Produce  json
// @Param address path string false "address"
// @Success 200 {object} auth.BaseAccount
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /ICS20/accounts/{address} [get]
func QueryKeysRequestHandlerFn(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder, ctx context.CoreContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		bech32addr := gtx.Param("address")

		addr, err := sdk.GetAccAddressBech32(bech32addr)
		if err != nil {
			httputil.NewError(gtx, http.StatusConflict, err)
			return
		}

		res, err := ctx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			gtx.Writer.Write([]byte(fmt.Sprintf("couldn't query account. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			httputil.NewError(gtx, http.StatusNoContent, err)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, err)
			gtx.Writer.Write([]byte(fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		httputil.Response(gtx,account)
	}
}