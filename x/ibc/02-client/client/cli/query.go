package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/exported"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
	ics23 "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/merkle"
	ibcversion "github.com/cosmos/cosmos-sdk/x/ibc/version"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GetQueryCmd returns the query commands for IBC clients
func GetQueryCmd(queryRouter string, cdc *codec.Codec) *cobra.Command {
	ics02ClientQueryCmd := &cobra.Command{
		Use:                        "client",
		Short:                      "IBC client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	ics02ClientQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryConsensusStateInit(queryRouter, cdc),
		GetCmdQueryConsensusState(queryRouter, cdc),
		GetCmdQueryPath(queryRouter, cdc),
		GetCmdQueryHeader(cdc),
		GetCmdQueryClientState(queryRouter, cdc),
		GetCmdQueryRoot(queryRouter, cdc),
		GetCmdQuerySelfConsensusState(cdc),
	)...)
	return ics02ClientQueryCmd
}

// GetCmdQueryClientState defines the command to query the state of a client with
// a given id as defined in https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics#query
func GetCmdQueryClientState(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "state [client-id]",
		Short: "Query a client state",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query stored client state

Example:
$ %s query ibc client state [client-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			clientID := args[0]
			if strings.TrimSpace(clientID) == "" {
				return errors.New("client ID can't be blank")
			}

			bz, err := cdc.MarshalJSON(types.NewQueryClientStateParams(clientID))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/client_state", queryRoute), bz)

			// println(types.ClientStatePath(clientID))
			// res, _, err := cliCtx.QueryWithData("/custom/ibc/client_state", bz)

			if err != nil {
				return err
			}

			var clientState types.ClientState
			if err := cdc.UnmarshalJSON(res, &clientState); err != nil {
				return err
			}

			return cliCtx.PrintOutput(clientState)
		},
	}
}

// GetCmdQueryRoot defines the command to query
func GetCmdQueryRoot(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "root [client-id] [height]",
		Short: "Query stored root",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a stored commitment root at a specific height for a particular client

Example:
$ %s query ibc client root [client-id] [height]
`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			clientID := args[0]
			if strings.TrimSpace(clientID) == "" {
				return errors.New("client ID can't be blank")
			}

			height, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("expected integer height, got: %v", args[1])
			}

			bz, err := cdc.MarshalJSON(types.NewQueryCommitmentRootParams(clientID, height))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/commitment_path", queryRoute), bz)
			if err != nil {
				return err
			}

			var root ics23.Root
			if err := cdc.UnmarshalJSON(res, &root); err != nil {
				return err
			}

			return cliCtx.PrintOutput(root)
		},
	}
}

func GetCmdQueryConsensusStateInit(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "consensus-state-init",
		Short: "Query the latest consensus state of the running chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)

			node, err := ctx.GetNode()
			if err != nil {
				return err
			}

			info, err := node.ABCIInfo()
			if err != nil {
				return err
			}

			height := info.Response.LastBlockHeight
			prevheight := height - 1

			commit, err := node.Commit(&height)
			if err != nil {
				return err
			}

			validators, err := node.Validators(&prevheight)
			if err != nil {
				return err
			}

			state := tendermint.ConsensusState{
				ChainID:          commit.ChainID,
				Height:           uint64(commit.Height),
				Root:             merkle.NewRoot(commit.AppHash),
				NextValidatorSet: tmtypes.NewValidatorSet(validators.Validators),
			}

			return ctx.PrintOutput(state)
		},
	}
}

// GetCmdQueryConsensusState defines the command to query the consensus state of
// the chain as defined in https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics#query
func GetCmdQueryConsensusState(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "consensus-state [client-id]",
		Short: "Query the latest consensus state of the client",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the consensus state for a particular client

Example:
$ %s query ibc client consensus-state [client-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			clientID := args[0]
			if strings.TrimSpace(clientID) == "" {
				return errors.New("client ID can't be blank")
			}

			bz, err := cdc.MarshalJSON(types.NewQueryClientStateParams(clientID))
			if err != nil {
				return err
			}

			// res, _, err := cliCtx.QueryWithData(types.ConsensusStatePath(clientID), bz)
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/consensus_state", queryRoute), bz)
			// TODO: bug fix

			if err != nil {
				return err
			}

			var consensusState exported.ConsensusState
			if err := cdc.UnmarshalJSON(res, &consensusState); err != nil {
				return err
			}

			return cliCtx.PrintOutput(consensusState)
		},
	}
}

// GetCmdQueryPath defines the command to query the commitment path
func GetCmdQueryPath(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Query the commitment path of the running chain",
		Long: strings.TrimSpace(fmt.Sprintf(`Query the commitment path
		
Example:
$ %s query ibc client path
		`, version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			prefixVersion := ibcversion.Prefix(ibcversion.Version)
			keyprefix := make([]byte, len(prefixVersion))
			copy(keyprefix, prefixVersion)
			path := merkle.NewPrefix([][]byte{[]byte(queryRoute)}, keyprefix)

			return cliCtx.PrintOutput(path)
		},
	}
}

// GetCmdQueryHeader defines the command to query the latest header on the chain
// TODO: do we really need this cmd ??
func GetCmdQueryHeader(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "header",
		Short: "Query the latest header of the running chain",
		Long: strings.TrimSpace(fmt.Sprintf(`Query the latest Tendermint header
		
Example:
$ %s query ibc client header
		`, version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			info, err := node.ABCIInfo()
			if err != nil {
				return err
			}

			height := info.Response.LastBlockHeight
			prevHeight := height - 1

			commit, err := node.Commit(&height)
			if err != nil {
				return err
			}

			validators, err := node.Validators(&prevHeight)
			if err != nil {
				return err
			}

			nextValidators, err := node.Validators(&height)
			if err != nil {
				return err
			}

			header := tendermint.Header{
				SignedHeader:     commit.SignedHeader,
				ValidatorSet:     tmtypes.NewValidatorSet(validators.Validators),
				NextValidatorSet: tmtypes.NewValidatorSet(nextValidators.Validators),
			}

			return cliCtx.PrintOutput(header)
		},
	}
}

// GetCmdQueryConsensusState defines the command to query the self ConsensusState
func GetCmdQuerySelfConsensusState(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "self-consensus-state",
		Short: "Query the self consensus state of the running chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext().WithCodec(cdc)

			node, err := ctx.GetNode()
			if err != nil {
				return err
			}

			info, err := node.ABCIInfo()
			if err != nil {
				return err
			}

			height := info.Response.LastBlockHeight
			prevheight := height - 1

			commit, err := node.Commit(&height)
			if err != nil {
				return err
			}

			validators, err := node.Validators(&prevheight)
			if err != nil {
				return err
			}

			state := tendermint.ConsensusState{
				ChainID:          commit.ChainID,
				Height:           uint64(commit.Height),
				Root:             merkle.NewRoot(commit.AppHash),
				NextValidatorSet: tmtypes.NewValidatorSet(validators.Validators),
			}

			return ctx.PrintOutput(state)
		},
	}
}
