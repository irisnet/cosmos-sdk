package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
)

var SubModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.PacketI)(nil), nil)
	cdc.RegisterConcrete(Packet{}, "ibc/channel/Packet", nil)
	cdc.RegisterConcrete(OpaquePacket{}, "ibc/channel/OpaquePacket", nil)

	cdc.RegisterConcrete(MsgChannelOpenInit{}, "ibc/channel/MsgChannelOpenInit", nil)
	cdc.RegisterConcrete(MsgChannelOpenTry{}, "ibc/channel/MsgChannelOpenTry", nil)
	cdc.RegisterConcrete(MsgChannelOpenAck{}, "ibc/channel/MsgChannelOpenAck", nil)
	cdc.RegisterConcrete(MsgChannelOpenConfirm{}, "ibc/channel/MsgChannelOpenConfirm", nil)
	cdc.RegisterConcrete(MsgChannelCloseInit{}, "ibc/channel/MsgChannelCloseInit", nil)
	cdc.RegisterConcrete(MsgChannelCloseConfirm{}, "ibc/channel/MsgChannelCloseConfirm", nil)

	cdc.RegisterConcrete(MsgPacketRecv{}, "ibc/channel/MsgPacketRecv", nil)
	// cdc.RegisterConcrete(MsgPacketAcknowledgement{}, "ibc/channel/MsgPacketAcknowledgement", nil)
	// cdc.RegisterConcrete(MsgPacketTimeout{}, "ibc/channel/MsgPacketTimeout", nil)
	// cdc.RegisterConcrete(MsgPacketTimeoutOnClose{}, "ibc/channel/MsgPacketTimeoutOnClose", nil)
	// cdc.RegisterConcrete(MsgPacketCleanup{}, "ibc/channel/MsgPacketCleanup", nil)

	SubModuleCdc = cdc
}
