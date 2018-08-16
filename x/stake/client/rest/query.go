package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/cosmos/cosmos-sdk/x/stake/tags"
	"github.com/cosmos/cosmos-sdk/x/stake/types"

	"github.com/gorilla/mux"
	"github.com/gin-gonic/gin"
	"errors"
	"github.com/cosmos/cosmos-sdk/client/httputils"
)

const storeName = "stake"

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec) {

	// Get all delegations (delegation, undelegation and redelegation) from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}",
		delegatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all staking txs (i.e msgs) from a delegator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/txs",
		delegatorTxsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all validators that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators",
		delegatorValidatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a validator that a delegator is bonded to
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/validators/{validatorAddr}",
		delegatorValidatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query a delegation between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/delegations/{validatorAddr}",
		delegationHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Query all unbonding_delegations between a delegator and a validator
	r.HandleFunc(
		"/stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}",
		unbondingDelegationsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get all validators
	r.HandleFunc(
		"/stake/validators",
		validatorsHandlerFn(cliCtx, cdc),
	).Methods("GET")

	// Get a single validator info
	r.HandleFunc(
		"/stake/validators/{addr}",
		validatorHandlerFn(cliCtx, cdc),
	).Methods("GET")

}

// already resolve the rational shares to not handle this in the client

// defines a delegation without type Rat for shares
type DelegationWithoutRat struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.AccAddress `json:"validator_addr"`
	Shares        string         `json:"shares"`
	Height        int64          `json:"height"`
}

// aggregation of all delegations, unbondings and redelegations
type DelegationSummary struct {
	Delegations          []DelegationWithoutRat      `json:"delegations"`
	UnbondingDelegations []stake.UnbondingDelegation `json:"unbonding_delegations"`
	Redelegations        []stake.Redelegation        `json:"redelegations"`
}

// HTTP request handler to query a delegator delegations
func delegatorHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var validatorAddr sdk.AccAddress
		var delegationSummary = DelegationSummary{}

		// read parameters
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Get all validators using key
		validators, statusCode, errMsg, err := getBech32Validators(storeName, cliCtx, cdc)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(fmt.Sprintf("%s%s", errMsg, err.Error())))
			return
		}

		for _, validator := range validators {
			validatorAddr = validator.Owner

			// Delegations
			delegations, statusCode, errMsg, err := getDelegatorDelegations(cliCtx, cdc, delegatorAddr, validatorAddr)
			if err != nil {
				w.WriteHeader(statusCode)
				w.Write([]byte(fmt.Sprintf("%s%s", errMsg, err.Error())))
				return
			}
			if statusCode != http.StatusNoContent {
				delegationSummary.Delegations = append(delegationSummary.Delegations, delegations)
			}

			// Undelegations
			unbondingDelegation, statusCode, errMsg, err := getDelegatorUndelegations(cliCtx, cdc, delegatorAddr, validatorAddr)
			if err != nil {
				w.WriteHeader(statusCode)
				w.Write([]byte(fmt.Sprintf("%s%s", errMsg, err.Error())))
				return
			}
			if statusCode != http.StatusNoContent {
				delegationSummary.UnbondingDelegations = append(delegationSummary.UnbondingDelegations, unbondingDelegation)
			}

			// Redelegations
			// only querying redelegations to a validator as this should give us already all relegations
			// if we also would put in redelegations from, we would have every redelegation double
			redelegations, statusCode, errMsg, err := getDelegatorRedelegations(cliCtx, cdc, delegatorAddr, validatorAddr)
			if err != nil {
				w.WriteHeader(statusCode)
				w.Write([]byte(fmt.Sprintf("%s%s", errMsg, err.Error())))
				return
			}
			if statusCode != http.StatusNoContent {
				delegationSummary.Redelegations = append(delegationSummary.Redelegations, redelegations)
			}
		}

		output, err := cdc.MarshalJSON(delegationSummary)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// nolint gocyclo
// HTTP request handler to query all staking txs (msgs) from a delegator
func delegatorTxsHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var output []byte
		var typesQuerySlice []string
		vars := mux.Vars(r)
		delegatorAddr := vars["delegatorAddr"]

		_, err := sdk.AccAddressFromBech32(delegatorAddr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		node, err := cliCtx.GetNode()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't get current Node information. Error: %s", err.Error())))
			return
		}

		// Get values from query

		typesQuery := r.URL.Query().Get("type")
		trimmedQuery := strings.TrimSpace(typesQuery)
		if len(trimmedQuery) != 0 {
			typesQuerySlice = strings.Split(trimmedQuery, " ")
		}

		noQuery := len(typesQuerySlice) == 0
		isBondTx := contains(typesQuerySlice, "bond")
		isUnbondTx := contains(typesQuerySlice, "unbond")
		isRedTx := contains(typesQuerySlice, "redelegate")
		var txs = []tx.Info{}
		var actions []string

		switch {
		case isBondTx:
			actions = append(actions, string(tags.ActionDelegate))
		case isUnbondTx:
			actions = append(actions, string(tags.ActionBeginUnbonding))
			actions = append(actions, string(tags.ActionCompleteUnbonding))
		case isRedTx:
			actions = append(actions, string(tags.ActionBeginRedelegation))
			actions = append(actions, string(tags.ActionCompleteRedelegation))
		case noQuery:
			actions = append(actions, string(tags.ActionDelegate))
			actions = append(actions, string(tags.ActionBeginUnbonding))
			actions = append(actions, string(tags.ActionCompleteUnbonding))
			actions = append(actions, string(tags.ActionBeginRedelegation))
			actions = append(actions, string(tags.ActionCompleteRedelegation))
		default:
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, action := range actions {
			foundTxs, errQuery := queryTxs(node, cdc, action, delegatorAddr)
			if errQuery != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("error querying transactions. Error: %s", errQuery.Error())))
			}
			txs = append(txs, foundTxs...)
		}

		output, err = cdc.MarshalJSON(txs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}

// HTTP request handler to query an unbonding-delegation
func unbondingDelegationsHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

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
		validatorAddrAcc := sdk.AccAddress(validatorAddr)

		key := stake.GetUBDKey(delegatorAddr, validatorAddrAcc)

		res, err := cliCtx.QueryStore(key, storeName)
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
			w.Write([]byte(fmt.Sprintf("couldn't unmarshall unbonding-delegation. Error: %s", err.Error())))
			return
		}

		// unbondings will be a list in the future but is not yet, but we want to keep the API consistent
		ubdArray := []stake.UnbondingDelegation{ubd}

		output, err := cdc.MarshalJSON(ubdArray)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't marshall unbonding-delegation. Error: %s", err.Error())))
			return
		}

		w.Write(output)
	}
}

// HTTP request handler to query a bonded validator
func delegationHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read parameters
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

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
		validatorAddrAcc := sdk.AccAddress(validatorAddr)

		key := stake.GetDelegationKey(delegatorAddr, validatorAddrAcc)

		res, err := cliCtx.QueryStore(key, storeName)
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

		outputDelegation := DelegationWithoutRat{
			DelegatorAddr: delegation.DelegatorAddr,
			ValidatorAddr: delegation.ValidatorAddr,
			Height:        delegation.Height,
			Shares:        delegation.Shares.String(),
		}

		output, err := cdc.MarshalJSON(outputDelegation)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

// HTTP request handler to query all delegator bonded validators
func delegatorValidatorsHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var validatorAccAddr sdk.AccAddress
		var bondedValidators []types.BechValidator

		// read parameters
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Get all validators using key
		kvs, err := cliCtx.QuerySubspace(stake.ValidatorsKey, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query validators. Error: %s", err.Error())))
			return
		} else if len(kvs) == 0 {
			// the query will return empty if there are no validators
			w.WriteHeader(http.StatusNoContent)
			return
		}

		validators, err := getValidators(kvs, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}

		for _, validator := range validators {
			// get all transactions from the delegator to val and append
			validatorAccAddr = validator.Owner

			validator, statusCode, errMsg, errRes := getDelegatorValidator(cliCtx, cdc, delegatorAddr, validatorAccAddr)
			if errRes != nil {
				w.WriteHeader(statusCode)
				w.Write([]byte(fmt.Sprintf("%s%s", errMsg, errRes.Error())))
				return
			} else if statusCode == http.StatusNoContent {
				continue
			}

			bondedValidators = append(bondedValidators, validator)
		}
		output, err := cdc.MarshalJSON(bondedValidators)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)
	}
}

// HTTP request handler to get information from a currently bonded validator
func delegatorValidatorHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read parameters
		var output []byte
		vars := mux.Vars(r)
		bech32delegator := vars["delegatorAddr"]
		bech32validator := vars["validatorAddr"]

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		validatorAccAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}

		// Check if there if the delegator is bonded or redelegated to the validator

		validator, statusCode, errMsg, err := getDelegatorValidator(cliCtx, cdc, delegatorAddr, validatorAccAddr)
		if err != nil {
			w.WriteHeader(statusCode)
			w.Write([]byte(fmt.Sprintf("%s%s", errMsg, err.Error())))
			return
		} else if statusCode == http.StatusNoContent {
			w.WriteHeader(statusCode)
			return
		}
		output, err = cdc.MarshalJSON(validator)
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
func validatorsHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kvs, err := cliCtx.QuerySubspace(stake.ValidatorsKey, storeName)
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

		validators, err := getValidators(kvs, cdc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
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

// HTTP request handler to query the validator information from a given validator address
func validatorHandlerFn(cliCtx context.CLIContext, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var output []byte
		// read parameters
		vars := mux.Vars(r)
		bech32validatorAddr := vars["addr"]
		valAddress, err := sdk.AccAddressFromBech32(bech32validatorAddr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
			return
		}

		key := stake.GetValidatorKey(valAddress)

		res, err := cliCtx.QueryStore(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("couldn't query validator, error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		validator, err := types.UnmarshalValidator(cdc, valAddress, res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		bech32Validator, err := validator.Bech32Validator()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		output, err = cdc.MarshalJSON(bech32Validator)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}

		if output == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Write(output)
	}
}

func RegisterSwaggerQueryRoutes(routerGroup *gin.RouterGroup, ctx context.CLIContext, cdc *wire.Codec) {
	routerGroup.GET("/stake/:delegator/delegation/:validator", delegationHandlerFun(cdc, ctx))
	routerGroup.GET("/stake/:delegator/ubd/:validator", ubdHandlerFun(cdc, ctx))
	routerGroup.GET("/stake/:delegator/red/:validator_src/:validator_dst", redHandlerFun(cdc, ctx))
	routerGroup.GET("/stake_validators", validatorsHandlerFun(cdc, ctx))
}

func delegationHandlerFun(cdc *wire.Codec, ctx context.CLIContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validator := gtx.Param("validator")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetDelegationKey(delegatorAddr, validatorAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, errors.New(fmt.Sprintf("couldn't query delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputils.Response(gtx,nil)
			return
		}

		delegation, err := types.UnmarshalDelegation(cdc, key, res)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		httputils.Response(gtx,delegation)
	}
}

// http request handler to query an unbonding-delegation
func ubdHandlerFun(cdc *wire.Codec, ctx context.CLIContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {

		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validator := gtx.Param("validator")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorAddr, err := sdk.AccAddressFromBech32(bech32validator)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetUBDKey(delegatorAddr, validatorAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputils.Response(gtx,nil)
			return
		}

		ubd, err := types.UnmarshalUBD(cdc, key, res)
		if err != nil {
			httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		httputils.Response(gtx,ubd)
	}
}

// http request handler to query an redelegation
func redHandlerFun(cdc *wire.Codec, ctx context.CLIContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		// read parameters
		bech32delegator := gtx.Param("delegator")
		bech32validatorSrc := gtx.Param("validator_src")
		bech32validatorDst := gtx.Param("validator_dst")

		delegatorAddr, err := sdk.AccAddressFromBech32(bech32delegator)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorSrcAddr, err := sdk.AccAddressFromBech32(bech32validatorSrc)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		validatorDstAddr, err := sdk.AccAddressFromBech32(bech32validatorDst)
		if err != nil {
			httputils.NewError(gtx, http.StatusBadRequest, err)
			return
		}

		key := stake.GetREDKey(delegatorAddr, validatorSrcAddr, validatorDstAddr)

		res, err := ctx.QueryStore(key, storeName)
		if err != nil {
			httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query redelegation. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this record
		if len(res) == 0 {
			httputils.Response(gtx,nil)
			return
		}

		red, err := types.UnmarshalRED(cdc, key, res)
		if err != nil {
			httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
			return
		}

		httputils.Response(gtx,red)
	}
}

func validatorsHandlerFun(cdc *wire.Codec, ctx context.CLIContext) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		kvs, err := ctx.QuerySubspace(stake.ValidatorsKey, storeName)
		if err != nil {
			httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query validators. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there are no validators
		if len(kvs) == 0 {
			httputils.Response(gtx,nil)
			return
		}

		// parse out the validators
		validators := make([]types.BechValidator, len(kvs))
		for i, kv := range kvs {

			addr := kv.Key[1:]
			validator, err := types.UnmarshalValidator(cdc, addr, kv.Value)
			if err != nil {
				httputils.NewError(gtx, http.StatusInternalServerError, errors.New(fmt.Sprintf("couldn't query unbonding-delegation. Error: %s", err.Error())))
				return
			}

			bech32Validator, err := validator.Bech32Validator()
			if err != nil {
				httputils.NewError(gtx, http.StatusBadRequest, err)
				return
			}
			validators[i] = bech32Validator
		}

		httputils.Response(gtx,validators)
	}
}