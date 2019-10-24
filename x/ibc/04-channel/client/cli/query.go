package cli

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	cli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/merkle"
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
		GetCmdQueryChannels(storeKey, cdc),
		GetCmdQueryChannelProof(storeKey, cdc),
		GetCmdQueryPacketProof(storeKey, cdc),
	)...)

	return ics04ChannelQueryCmd
}

// GetCmdQueryChannel defines the command to query a channel end
func GetCmdQueryChannel(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end [port-id] [channel-id]",
		Short: "Query stored channel",
		Long: strings.TrimSpace(fmt.Sprintf(`Query stored channel
		
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

// GetCmdQueryChannel defines the command to query channel ends
func GetCmdQueryChannels(storeKey string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ends [port-id]",
		Short: "Query stored channels",
		Long: strings.TrimSpace(fmt.Sprintf(`Query stored channels
		
Example:
$ %s query ibc channel ends [port-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]

			subspace := []byte(fmt.Sprintf("channels/ports/%s/channels/", portID))
			resKVs, _, err := cliCtx.QuerySubspace(subspace, "ibc")
			if err != nil {
				return err
			}

			var channels []channel.Channel
			for _, kv := range resKVs {
				key := kv.Key[len(subspace):]
				if !bytes.Contains(key, []byte("/")) {
					var channel channel.Channel
					cliCtx.Codec.MustUnmarshalBinaryLengthPrefixed(kv.Value, &channel)
					channels = append(channels, channel)
				}
			}

			return cliCtx.PrintOutput(channels)
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
			proofHeight, _ := strconv.ParseInt(args[2], 10, 64)

			key := append([]byte(channel.SubModuleName+"/"), channel.KeyChannel(portID, channelID)...)
			channProof, err := cliCtx.QueryStoreProof(key, "ibc", proofHeight-1)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(merkle.Proof{Proof: channProof, Key: key})
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

			key := append([]byte(channel.SubModuleName+"/"), channel.KeyPacketCommitment(portID, channelID, Sequence)...)
			packetProof, err := cliCtx.QueryStoreProof(key, "ibc", proofHeight-1)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(merkle.Proof{Proof: packetProof, Key: key})
		},
	}

	return cmd
}
