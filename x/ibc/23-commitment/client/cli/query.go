package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/merkle"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for IBC connections
func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	ics23CommitmentQueryCmd := &cobra.Command{
		Use:                        "commitment",
		Short:                      "IBC connection query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	ics23CommitmentQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryProof(storeKey, cdc),
	)...)
	return ics23CommitmentQueryCmd
}

func GetCmdQueryProof(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proof [path] [height]",
		Short: "Query the commitment path proof of the running chain",
		Long: strings.TrimSpace(fmt.Sprintf(`Query the commitment path
		
Example:
$ %s query ibc commit path height
		`, version.ClientName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proofHeight, _ := strconv.ParseInt(args[1], 10, 64)

			proof, err := cliCtx.QueryStoreProof([]byte(args[0]), storeKey, proofHeight-1)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(merkle.Proof{Proof: proof, Key: []byte(args[0])})
		},
	}
}
