package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	"github.com/spf13/cobra"
)

// TODO: get proofs
// const (
// 	FlagProve = "prove"
// )

// GetQueryCmd returns the query commands for IBC connections
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ics03ConnectionQueryCmd := &cobra.Command{
		Use:                        "connection",
		Short:                      "IBC connection query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	ics03ConnectionQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryConnection(storeKey, cdc),
		GetCmdQueryConnectionProof(storeKey, cdc),
		GetCmdQueryClientConnections(storeKey, cdc),
	)...)
	return ics03ConnectionQueryCmd
}

// GetCmdQueryConnection defines the command to query a connection end
func GetCmdQueryConnection(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end [connection-id]",
		Short: "Query stored connection end",
		Long: strings.TrimSpace(fmt.Sprintf(`Query stored connection end
		
Example:
$ %s query ibc connection end [connection-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			connectionID := args[0]

			//bz, err := cdc.MarshalJSON(types.NewQueryConnectionParams(connectionID))
			//if err != nil {
			//	return err
			//}

			res, _, err := cliCtx.QueryStore(append([]byte("connection/"), connection.KeyConnection(connectionID)...), "ibc")
			if err != nil {
				return err
			}

			var connection connection.ConnectionEnd
			if err := cdc.UnmarshalBinaryLengthPrefixed(res, &connection); err != nil {
				return err
			}

			return cliCtx.PrintOutput(connection)
		},
	}
	// cmd.Flags().Bool(FlagProve, false, "(optional) show proofs for the query results")

	return cmd
}

// GetCmdQueryClientConnections defines the command to query a client connections
func GetCmdQueryClientConnections(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "client [client-id]",
		Short: "Query stored client connection paths",
		Long: strings.TrimSpace(fmt.Sprintf(`Query stored client connection paths
		
Example:
$ %s query ibc connection client [client-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			clientID := args[0]

			//bz, err := cdc.MarshalJSON(types.NewQueryClientConnectionsParams(clientID))
			//if err != nil {
			//	return err
			//}

			res, _, err := cliCtx.QueryStore(append([]byte("connection/"), connection.KeyClientConnections(clientID)...), "ibc")

			if err != nil {
				return err
			}

			var connectionPaths []string
			if err := cdc.UnmarshalBinaryLengthPrefixed(res, &connectionPaths); err != nil {
				return err
			}

			return cliCtx.PrintOutput(connectionPaths)
		},
	}
}

// GetCmdQueryConnectionProof defines the command to query a connection end proof
func GetCmdQueryConnectionProof(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [connection-id] [proof-height]",
		Short: "Query connection end proof",
		Long: strings.TrimSpace(fmt.Sprintf(`Query connection end proof
		
Example:
$ %s query ibc connection proof [connection-id] [proof-height]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			connectionID := args[0]
			proofHeight, _ := strconv.ParseInt(args[1], 10, 64)

			connProof, err := cliCtx.QueryStoreProof(append([]byte("connection/"), connection.KeyConnection(connectionID)...), "ibc", proofHeight)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(connProof)
		},
	}

	return cmd
}
