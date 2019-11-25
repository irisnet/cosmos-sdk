package rest_test

import (
	"testing"

	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	channel_types "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	rest_transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/client/rest"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

func TestTransferTxReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	address, _ := sdk.AccAddressFromBech32("cosmos15keln30v6n7p6pyvhfw6y4wth2h9vmmse4lmce")

	req := rest_transfer.TransferTxReq{
		BaseReq: rest.BaseReq{
			From:          "cosmos15keln30v6n7p6pyvhfw6y4wth2h9vmmse4lmce",
			Memo:          "memo",
			ChainID:       "chann-2",
			AccountNumber: 3,
			Sequence:      2,
			Gas:           "200000",
			GasAdjustment: "1.2",
			Fees: []types.Coin{
				{
					Denom:  "stake",
					Amount: types.NewInt(50),
				},
			},
			Simulate: false,
		},
		Amount: []types.Coin{
			{
				Denom:  "stake",
				Amount: types.NewInt(50),
			},
		},
		Receiver: address,
		Source:   false,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestRecvPacketReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	proof := commitment.Proof{
		Proof: &merkle.Proof{
			Ops: []merkle.ProofOp{
				{
					Type: "iavl:v",
					Key:  []byte("test_key"),
					Data: []byte("test_data"),
					//
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     []byte(""),
					XXX_sizecache:        1,
				},
				{
					Type: "multistore",
					Key:  []byte("test_key"),
					Data: []byte("test_data"),
					//
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     []byte("test"),
					XXX_sizecache:        1,
				},
			},
			//
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     []byte("test"),
			XXX_sizecache:        1,
		},
	}

	req := rest_transfer.RecvPacketReq{
		BaseReq: rest.BaseReq{
			From:          "cosmos15keln30v6n7p6pyvhfw6y4wth2h9vmmse4lmce",
			Memo:          "memo",
			ChainID:       "chann-2",
			AccountNumber: 3,
			Sequence:      2,
			Gas:           "200000",
			GasAdjustment: "1.2",
			Fees: []types.Coin{
				{
					Denom:  "stake",
					Amount: types.NewInt(50),
				},
			},
			Simulate: false,
		},
		Packet: channel_types.Packet{
			Sequence:           100,
			Timeout:            100,
			SourcePort:         "source_port",
			SourceChannel:      "source_channel",
			DestinationPort:    "",
			DestinationChannel: "",
			Data:               []byte("data.."),
		},
		Proofs: []commitment.Proof{proof},
		Height: 83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}
