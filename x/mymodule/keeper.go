package mymodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"

)

var (
	countKey = []byte{0x00} //
	//Types    = []string{"TextProposal"}
)

type Keeper struct {
	// The reference to the CoinKeeper to modify balances
	ck bank.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	mymoduleStoreKey sdk.StoreKey

	// The wire codec for binary encoding/decoding.
	cdc *wire.Codec

	//
	codespace sdk.CodespaceType
}

// NewGovernanceMapper returns a mapper that uses go-wire to (binary) encode and decode gov types.
func NewKeeper(cdc *wire.Codec,key sdk.StoreKey, ck bank.Keeper,codespace sdk.CodespaceType) Keeper {
	return Keeper{
		mymoduleStoreKey: key,
		ck:               ck,
		cdc:              cdc,
		codespace: codespace,
	}
}

// Returns the go-wire codec.
func (keeper Keeper) WireCodec() *wire.Codec {
	return keeper.cdc
}

func (keeper Keeper) GetCounter(ctx sdk.Context, addr sdk.Address) int64{
	store := ctx.KVStore(keeper.mymoduleStoreKey)
	b := store.Get(addr.Bytes())
	if b== nil{
		return 0
	}
	var counter int64
	err := keeper.cdc.UnmarshalBinary(b,&counter)
	if err != nil{
		panic(err)
	}
	return  counter
}

func (keeper Keeper) SetCounter(ctx sdk.Context, addr sdk.Address,counter int64){
	store := ctx.KVStore(keeper.mymoduleStoreKey)
	bz,err := keeper.cdc.MarshalBinary(counter)

	if err != nil{
		panic(err)
	}
	store.Set(addr.Bytes(),bz)
}