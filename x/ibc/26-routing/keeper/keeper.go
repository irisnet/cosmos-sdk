package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/types"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	codespace sdk.CodespaceType
	prefix    []byte

	ClientKeeper     client.Keeper
	ConnectionKeeper connection.Keeper
	ChannelKeeper    channel.Keeper
	PortKeeper       port.Keeper
	TransferKeeper   transfer.Keeper
}

func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	codespace sdk.CodespaceType,
	bk transfer.BankKeeper,
	sk transfer.SupplyKeeper,
) Keeper {
	clientKeeper := client.NewKeeper(cdc, key, codespace)
	connectionKeeper := connection.NewKeeper(cdc, key, codespace, clientKeeper)
	portKeeper := port.NewKeeper(cdc, key, codespace)
	channelKeeper := channel.NewKeeper(cdc, key, codespace, clientKeeper, connectionKeeper, portKeeper)
	transferKeeper := transfer.NewKeeper(cdc, key, codespace, clientKeeper, connectionKeeper, channelKeeper, bk, sk)

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		codespace:        sdk.CodespaceType(fmt.Sprintf("%s/%s", codespace, types.DefaultCodespace)), // "ibc/routing",
		prefix:           []byte(types.SubModuleName + "/"),                                          // "routing/"
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
		TransferKeeper:   transferKeeper,
	}
}

func (k Keeper) setCapability(ctx sdk.Context, portID string, capability sdk.CapabilityKey) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(capability)
	store.Set(types.KeyAuthenticationPath(portID), bz)
}

func (k Keeper) deleteCapability(ctx sdk.Context, portID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	store.Delete(types.KeyAuthenticationPath(portID))
}

func (k Keeper) getCapability(ctx sdk.Context, portID string) (sdk.CapabilityKey, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := store.Get(types.KeyAuthenticationPath(portID))
	if bz == nil {
		return sdk.NewKVStoreKey(""), false
	}
	var capability sdk.CapabilityKey
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &capability)
	return capability, true
}

func (k Keeper) setCallbacks(ctx sdk.Context, portID string, callbacks types.ModuleCallbacks) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(callbacks)
	store.Set(types.KeyCallbackPath(portID), bz)
}

func (k Keeper) deleteCallbacks(ctx sdk.Context, portID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	store.Delete(types.KeyCallbackPath(portID))
}

func (k Keeper) getCallbacks(ctx sdk.Context, portID string) (types.ModuleCallbacks, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), k.prefix)
	bz := store.Get(types.KeyCallbackPath(portID))
	if bz == nil {
		return nil, false
	}
	var callbacks types.ModuleCallbacks
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &callbacks)
	return callbacks, true
}
