package rest_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	rest_client "github.com/cosmos/cosmos-sdk/x/ibc/02-client/client/rest"
	tendermint "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types/tendermint"
)

func TestCreateClientReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	req := rest_client.CreateClientReq{
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
		ClientID:       "client_id",
		ConsensusState: tendermint.ConsensusState{},
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestUpdateClientReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	req := rest_client.UpdateClientReq{
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
		Header: tendermint.Header{},
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}

func TestSubmitMisbehaviourReq(t *testing.T) {
	cdc := codec.New()
	ibc.AppModuleBasic{}.RegisterCodec(cdc)

	req := rest_client.SubmitMisbehaviourReq{
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
		Evidence: tendermint.Evidence{},
	}

	bz, _ := cdc.MarshalJSON(req)
	println(string(bz))
}
