package rest

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/gorilla/mux"
)

var (
	ChannelKeys = []string{"keys", "nextSequenceSend", "nextSequenceRecv", "packets", "acknowledgements"}
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

		subspace := []byte(fmt.Sprintf("channels/ports/%s/channels/", portID))

		resKVs, _, err := cliCtx.QuerySubspace(subspace, "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var channels []types.Channel

		for _, kv := range resKVs {
			key := kv.Key[len(subspace):]
			if bytes.Index(key, []byte("/")) == 0 {
				var channel types.Channel
				cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(kv.Value, &channel)

				channels = append(channels, channel)
			}
		}

		rest.PostProcessResponse(w, cliCtx, channels)
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

		res, _, err := cliCtx.QueryStore(append([]byte(types.SubModuleName+"/"), types.KeyChannel(portID, channelID)...), "ibc")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var channel types.Channel
		if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &channel); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, channel)
	}
}
