package mymodule

import (
"encoding/json"

sdk "github.com/cosmos/cosmos-sdk/types"
)

// name to idetify transaction types
const MsgType = "mymodule"

// XXX remove: think it makes more sense belonging with the Params so we can
// initialize at genesis - to allow for the same tests we should should make
// the ValidateBasic() function a return from an initializable function
// ValidateBasic(bondDenom string) function
const mymoduleToken = "steak"

//Verify interface at compile time
var _, _ sdk.Msg = &MsgDo{}, &MsgUndo{}

//var msgCdc = wire.NewCodec()
//
//func init() {
//	wire.RegisterCrypto(msgCdc)
//}

//______________________________________________________________________

type MsgDo struct {
	Addr sdk.Address   `json:"address"`
	ValueNum ValueNum
}

func NewMsgDo(addr sdk.Address,valueNum ValueNum) MsgDo {
	return MsgDo{
           Addr :addr,
		   ValueNum: valueNum,
	}
}

//nolint
func (msg MsgDo) Type() string              { return MsgType } //TODO update "stake/declarecandidacy"
func (msg MsgDo) GetSigners() []sdk.Address { return []sdk.Address{msg.Addr} }

// get the bytes for the message signer to sign on
func (msg MsgDo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgDo) ValidateBasic() sdk.Error {
	if msg.ValueNum.Num == 0 {
		return ErrValueNumEmpty(DefaultCodespace)
	}
	if msg.Addr == nil {
		return ErrAddrEmpty(DefaultCodespace)
	}
	return nil
}

// Implements Msg.
func (msg MsgDo) Get(key interface{}) (value interface{}) {
	return nil
}
//______________________________________________________________________

// MsgEditCandidacy - struct for editing a candidate
type MsgUndo struct {
	Addr sdk.Address `json:"address"`
}

func NewMsgUndo(addr sdk.Address) MsgUndo {
	return MsgUndo{
		Addr:   addr,
	}
}

//nolint
func (msg MsgUndo) Type() string              { return MsgType } //TODO update "stake/msgeditcandidacy"
func (msg MsgUndo) GetSigners() []sdk.Address { return []sdk.Address{msg.Addr} }

// get the bytes for the message signer to sign on
func (msg MsgUndo) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgUndo) ValidateBasic() sdk.Error {
	if msg.Addr == nil {
		return ErrAddrEmpty(DefaultCodespace)
	}
	return nil
}

func (msg MsgUndo) Get(key interface{}) (value interface{}) {
	return nil
}