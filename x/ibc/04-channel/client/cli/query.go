package cli

import (
	"fmt"
	"strconv"
	"strings"

	cli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/spf13/cobra"
)

// TODO: get proofs
// const (
// 	FlagProve = "prove"
// )

// GetQueryCmd returns the query commands for IBC channels
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ics04ChannelQueryCmd := &cobra.Command{
		Use:                "channel",
		Short:              "IBC channel query subcommands",
		DisableFlagParsing: true,
	}

	ics04ChannelQueryCmd.AddCommand(cli.GetCommands(
		GetCmdQueryChannel(storeKey, cdc),
		GetCmdQueryChannelProof(storeKey, cdc),
	)...)

	return ics04ChannelQueryCmd
}

// GetCmdQueryChannel defines the command to query a channel end
func GetCmdQueryChannel(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end [port-id] [channel-id]",
		Short: "Query stored connection",
		Long: strings.TrimSpace(fmt.Sprintf(`Query stored connection end
		
Example:
$ %s query ibc channel end [port-id] [channel-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]
			channelID := args[1]

			res, _, err := cliCtx.QueryStore(append([]byte("channels/"), channel.KeyChannel(portID, channelID)...), "ibc")
			if err != nil {
				return err
			}

			var channel channel.Channel
			if err := cdc.UnmarshalBinaryLengthPrefixed(res, &channel); err != nil {
				return err
			}

			return cliCtx.PrintOutput(channel)
		},
	}

	// cmd.Flags().Bool(FlagProve, false, "(optional) show proofs for the query results")

	return cmd
}

// GetCmdQueryChannelProof defines the command to query a channel proof
func GetCmdQueryChannelProof(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [port-id] [channel-id] [proof-height]",
		Short: "Query channel proof",
		Long: strings.TrimSpace(fmt.Sprintf(`Query channel proof
		
Example:
$ %s query ibc channel proof [channel-id] [proof-height]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]
			channelID := args[1]
			proofHeight, _ := strconv.ParseInt(args[1], 10, 64)

			channProof, err := cliCtx.QueryStoreProof(append([]byte("connection/"), channel.KeyChannel(portID, channelID)...), "ibc", proofHeight)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(channProof)
		},
	}

	return cmd
}
