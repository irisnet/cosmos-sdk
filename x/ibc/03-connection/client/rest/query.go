package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/ibc/connections/{%s}", RestConnectionID), queryConnectionHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/ibc/clients/{%s}/connections", RestClientID), queryConnectionsHandlerFn(cliCtx)).Methods("GET")
}

func queryConnectionHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		connectionID := vars[RestConnectionID]

		if len(connectionID) == 0 {
			err := errors.New("connection ID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryStore(append([]byte(types.SubModuleName+"/"), types.KeyConnection(connectionID)...), "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var connection types.ConnectionEnd
		if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &connection); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, connection)
	}
}

func queryConnectionsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars[RestClientID]

		if len(clientID) == 0 {
			err := errors.New("client ID required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryStore(append([]byte(types.SubModuleName+"/"), types.KeyClientConnections(clientID)...), "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var connectionPaths []string
		if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &connectionPaths); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var connections []types.ConnectionEnd

		if len(connectionPaths) > 0 {
			for _, p := range connectionPaths {
				res, _, err := cliCtx.QueryStore(append([]byte(types.SubModuleName+"/"), types.KeyConnection(p)...), "ibc")
				if err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}

				var connection types.ConnectionEnd
				if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &connection); err != nil {
					rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
					return
				}

				connections = append(connections, connection)
			}
		}

		rest.PostProcessResponse(w, cliCtx, connections)
	}
}
