package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/keeper"
)

func HandleChanOpenInit(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelOpenInit) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanOpenInit(
		ctx,
		msg.Channel.Ordering,
		msg.Channel.ConnectionHops,
		msg.PortID,
		msg.ChannelID,
		msg.Channel.Counterparty,
		msg.Channel.Version,
	)
	_ = channel.HandleMsgChannelOpenInit(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleChanOpenTry(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelOpenTry) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanOpenTry(
		ctx,
		msg.Channel.Ordering,
		msg.Channel.ConnectionHops,
		msg.PortID,
		msg.ChannelID,
		msg.Channel.Counterparty,
		msg.Channel.Version,
		msg.CounterpartyVersion,
	)
	_ = channel.HandleMsgChannelOpenTry(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleChanOpenAck(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelOpenAck) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanOpenAck(
		ctx,
		msg.PortID,
		msg.ChannelID,
		msg.CounterpartyVersion,
	)
	_ = channel.HandleMsgChannelOpenAck(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleChanOpenConfirm(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelOpenConfirm) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanOpenConfirm(
		ctx,
		msg.PortID,
		msg.ChannelID,
	)
	_ = channel.HandleMsgChannelOpenConfirm(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleChanCloseInit(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelCloseInit) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanCloseInit(
		ctx,
		msg.PortID,
		msg.ChannelID,
	)
	_ = channel.HandleMsgChannelCloseInit(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleChanCloseConfirm(ctx sdk.Context, k keeper.Keeper, msg types.MsgChannelCloseConfirm) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.PortID)
	// TODO: not found
	module.OnChanCloseConfirm(
		ctx,
		msg.PortID,
		msg.ChannelID,
	)
	_ = channel.HandleMsgChannelCloseConfirm(ctx, k.ChannelKeeper, msg)
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
