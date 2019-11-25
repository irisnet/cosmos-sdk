package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// SubModuleCdc defines the IBC transfer codec.
var SubModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "ibc/transfer/MsgTransfer", nil)
	cdc.RegisterConcrete(MsgRecvPacket{}, "ibc/transfer/MsgRecvPacket", nil)
	cdc.RegisterConcrete(PacketData{}, "ibc/transfer/PacketData", nil)

	SetSubModuleCodec(cdc)
}

func SetSubModuleCodec(cdc *codec.Codec) {
	SubModuleCdc = cdc
}
