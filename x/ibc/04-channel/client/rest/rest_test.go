package rest_test

import (
	"testing"

	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	rest_channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/client/rest"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

func TestChannelOpenInitReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)
	req := rest_channel.ChannelOpenInitReq{
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
		PortID:                "portid",
		ChannelID:             "channelid",
		Version:               "version",
		ChannelOrder:          1,
		ConnectionHops:        []string{"connectionid"},
		CounterpartyPortID:    "counterparty_port_id",
		CounterpartyChannelID: "counterparty_channel_id",
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestChannelOpenTryReq(t *testing.T) {
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

	req := rest_channel.ChannelOpenTryReq{
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
		PortID:                "portid",
		ChannelID:             "channelid",
		Version:               "version",
		ChannelOrder:          1,
		ConnectionHops:        []string{"connectionid"},
		CounterpartyPortID:    "counterparty_port_id",
		CounterpartyChannelID: "counterparty_channel_id",
		ProofInit:             proof,
		ProofHeight:           83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestChannelOpenAckReq(t *testing.T) {
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

	req := rest_channel.ChannelOpenAckReq{
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
		ProofTry:    proof,
		ProofHeight: 83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestChannelOpenConfirmReq(t *testing.T) {
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

	req := rest_channel.ChannelOpenConfirmReq{
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
		ProofAck:    proof,
		ProofHeight: 83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestChannelCloseInitReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	req := rest_channel.ChannelCloseInitReq{
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
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestChannelCloseConfirmReq(t *testing.T) {
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

	req := rest_channel.ChannelCloseConfirmReq{
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
		ProofInit:   proof,
		ProofHeight: 83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}
