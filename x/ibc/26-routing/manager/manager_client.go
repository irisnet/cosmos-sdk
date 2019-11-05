package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/keeper"
)

func HandleClientCreate(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateClient) sdk.Result {
	return client.HandleMsgCreateClient(ctx, k.ClientKeeper, msg)
}

func HandleClientUpdate(ctx sdk.Context, k keeper.Keeper, msg types.MsgUpdateClient) sdk.Result {
	return client.HandleMsgUpdateClient(ctx, k.ClientKeeper, msg)
}

func HandleClientMisbehaviour(ctx sdk.Context, k keeper.Keeper, msg types.MsgSubmitMisbehaviour) sdk.Result {
	return client.HandleMsgSubmitMisbehaviour(ctx, k.ClientKeeper, msg)
}
