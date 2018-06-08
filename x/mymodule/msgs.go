package mymodule

import (
"encoding/json"

sdk "github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/wire"
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

var msgCdc = wire.NewCodec()

func init() {
	wire.RegisterCrypto(msgCdc)
}

//______________________________________________________________________

// MsgDeclareCandidacy - struct for unbonding transactions
type MsgDo struct {
	addr sdk.Address   `json:"address"`
	valueNum ValueNum
}

func NewMsgDo(addr sdk.Address,valueNum ValueNum) MsgDo {
	return MsgDo{
           addr :addr,
		   valueNum: valueNum,
	}
}

//nolint
func (msg MsgDo) Type() string              { return MsgType } //TODO update "stake/declarecandidacy"
func (msg MsgDo) GetSigners() []sdk.Address { return []sdk.Address{msg.addr} }

// get the bytes for the message signer to sign on
func (msg MsgDo) GetSignBytes() []byte {
	return msgCdc.MustMarshalBinary(msg)
}

// quick validity check
func (msg MsgDo) ValidateBasic() sdk.Error {
	if msg.valueNum.num == 0 {
		return ErrValueNumEmpty(DefaultCodespace)
	}
	if msg.addr == nil {
		return ErrAddrEmpty(DefaultCodespace)
	}
	return nil
}

//______________________________________________________________________

// MsgEditCandidacy - struct for editing a candidate
type MsgUndo struct {
	addr sdk.Address `json:"address"`
}

func NewMsgUndo(addr sdk.Address) MsgUndo {
	return MsgUndo{
		addr:   addr,
	}
}

//nolint
func (msg MsgUndo) Type() string              { return MsgType } //TODO update "stake/msgeditcandidacy"
func (msg MsgUndo) GetSigners() []sdk.Address { return []sdk.Address{msg.addr} }

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
	if msg.addr == nil {
		return ErrAddrEmpty(DefaultCodespace)
	}
	return nil
}
