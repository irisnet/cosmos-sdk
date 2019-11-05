package manager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/keeper"
)

func HandlePacketRecv(ctx sdk.Context, k keeper.Keeper, msg types.MsgPacketRecv) sdk.Result {
	module, _ := k.LookupModule(ctx, msg.Packet.GetSourcePort())
	packet := msg.Packet.(types.Packet)
	acknowledgement, _ := module.OnRecvPacket(ctx, packet)
	// TODO: error
	_ = channel.HandleMsgRecvPacket(ctx, k.ChannelKeeper, msg, acknowledgement, sdk.NewKVStoreKey("")) // TODO: portCapability
	ctx.EventManager().EmitEvents(sdk.Events{
		// TODO:
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// func HandlePacketAcknowledgement(ctx sdk.Context, k keeper.Keeper, msg types.MsgPacketAcknowledgement) sdk.Result {
// 	// module = lookupModule(msg.packet.sourcePort)
// 	// module.onAcknowledgePacket(
// 	//   msg.packet,
// 	//   msg.acknowledgement
// 	// )
// 	_ = channel.HandleMsgAcknowledgePacket(ctx, k.ChannelKeeper, msg)

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		// TODO:
// 	})

// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// func HandlePacketTimeout(ctx sdk.Context, k keeper.Keeper, msg types.MsgPacketTimeout) sdk.Result {
// 	// module = lookupModule(msg.packet.sourcePort)
// 	// module.onTimeoutPacket(msg.packet)
// 	_ = channel.HandleMsgTimeoutPacket(ctx, k.ChannelKeeper, msg)

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		// TODO:
// 	})

// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// func HandlePacketTimeoutOnClose(ctx sdk.Context, k keeper.Keeper, msg types.MsgPacketTimeoutOnClose) sdk.Result {
// 	// module = lookupModule(msg.packet.sourcePort)
// 	// module.onTimeoutPacket(msg.packet)
// 	_ = channel.HandleMsgTimeoutOnClose(ctx, k.ChannelKeeper, msg)

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		// TODO:
// 	})

// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }

// func HandlePacketCleanup(ctx sdk.Context, k keeper.Keeper, msg types.MsgPacketCleanup) sdk.Result {
// 	_ = channel.HandleMsgCleanupPacket(ctx, k.ChannelKeeper, msg)

// 	ctx.EventManager().EmitEvents(sdk.Events{
// 		// TODO:
// 	})

// 	return sdk.Result{Events: ctx.EventManager().Events()}
// }
