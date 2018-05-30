package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/cosmos-sdk/x/stake"
)

var (
	NewProposalIDKey = []byte{0x00} //
	ProposalQueueKey = []byte{0x01} //
	ProposalListKey  = []byte{0x02}
	ProposalTypes    = []string{"TextProposal"}
)

type Keeper struct {
	// The reference to the CoinKeeper to modify balances
	ck bank.Keeper

	// The reference to the StakeMapper to get information about stakers
	sm stake.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	proposalStoreKey sdk.StoreKey

	// The wire codec for binary encoding/decoding.
	cdc *wire.Codec
}

// NewGovernanceMapper returns a mapper that uses go-wire to (binary) encode and decode gov types.
func NewKeeper(key sdk.StoreKey, ck bank.Keeper, sk stake.Keeper) Keeper {
	cdc := wire.NewCodec()
	return Keeper{
		proposalStoreKey: key,
		ck:               ck,
		cdc:              cdc,
		sm:               sk,
	}
}

// Returns the go-wire codec.
func (keeper Keeper) WireCodec() *wire.Codec {
	return keeper.cdc
}

func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID int64) *Proposal {
	store := ctx.KVStore(keeper.proposalStoreKey)
	key, _ := keeper.cdc.MarshalBinary(proposalID)
	bz := store.Get(key)
	if bz == nil {
		return nil
	}

	proposal := &Proposal{}
	err := keeper.cdc.UnmarshalBinary(bz, proposal)
	if err != nil {
		panic(err)
	}

	return proposal
}

// Implements sdk.AccountMapper.
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.proposalStoreKey)

	bz, err := keeper.cdc.MarshalBinary(proposal)
	if err != nil {
		panic(err)
	}

	key, _ := keeper.cdc.MarshalBinary(proposal.ProposalID)

	store.Set(key, bz)
}

func (keeper Keeper) getNewProposalID(ctx sdk.Context) int64 {
	store := ctx.KVStore(keeper.proposalStoreKey)
	bz := store.Get(NewProposalIDKey)

	proposalID := new(int64)
	if bz == nil {
		bz, _ = keeper.cdc.MarshalBinary(int64(0))
	}

	err := keeper.cdc.UnmarshalBinary(bz, proposalID) // TODO: switch to UnmarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic("should not happen")
	}

	bz, err = keeper.cdc.MarshalBinary(*proposalID + 1) // TODO: switch to MarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic("should not happen")
	}

	store.Set(NewProposalIDKey, bz)

	return *proposalID
}

func (keeper Keeper) getProposalQueue(ctx sdk.Context) ProposalQueue {
	store := ctx.KVStore(keeper.proposalStoreKey)
	bz := store.Get(ProposalQueueKey)
	if bz == nil {
		return nil
	}

	proposalQueue := &ProposalQueue{}
	err := keeper.cdc.UnmarshalBinary(bz, proposalQueue) // TODO: switch to UnmarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic(err)
	}

	return *proposalQueue
}

func (keeper Keeper) setProposalQueue(ctx sdk.Context, proposalQueue ProposalQueue) {
	store := ctx.KVStore(keeper.proposalStoreKey)

	bz, err := keeper.cdc.MarshalBinary(proposalQueue) // TODO: switch to MarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic(err)
	}

	store.Set(ProposalQueueKey, bz)
}

func (keeper Keeper) ProposalQueuePeek(ctx sdk.Context) *Proposal {
	proposalQueue := keeper.getProposalQueue(ctx)
	if len(proposalQueue) == 0 {
		return nil
	}
	return keeper.GetProposal(ctx, proposalQueue[0])
}

func (keeper Keeper) ProposalQueuePop(ctx sdk.Context) *Proposal {
	proposalQueue := keeper.getProposalQueue(ctx)
	if len(proposalQueue) == 0 {
		return nil
	}
	frontElement, proposalQueue := proposalQueue[0], proposalQueue[1:]
	keeper.setProposalQueue(ctx, proposalQueue)
	return keeper.GetProposal(ctx, frontElement)
}

func (keeper Keeper) ProposalQueuePush(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.proposalStoreKey)

	proposalQueue := append(keeper.getProposalQueue(ctx), proposal.ProposalID)
	bz, err := keeper.cdc.MarshalBinary(proposalQueue)
	if err != nil {
		panic(err)
	}
	store.Set(ProposalQueueKey, bz)
}

func (keeper Keeper) GetActiveProcedure() *Procedure { // TODO: move to param store and allow for updating of this
	return &Procedure{
		VotingPeriod:      200,
		MinDeposit:        sdk.Coins{{stake.StakingToken, 2}},
		ProposalTypes:     ProposalTypes,
		Threshold:         sdk.NewRat(1, 2),
		Veto:              sdk.NewRat(1, 3),
		FastPass:          sdk.NewRat(2, 3),
		MaxDepositPeriod:  200,
		GovernancePenalty: sdk.NewRat(1, 100),
	}
}

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, proposal *Proposal) {
	proposal.VotingStartBlock = ctx.BlockHeight()

	pool := keeper.sm.GetPool(ctx)
	proposal.TotalVotingPower = pool.BondedPool

	validatorList := keeper.sm.GetValidators(ctx)
	for _, validator := range validatorList {

		validatorGovInfo := ValidatorGovInfo{
			ProposalID:      proposal.ProposalID,
			ValidatorAddr:   validator.Address,
			InitVotingPower: validator.Power.Evaluate(),
			Minus:           0,
			LastVoteWeight:  -1,
		}

		proposal.ValidatorGovInfos = append(proposal.ValidatorGovInfos, validatorGovInfo)
	}

	keeper.ProposalQueuePush(ctx, *proposal)
	keeper.ProposalListDelete(ctx, *proposal)
	//fmt.Print("mikexu--Delete--")
	//fmt.Println(keeper.getProposalList(ctx))
}

func (keeper Keeper) getProposalList(ctx sdk.Context) ProposalList {
	store := ctx.KVStore(keeper.proposalStoreKey)
	bz := store.Get(ProposalListKey)
	if bz == nil {
		return nil
	}

	proposalList := &ProposalList{}
	err := keeper.cdc.UnmarshalBinary(bz, proposalList) // TODO: switch to UnmarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic(err)
	}

	return *proposalList
}

func (keeper Keeper) setProposalList(ctx sdk.Context, proposalList ProposalList) {
	store := ctx.KVStore(keeper.proposalStoreKey)

	bz, err := keeper.cdc.MarshalBinary(proposalList) // TODO: switch to MarshalBinaryBare when new go-amino gets added
	if err != nil {
		panic(err)
	}

	store.Set(ProposalListKey, bz)
}

func (keeper Keeper) ProposalListPeek(ctx sdk.Context) *Proposal {
	proposalList := keeper.getProposalList(ctx)
	if len(proposalList) == 0 {
		return nil
	}
	return keeper.GetProposal(ctx, proposalList[0])
}

func (keeper Keeper) ProposalListPop(ctx sdk.Context) *Proposal {
	proposalList := keeper.getProposalList(ctx)
	if len(proposalList) == 0 {
		return nil
	}
	frontElement, proposalList := proposalList[0], proposalList[1:]
	keeper.setProposalList(ctx, proposalList)
	return keeper.GetProposal(ctx, frontElement)
}


func (keeper Keeper) ProposalListAppend(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.proposalStoreKey)

	proposalList := append(keeper.getProposalList(ctx), proposal.ProposalID)
	bz, err := keeper.cdc.MarshalBinary(proposalList)
	if err != nil {
		panic(err)
	}
	store.Set(ProposalListKey, bz)
}

func (keeper Keeper) ProposalListDelete(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.proposalStoreKey)
	proposalList := keeper.getProposalList(ctx)

	for index,proposalID :=range keeper.getProposalList(ctx){
		if proposalID == proposal.ProposalID{
			proposalList = append(proposalList[:index],proposalList[index+1:]...)
			index = index
			break
		}

	}
	bz, err := keeper.cdc.MarshalBinary(proposalList)
	if err != nil {
		panic(err)
	}
	store.Set(ProposalListKey, bz)

}