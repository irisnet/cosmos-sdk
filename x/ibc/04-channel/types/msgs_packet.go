package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

var _ sdk.Msg = MsgPacketRecv{}

type MsgPacketRecv struct {
	Packet exported.PacketI   `json:"packet" yaml:"packet"`
	Proofs []commitment.Proof `json:"proofs" yaml:"proofs"`
	Height uint64             `json:"height" yaml:"height"`
	Signer sdk.AccAddress     `json:"signer" yaml:"signer"`
}

// NewMsgRecvPacket creates a new MsgRecvPacket instance
func NewMsgRecvPacket(packet exported.PacketI, proofs []commitment.Proof, height uint64, signer sdk.AccAddress) MsgPacketRecv {
	return MsgPacketRecv{
		Packet: packet,
		Proofs: proofs,
		Height: height,
		Signer: signer,
	}
}

// Route implements sdk.Msg
func (MsgPacketRecv) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgPacketRecv) Type() string {
	return "recv_packet"
}

// ValidateBasic implements sdk.Msg
func (msg MsgPacketRecv) ValidateBasic() sdk.Error {
	if msg.Height < 1 {
		return sdk.NewError(DefaultCodespace, CodeInvalidHeight, "invalid height")
	}

	if msg.Proofs == nil {
		return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
	}

	for _, proof := range msg.Proofs {
		if proof.Proof == nil {
			return sdk.NewError(DefaultCodespace, CodeProofMissing, "proof missing")
		}
	}

	if msg.Signer.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidAddress, "invalid signer")
	}

	return msg.Packet.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (msg MsgPacketRecv) GetSignBytes() []byte {
	return sdk.MustSortJSON(SubModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgPacketRecv) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
