package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

// HandleMsgRecvPacket defines the sdk.Handler for MsgRecvPacket
func HandleMsgRecvPacket(ctx sdk.Context, k Keeper, msg types.MsgPacketRecv, acknowledgement []byte, portCapability sdk.CapabilityKey) (res sdk.Result) {
	_, err := k.RecvPacket(ctx, msg.Packet, msg.Proofs[0], msg.Height, acknowledgement, portCapability)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
