package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	ics02 "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	ics03 "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	ics04 "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	ics23 "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/client/cli"
	mockbank "github.com/cosmos/cosmos-sdk/x/ibc/mock/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(
		ics02.GetTxCmd(cdc, storeKey),
		ics03.GetTxCmd(cdc, storeKey),
		ics04.GetTxCmd(cdc, storeKey),
		mockbank.GetTxCmd(cdc),
	)
	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group ibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the IBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(
		ics02.GetQueryCmd(cdc, queryRoute),
		ics03.GetQueryCmd(cdc, queryRoute),
		ics04.GetQueryCmd(cdc, queryRoute),
		ics23.GetQueryCmd(queryRoute, cdc),
	)
	return ibcQueryCmd
}
