package channel

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/client/cli"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/client/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// Name returns the ibc-channel module's name
func Name() string {
	return SubModuleName
}

// RegisterRESTRoutes registers the REST routes for the ibc-channel module.
func RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the ibc-channel module.
func GetTxCmd(cdc *codec.Codec, storeKey string) *cobra.Command {
	return cli.GetTxCmd(fmt.Sprintf("%s/%s", storeKey, SubModuleName), cdc)
}

// GetQueryCmd returns no root query command for the ibc-channel module.
func GetQueryCmd(cdc *codec.Codec, queryRoute string) *cobra.Command {
	return cli.GetQueryCmd(fmt.Sprintf("%s/%s", queryRoute, SubModuleName), cdc)
}
