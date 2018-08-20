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
		sdk.Coins{sdk.NewCoin("iris", 200)})

	assert.Equal(t, keeper.GetDepositProcedure(ctx).MaxDepositPeriod,int64(20))




	assert.Equal(t, keeper.GetTallyingProcedure(ctx),
		TallyingProcedure{
			Threshold:         sdk.NewRat(2, 8),
			Veto:              sdk.NewRat(1, 4),
			GovernancePenalty: sdk.NewRat(1, 50),
		})
}

func TestSoftwareUpgradeProposal(t *testing.T) {
	mapp, keeper,pk, _, _, _,_ := getMockAppPK(t, 0)

	mapp.BeginBlock(abci.RequestBeginBlock{})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{Height:0})

	var proposalID int64
	var height int64
	sp := SoftwareUpgradeProposal{TextProposal{ProposalID:28}}

	pk.Setter().GovSetter().Set(ctx,"upgrade/proposalId",int64(-1))
	sp.Execute(ctx, keeper)
	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalId",&proposalID)
	assert.Equal(t,proposalID,int64(28))


	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalAcceptHeight",&height)
	assert.Equal(t,height,int64(0))


	pk.Setter().GovSetter().Set(ctx,"upgrade/proposalId",int64(5))
	pk.Setter().GovSetter().Set(ctx,"upgrade/proposalAcceptHeight",int64(-1))

	sp.Execute(ctx, keeper)
	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalId",&proposalID)
	assert.Equal(t,proposalID,int64(5))
	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalAcceptHeight",&height)
	assert.Equal(t,height,int64(-1))

	ctx = mapp.BaseApp.NewContext(false, abci.Header{Height:64})
	pk.Setter().GovSetter().Set(ctx,"upgrade/proposalId",int64(-1))
	pk.Setter().GovSetter().Set(ctx,"upgrade/proposalAcceptHeight",int64(-1))

	sp.Execute(ctx, keeper)
	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalId",&proposalID)
	assert.Equal(t,proposalID,int64(28))
	pk.Getter().GovGetter().Get(ctx,"upgrade/proposalAcceptHeight",&height)
	assert.Equal(t,height,int64(64))

	}