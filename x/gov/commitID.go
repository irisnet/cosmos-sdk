package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var(
	KVstoreKeyListKey = []byte("k/")
	)

// Get Proposal from store by ProposalID

func (keeper Keeper) GetKVstoreKeylist(ctx sdk.Context, proposalID int64) string {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KVstoreKeyListKey)
	if bz == nil {
		return " "
	}
	KVstoreKeylist := string(bz)
	return KVstoreKeylist
}

// Implements sdk.AccountMapper.
func (keeper Keeper) setKVstoreKeylist(ctx sdk.Context, KVstoreKeyList string) {
	store := ctx.KVStore(keeper.storeKey)
	bz := []byte(KVstoreKeyList)
	store.Set(KVstoreKeyListKey, bz)
}

// InitGenesis - store genesis parameters
func InitGenesis_commitID(ctx sdk.Context, k Keeper, KVstoreKeyList string) {
	k.setKVstoreKeylist(ctx,KVstoreKeyList)
}