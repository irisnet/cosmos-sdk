package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	ics04 "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

// RegisterRestRoutes registers the REST routes for this module
func RegisterRestRoutes(cliCtx context.CLIContext, r *mux.Router) {
	ics04.RegisterRESTRoutes(cliCtx, r)
}
