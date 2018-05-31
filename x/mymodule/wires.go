package mymodule

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgDo{}, "cosmos-sdk/MsgDo", nil)
	cdc.RegisterConcrete(MsgUndo{}, "cosmos-sdk/MsgUndo", nil)
}