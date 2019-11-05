package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
)

type ModuleCallbacks interface {
	OnChanOpenInit(
		ctx sdk.Context,
		order types.Order,
		connectionHops []string,
		portID,
		channelID string,
		counterparty types.Counterparty,
		version string,
	) error

	OnChanOpenTry(
		ctx sdk.Context,
		order types.Order,
		connectionHops []string,
		portID,
		channelID string,
		counterparty types.Counterparty,
		version string,
		counterpartyVersion string,
	) error

	OnChanOpenAck(
		ctx sdk.Context,
		portID,
		channelID string,
		version string,
	) error

	OnChanOpenConfirm(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnChanCloseInit(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnChanCloseConfirm(
		ctx sdk.Context,
		portID,
		channelID string,
	) error

	OnRecvPacket(
		ctx sdk.Context,
		packet types.Packet,
	) ([]byte, error)

	OnTimeoutPacket(
		ctx sdk.Context,
		packet types.Packet,
	) error

	OnAcknowledgePacket(
		ctx sdk.Context,
		packet types.Packet,
		acknowledgement []byte,
	) error

	OnTimeoutPacketClose(
		ctx sdk.Context,
		packet types.Packet,
	)
}
