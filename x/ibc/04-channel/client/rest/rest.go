package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

// REST Variable names
// nolint
const (
	RestPortID    = "port-id"
	RestChannelID = "channel-id"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
