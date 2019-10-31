package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"

	cli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

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
		GetCmdQueryPacketProof(storeKey, cdc),
	)...)

	return ics04ChannelQueryCmd
}

// GetCmdQueryChannel defines the command to query a channel end
func GetCmdQueryChannel(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end [port-id] [channel-id]",
		Short: "Query a channel end",
		Long: strings.TrimSpace(fmt.Sprintf(`Query an IBC channel end
		
Example:
$ %s query ibc channel end [port-id] [channel-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]
			channelID := args[1]

			bz, err := cdc.MarshalJSON(types.NewQueryChannelParams(portID, channelID))
			if err != nil {
				return err
			}

			req := abci.RequestQuery{
				Path:  fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryChannel),
				Data:  bz,
				Prove: viper.GetBool(flags.FlagProve),
			}

			res, err := cliCtx.QueryABCI(req)
			if err != nil {
				return err
			}

			var channel types.Channel
			if err := cdc.UnmarshalJSON(res.Value, &channel); err != nil {
				return err
			}

			if res.Proof == nil {
				return cliCtx.PrintOutput(channel)
			}

			channelRes := types.NewChannelResponse(portID, channelID, channel, res.Proof, res.Height)
			return cliCtx.PrintOutput(channelRes)
		},
	}
	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")

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
			proofHeight, _ := strconv.ParseInt(args[2], 10, 64)

			key := append([]byte(types.SubModuleName+"/"), types.KeyChannel(portID, channelID)...)
			channProof, err := cliCtx.QueryStoreProof(key, "ibc", proofHeight-1)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(commitment.Proof{Proof: channProof})
		},
	}

	return cmd
}

// GetCmdQueryPakcerProof defines the command to query a packet proof
func GetCmdQueryPacketProof(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet-proof [port-id] [channel-id] [sequence] [proof-height]",
		Short: "Query packet proof",
		Long: strings.TrimSpace(fmt.Sprintf(`Query packet proof
		
Example:
$ %s query ibc channel proof [channel-id] [sequence] [proof-height]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]
			channelID := args[1]
			Sequence, _ := strconv.ParseUint(args[2], 10, 64)
			proofHeight, _ := strconv.ParseInt(args[3], 10, 64)

			key := append([]byte(types.SubModuleName+"/"), types.KeyPacketCommitment(portID, channelID, Sequence)...)
			packetProof, err := cliCtx.QueryStoreProof(key, "ibc", proofHeight-1)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(commitment.Proof{Proof: packetProof})
		},
	}

	return cmd
}
