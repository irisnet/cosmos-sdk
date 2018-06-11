package mymodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

var amt sdk.Coins = []sdk.Coin{{Denom:"steak",Amount:10}}
// Handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDo:
			return handleMsgDo(ctx, keeper, msg)
		case MsgUndo:
			return handleMsgUndo(ctx, keeper, msg)
		default:
			errMsg := "Unrecognized mymodule Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDo(ctx sdk.Context, keeper Keeper, msg MsgDo) sdk.Result{
	_,_, err := keeper.ck.AddCoins(ctx,msg.Addr,amt)
	if err!= nil{
		return err.Result()
	}
	if ctx.IsCheckTx() {
		return sdk.Result{}
	}
	var i int64 =0
	for{
		_,_, err := keeper.ck.AddCoins(ctx,msg.Addr,amt)
		if err!= nil{
			return err.Result()
		}
		i=i+1
		if i < msg.ValueNum.Num{
			break
		}

	}

	keeper.SetCounter(ctx,msg.Addr,keeper.GetCounter(ctx,msg.Addr)+msg.ValueNum.Num)

	return sdk.Result{}
}

func handleMsgUndo(ctx sdk.Context, keeper Keeper, msg MsgUndo) sdk.Result{
	_,_, err := keeper.ck.SubtractCoins(ctx,msg.Addr,amt)
	if err!= nil{
		return err.Result()
	}
	if ctx.IsCheckTx() {
		return sdk.Result{}
	}

	num := keeper.GetCounter(ctx,msg.Addr)

	var i int64 =0
	for{
		_,_, err := keeper.ck.SubtractCoins(ctx,msg.Addr,amt)
		if err!= nil{
			return err.Result()
		}
		i=i+1
		if i < num{
			break
		}

	}
	keeper.SetCounter(ctx,msg.Addr,int64(0))
	return sdk.Result{}
}