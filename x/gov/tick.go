package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func NewBeginBlocker(keeper Keeper) sdk.BeginBlocker {
	return func(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		checkProposal(ctx,keeper)
		proposal := keeper.getExecuteQueue(ctx).Pop()
		if proposal == nil {
			return abci.ResponseBeginBlock{}
		}

		ctx.Logger().Info("execute Proposal", "Proposal", proposal.ProposalID)
		refund(ctx, proposal, keeper)

		if proposal.isTimeOutPass(ctx.BlockHeight()){
			slash(ctx, proposal)
		}

		//TODO proposal.execute
		return abci.ResponseBeginBlock{}
	}
}

// refund Deposit
func refund(ctx sdk.Context, proposal *Proposal, keeper Keeper) {
	for _, deposit := range proposal.Deposits {
		ctx.Logger().Info("Execute Refund", "Depositer", deposit.Depositer, "Amount", deposit.Amount)
		_, _, err := keeper.ck.AddCoins(ctx, deposit.Depositer, deposit.Amount)
		if err != nil {
			panic("should not happen")
		}
	}
}

// Slash validators if not voted
func slash(ctx sdk.Context, proposal *Proposal) {
	ctx.Logger().Info("Begin to Execute Slash")
	for _, validatorGovInfo := range proposal.ValidatorGovInfos {
		if validatorGovInfo.LastVoteWeight < 0 {
			// TODO: SLASH MWAHAHAHAHAHA
			ctx.Logger().Info("Execute Slash", "validator", validatorGovInfo.ValidatorAddr,"ProposalId",proposal.ProposalID)
		}
	}
}

//check Deposit timeout
func checkProposal(ctx sdk.Context,keeper Keeper){
	proposals := keeper.popExpiredProposal(ctx)
	for _,proposal := range proposals {
		if proposal.isActive(){
			if proposal.isTimeOutPass(ctx.BlockHeight()) {
				keeper.getExecuteQueue(ctx).Push(proposal.ProposalID)
			}else {
				refund(ctx, proposal, keeper)
				slash(ctx, proposal)
			}
		}else {
			refund(ctx, proposal, keeper)
		}
	}
}
