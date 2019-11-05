package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	"github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/keeper"
)

func HandleConnOpenInit(ctx sdk.Context, k keeper.Keeper, msg types.MsgConnectionOpenInit) sdk.Result {
	return connection.HandleMsgConnectionOpenInit(ctx, k.ConnectionKeeper, msg)
}

func HandleConnOpenTry(ctx sdk.Context, k keeper.Keeper, msg types.MsgConnectionOpenTry) sdk.Result {
	return connection.HandleMsgConnectionOpenTry(ctx, k.ConnectionKeeper, msg)
}

func HandleConnOpenAck(ctx sdk.Context, k keeper.Keeper, msg types.MsgConnectionOpenAck) sdk.Result {
	return connection.HandleMsgConnectionOpenAck(ctx, k.ConnectionKeeper, msg)
}

func HandleConnOpenConfirm(ctx sdk.Context, k keeper.Keeper, msg types.MsgConnectionOpenConfirm) sdk.Result {
	return connection.HandleMsgConnectionOpenConfirm(ctx, k.ConnectionKeeper, msg)
}
