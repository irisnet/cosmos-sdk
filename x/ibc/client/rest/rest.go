package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	ics03 "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	ics04 "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
)

// RegisterRestRoutes registers the REST routes for this module
func RegisterRestRoutes(cliCtx context.CLIContext, r *mux.Router) {
	ics03.RegisterRESTRoutes(cliCtx, r)
	ics04.RegisterRESTRoutes(cliCtx, r)
}
