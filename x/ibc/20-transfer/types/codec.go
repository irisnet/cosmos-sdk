package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var SubModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "ibc/transfer/MsgTransfer", nil)
	cdc.RegisterConcrete(PacketData{}, "ibc/transfer/PacketData", nil)

	SubModuleCdc = cdc
}
