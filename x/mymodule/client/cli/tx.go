package cli

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client/context"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/wire"
    authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"strconv"
	"github.com/cosmos/cosmos-sdk/x/mymodule"
)

func GetCmdDo(cdc *wire.Codec) *cobra.Command{
	 cmd := &cobra.Command{
	     Use: "do [addr] [num]",
	     Short: "just do it",
	     RunE: func(cmd *cobra.Command, args []string) error {
			 addr,err := sdk.GetAddress(args[0])
			 if err != nil{
			 	return  err
			 }

			 num, err :=strconv.ParseInt(args[1],10,64)

			 msg := mymodule.NewMsgDo(addr,mymodule.NewValueNum(num))
			 fmt.Println(msg)
			 ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			 res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg ,cdc)
			 if err != nil{
			 	return  err
			 }
			 fmt.Println("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			 return  nil

		 },
	 }
	return cmd
}

func GetCmdUndo(cdc *wire.Codec) *cobra.Command{
	cmd := &cobra.Command{
		Use: "undo [addr]",
		Short: "just undo it",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr,err := sdk.GetAddress(args[0])
			if err != nil{
				return  err
			}
			msg := mymodule.NewMsgUndo(addr)

			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg ,cdc)
			if err != nil{
				return  err
			}
			fmt.Println("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
			return  nil

		},
	}
	return cmd
}