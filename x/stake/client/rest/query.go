package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
	"github.com/gin-gonic/gin"
	"github.com/cosmos/cosmos-sdk/client/httputil"
	"errors"
	"github.com/gorilla/mux"
)

const storeName = "stake"

func registerQueryRoutes(queryCtx context.QueryContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc(
		"/stake/{delegator}/delegation/{validator}",
		delegationHandlerFn(queryCtx, cdc),
	).Methods("GET")

	r.HandleFunc(
		"/stake/{delegator}/ubd/{validator}",
		ubdHandlerFn(queryCtx, cdc),
	).Methods("GET")

	r.HandleFunc(
		"/stake/{delegator}/red/{validator_src}/{validator_dst}",
		redHandlerFn(queryCtx, cdc),
	).Methods("GET")

	r.HandleFunc(
		"/stake/validators",
		validatorsHandlerFn(queryCtx, cdc),
	).Methods("GET")
}

// http request handler to query a delegation
func delegationHandlerFn(queryCtx context.QueryContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegator"]
		bech32validator := vars["validator"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key := stake.GetDelegationKey(delegatorAddr, validatorAddr)

		res, err := queryCtx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		delegation, err := types.UnmarshalDelegation(cdc, key, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		output, err := cdc.MarshalJSON(delegation)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// http request handler to query an unbonding-delegation
func ubdHandlerFn(queryCtx context.QueryContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegator"]
		bech32validator := vars["validator"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key := stake.GetUBDKey(delegatorAddr, validatorAddr)

		res, err := queryCtx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ubd, err := types.UnmarshalUBD(cdc, key, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(ubd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// http request handler to query an redelegation
func redHandlerFn(queryCtx context.QueryContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read parameters
		vars := mux.Vars(r)
		bech32delegator := vars["delegator"]
		bech32validatorSrc := vars["validator_src"]
		bech32validatorDst := vars["validator_dst"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		validatorSrcAddr, err := sdk.AccAddressFromBech32(bech32validatorSrc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		validatorDstAddr, err := sdk.AccAddressFromBech32(bech32validatorDst)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key := stake.GetREDKey(delegatorAddr, validatorSrcAddr, validatorDstAddr)

		res, err := queryCtx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query redelegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		red, err := types.UnmarshalRED(cdc, key, res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(red)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// TODO bech32
// http request handler to query list of validators
func validatorsHandlerFn(queryCtx context.QueryContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kvs, err := queryCtx.QuerySubspace(stake.ValidatorsKey, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query validators. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there are no validators
		if len(kvs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// parse out the validators
		validators := make([]types.BechValidator, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
				return
			}

			bech32Validator, err := validator.Bech32Validator()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			validators[i] = bech32Validator
		}

		output, err := cdc.MarshalJSON(validators)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

func RegisterQueryLCDRoutes(routerGroup *gin.RouterGroup, ctx context.QueryContext, cdc *wire.Codec) {
	routerGroup.GET("/stake/:delegator/delegation/:validator", delegationHandlerFun(cdc, ctx))
	routerGroup.GET("/stake/:delegator/ubd/:validator", ubdHandlerFun(cdc, ctx))
	routerGroup.GET("/stake/:delegator/red/:validator_src/:validator_dst", redHandlerFun(cdc, ctx))
	routerGroup.GET("/stake_validators", validatorsHandlerFun(cdc, ctx))
}

func delegationHandlerFun(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
		return func(gtx *gin.Context) {

		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validator := gtx.Param("validator")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetDelegationKey(delegatorAddr, validatorAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, errors.New(fmt.Sprintf("couldn't query delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputil.Response(gtx,nil)
			return
		}

		delegation, err := types.UnmarshalDelegation(cdc, key, res)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		httputil.Response(gtx,delegation)
	}
}

// http request handler to query an unbonding-delegation
func ubdHandlerFun(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validator := gtx.Param("validator")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetUBDKey(delegatorAddr, validatorAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputil.Response(gtx,nil)
			return
		}

		ubd, err := types.UnmarshalUBD(cdc, key, res)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		httputil.Response(gtx,ubd)
	}
}

// http request handler to query an redelegation
func redHandlerFun(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validatorSrc := gtx.Param("validator_src")
		bech32validatorDst := gtx.Param("validator_dst")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorSrcAddr, err := sdk.AccAddressFromBech32(bech32validatorSrc)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorDstAddr, err := sdk.AccAddressFromBech32(bech32validatorDst)
		if err != nil {
			httputil.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetREDKey(delegatorAddr, validatorSrcAddr, validatorDstAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query redelegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputil.Response(gtx,nil)
			return
		}

		red, err := types.UnmarshalRED(cdc, key, res)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		httputil.Response(gtx,red)
	}
}

func validatorsHandlerFun(cdc *wire.Codec, ctx context.QueryContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		kvs, err := ctx.QuerySubspace(stake.ValidatorsKey, storeName)
		if err != nil {
			httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query validators. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there are no validators
		if len(kvs) == 0 {
			httputil.Response(gtx,nil)
			return
		}

		// parse out the validators
		validators := make([]types.BechValidator, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
			if err != nil {
				httputil.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
				return
			}

			bech32Validator, err := validator.Bech32Validator()
			if err != nil {
				httputil.NewError(gtx, http.StatusBadRequest, err)
				return
			}
			validators[i] = bech32Validator
		}

		httputil.Response(gtx,validators)
	}
}