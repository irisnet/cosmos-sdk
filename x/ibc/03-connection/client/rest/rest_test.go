package rest_test

import (
	"testing"

	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	rest_connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/client/rest"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
)

func TestConnectionOpenInitReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)
	req := rest_connection.ConnectionOpenInitReq{
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
		ConnectionID:             "conntochaintwo",
		ClientID:                 "conntochaintwo",
		CounterpartyClientID:     "conntochaintwo",
		CounterpartyConnectionID: "conntochaintwo",
		CounterpartyPrefix:       commitment.NewPrefix([]byte("ibc")),
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestConnectionOpenTryReq(t *testing.T) {
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

	req := rest_connection.ConnectionOpenTryReq{
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
		ConnectionID:             "conntochaintwo",
		ClientID:                 "conntochaintwo",
		CounterpartyClientID:     "conntochaintwo",
		CounterpartyConnectionID: "conntochaintwo",
		CounterpartyPrefix:       commitment.NewPrefix([]byte("ibc")),
		CounterpartyVersions:     []string{"1.0.0"},
		ProofInit:                proof,
		ProofConsensus:           proof,
		ProofHeight:              83,
		ConsensusHeight:          83,
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestConnectionOpenAckReq(t *testing.T) {
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

	req := rest_connection.ConnectionOpenAckReq{
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
		ProofTry:        proof,
		ProofConsensus:  proof,
		ProofHeight:     83,
		ConsensusHeight: 83,
		Version:         "1.0.0",
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestConnectionOpenConfirmReq(t *testing.T) {
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

	req := rest_connection.ConnectionOpenConfirmReq{
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
