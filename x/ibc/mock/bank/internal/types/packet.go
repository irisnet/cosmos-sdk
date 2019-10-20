package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
)

var _ exported.PacketI = Packet{}

// Packet defines a type that carries data across different chains through IBC
type Packet struct {
	Msequence           uint64 // number corresponds to the order of sends and receives, where a packet with an earlier sequence number must be sent and received before a packet with a later sequence number.
	Mtimeout            uint64 // indicates a consensus height on the destination chain after which the packet will no longer be processed, and will instead count as having timed-out.
	MsourcePort         string // identifies the port on the sending chain.
	MsourceChannel      string // identifies the channel end on the sending chain.
	MdestinationPort    string // identifies the port on the receiving chain.
	MdestinationChannel string // identifies the channel end on the receiving chain.
	Mdata               []byte // opaque value which can be defined by the application logic of the associated modules.
}

// newPacket creates a new Packet instance
func NewPacket(
	sequence, timeout uint64, sourcePort, sourceChannel,
	destinationPort, destinationChannel string, data []byte,
) Packet {
	return Packet{
		sequence,
		timeout,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		data,
	}
}

// Sequence implements PacketI interface
func (p Packet) Sequence() uint64 { return p.Msequence }

// TimeoutHeight implements PacketI interface
func (p Packet) TimeoutHeight() uint64 { return p.Mtimeout }

// SourcePort implements PacketI interface
func (p Packet) SourcePort() string { return p.MsourcePort }

// SourceChannel implements PacketI interface
func (p Packet) SourceChannel() string { return p.MsourceChannel }

// DestPort implements PacketI interface
func (p Packet) DestPort() string { return p.MdestinationPort }

// DestChannel implements PacketI interface
func (p Packet) DestChannel() string { return p.MdestinationChannel }

// Data implements PacketI interface
func (p Packet) Data() []byte { return p.Mdata }

type PacketAlias Packet

func (p Packet) MarshalJSON() ([]byte, error) {
	return MouduleCdc.MarshalJSON(PacketAlias(p))
}

func (p *Packet) UnmarshalJSON(bz []byte) (err error) {
	return MouduleCdc.UnmarshalJSON(bz, (*PacketAlias)(p))
}

// TransferPacketData defines a struct for the packet payload
type TransferPacketData struct {
	Denomination string         `json:"denomination" yaml:"denomination"`
	Amount       sdk.Int        `json:"amount" yaml:"amount"`
	Sender       sdk.AccAddress `json:"sender" yaml:"sender"`
	Receiver     string         `json:"receiver" yaml:"receiver"`
	Source       bool           `json:"source" yaml:"source"`
}

func (tpd TransferPacketData) String() string {
	return fmt.Sprintf(`TransferPacketData:
	Denomination          %s
	Amount:               %s
	Sender:               %s
	Receiver:             %s
	Source:               %v`,
		tpd.Denomination,
		tpd.Amount.String(),
		tpd.Sender.String(),
		tpd.Receiver,
		tpd.Source,
	)
}

func (tpd TransferPacketData) Validate() error {
	if !tpd.Amount.IsPositive() {
		return sdk.NewError(sdk.CodespaceType(DefaultCodespace), CodeInvalidAmount, "invalid amount")
	}

	if tpd.Sender.Empty() || len(tpd.Receiver) == 0 {
		return sdk.NewError(sdk.CodespaceType(DefaultCodespace), CodeInvalidAddress, "invalid address")
	}

	return nil
}
