package gov

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestParameterProposal(t *testing.T) {
	mapp, keeper, _, _, _, _ := getMockApp(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{})

	fmt.Println(keeper.GetDepositProcedure(ctx))
	fmt.Println(keeper.GetTallyingProcedure(ctx))
	fmt.Println(keeper.GetVotingProcedure(ctx))

	pp := ParameterProposal{
		Params: []Param{
			{Key: Prefix + ParamStoreKeyDepositProcedureDeposit, Value: "200iris", Op: Update},
			{Key: Prefix + ParamStoreKeyDepositProcedureMaxDepositPeriod, Value: "20", Op: Update},
			{Key: Prefix + ParamStoreKeyTallyingProcedurePenalty, Value: "1/50", Op: Update},
			{Key: Prefix + ParamStoreKeyTallyingProcedureVeto, Value: "1/4", Op: Update},
			{Key: Prefix + ParamStoreKeyTallyingProcedureThreshold, Value: "2/8", Op: Update},
		},
	}

	pp.Execute(ctx, keeper)
	assert.Equal(t, keeper.GetDepositProcedure(ctx).MinDeposit,
		sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(200))})

	assert.Equal(t, keeper.GetDepositProcedure(ctx).MaxDepositPeriod,int64(20))




	assert.Equal(t, keeper.GetTallyingProcedure(ctx),
		TallyingProcedure{
			Threshold:         sdk.NewRat(2, 8),
			Veto:              sdk.NewRat(1, 4),
			GovernancePenalty: sdk.NewRat(1, 50),
		})
}