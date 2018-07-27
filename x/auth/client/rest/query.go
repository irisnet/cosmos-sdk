package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
	"errors"
	"github.com/gorilla/mux"
)

// register REST routes
func RegisterRoutes(queryCtx context.QueryContext, r *mux.Router, cdc *wire.Codec, storeName string) {
	r.HandleFunc(
		"/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cdc, authcmd.GetAccountDecoder(cdc), queryCtx),
	).Methods("GET")
}

// query accountREST Handler
func QueryAccountRequestHandlerFn(
	storeName string, cdc *wire.Codec,
	decoder auth.AccountDecoder, queryCtx context.QueryContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		res, err := queryCtx.QueryStore(auth.AddressStoreKey(addr), storeName)
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
func RegisterLCDRoutes(routerGroup *gin.RouterGroup, ctx context.QueryContext, cdc *wire.Codec, storeName string) {
	routerGroup.GET("accounts/:address",QueryKeysRequestHandlerFn(storeName,cdc,authcmd.GetAccountDecoder(cdc),ctx))
}

func QueryKeysRequestHandlerFn(storeName string, cdc *wire.Codec, decoder auth.AccountDecoder, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		bech32addr := gtx.Param("address")

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			httputil.NewError(gtx, http.StatusConflict, err)
			return
		}

		res, err := ctx.QueryStore(auth.AddressStoreKey(addr), storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query account. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			httputil.Response(gtx,nil)
			return
		}

		// decode the value
		account, err := decoder(res)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't parse query result. Result: %s. Error: %s", res, err.Error())))
			return
		}

		httputil.Response(gtx,account)
	}
}