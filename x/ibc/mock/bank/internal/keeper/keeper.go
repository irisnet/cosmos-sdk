package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	ics04 "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ics23 "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"github.com/cosmos/cosmos-sdk/x/ibc/mock/bank/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	cdc  *codec.Codec
	key  sdk.StoreKey
	ibck ibc.Keeper
	bk   types.BankKeeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ibck ibc.Keeper, bk types.BankKeeper) Keeper {
	return Keeper{
		cdc:  cdc,
		key:  key,
		ibck: ibck,
		bk:   bk,
	}
}

// SendTransfer handles transfer sending logic
func (k Keeper) SendTransfer(ctx sdk.Context, srcPort, srcChan string, denom string, amount sdk.Int, sender sdk.AccAddress, receiver string, source bool) sdk.Error {
	// get the port and channel of the counterparty
	channel, ok := k.ibck.ChannelKeeper.GetChannel(ctx, srcPort, srcChan)
	if !ok {
		return sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), ics04.CodeChannelNotFound, "failed to get channel")
	}

	dstPort := channel.Counterparty.PortID
	dstChan := channel.Counterparty.ChannelID

	// get the next sequence
	sequence, ok := k.ibck.ChannelKeeper.GetNextSequenceSend(ctx, srcPort, srcChan)
	if !ok {
		return sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), ics04.CodeSequenceNotFound, "failed to retrieve sequence")
	}

	return k.createOutgoingPacket(ctx, sequence, srcPort, srcChan, dstPort, dstChan, denom, amount, sender, receiver, source)
}

// ReceiveTransfer handles transfer receiving logic
func (k Keeper) ReceiveTransfer(ctx sdk.Context, packet exported.PacketI, proofs ics23.Proof, height uint64) sdk.Error {
	_, err := k.ibck.ChannelKeeper.RecvPacket(ctx, packet, proofs, height, nil)
	if err != nil {
		return sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), types.CodeErrReceivePacket, "failed to receive packet")
	}

	var data types.TransferPacketData
	err = data.UnmarshalJSON(packet.Data)
	if err != nil {
		return sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), types.CodeInvalidPacketData, "invalid packet data")
	}

	receiverAddr, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), types.CodeInvalidReceiver, "invalid receiver address")
	}

	if data.Source {
		// mint tokens
		_, err := k.bk.AddCoins(ctx, receiverAddr, sdk.Coins{sdk.NewCoin(data.Denomination, data.Amount)})
		if err != nil {
			return err
		}

	} else {
		// unescrow tokens
		escrowAddress := k.GetEscrowAddress(packet.DestChannel())
		err := k.bk.SendCoins(ctx, escrowAddress, receiverAddr, sdk.Coins{sdk.NewCoin(data.Denomination, data.Amount)})
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) createOutgoingPacket(ctx sdk.Context, seq uint64, srcPort, srcChan, dstPort, dstChan string, denom string, amount sdk.Int, sender sdk.AccAddress, receiver string, source bool) sdk.Error {
	if source {
		// escrow tokens
		escrowAddress := k.GetEscrowAddress(srcChan)
		err := k.bk.SendCoins(ctx, sender, escrowAddress, sdk.Coins{sdk.NewCoin(denom, amount)})
		if err != nil {
			return err
		}

	} else {
		// burn vouchers from sender
		err := k.bk.BurnCoins(ctx, sender, sdk.Coins{sdk.NewCoin(denom, amount)})
		if err != nil {
			return err
		}
	}

	// build packet
	packetData := types.TransferPacketData{
		Denomination: denom,
		Amount:       amount,
		Sender:       sender,
		Receiver:     receiver,
		Source:       source,
	}

	packetDataBz, err := packetData.MarshalJSON()
	if err != nil {
		sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), types.CodeInvalidPacketData, "invalid packet data")
	}

	packet := ics04.NewPacket(seq, 0, srcPort, srcChan, dstPort, dstChan, packetDataBz)

	err = k.ibck.ChannelKeeper.SendPacket(ctx, packet)
	if err != nil {
		return sdk.NewError(sdk.CodespaceType(types.DefaultCodespace), types.CodeErrSendPacket, "failed to send packet")
	}

	return nil
}

// GetEscrowAddress returns the escrow address for the specified channel
func (k Keeper) GetEscrowAddress(chanID string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(chanID)))
}
