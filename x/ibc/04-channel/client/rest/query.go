package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/ibc/ports/{%s}/channels", RestPortID), queryChannelsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/ibc/ports/{%s}/channels/{%s}", RestPortID, RestChannelID), queryChannelHandlerFn(cliCtx)).Methods("GET")
}

func queryChannelsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		portID := vars[RestPortID]

		if len(portID) == 0 {
			err := errors.New("port ID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QuerySubspace([]byte(fmt.Sprintf("ports/%s/channels/", portID)), "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryChannelHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		portID := vars[RestPortID]
		channelID := vars[RestChannelID]

		if len(portID) == 0 {
			err := errors.New("port ID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(channelID) == 0 {
			err := errors.New("channel ID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryStore(append([]byte(types.SubModuleName), types.KeyChannel(portID, channelID)...), "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
