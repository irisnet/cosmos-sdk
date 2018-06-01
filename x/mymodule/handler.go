package mymodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDo:
			return handleMsgDo(ctx, keeper, msg)
		case MsgUndo:
			return handleMsgUndo(ctx, keeper, msg)
			errMsg := "Unrecognized mymodule Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDo(ctx sdk.Context, keeper Keeper, msg MsgDo) sdk.Result{
	keeper = keeper
	return sdk.Result{}
}

func handleMsgUndo(ctx sdk.Context, keeper Keeper, msg MsgUndo) sdk.Result{
    keeper = keeper
	return sdk.Result{}
}