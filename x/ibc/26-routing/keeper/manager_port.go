package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/26-routing/types"
)

func (k Keeper) BindPort(ctx sdk.Context, portID string, callbacks types.ModuleCallbacks) error {
	_, existed := k.LookupModule(ctx, portID)
	if existed {
		return types.ErrConnectionExists(types.DefaultCodespace, portID)
	}
	capability := k.PortKeeper.BindPort(portID)
	k.setCapability(ctx, portID, capability)
	k.setCallbacks(ctx, portID, callbacks)
	return nil
}

func (k Keeper) UpdatePort(ctx sdk.Context, portID string, callbacks types.ModuleCallbacks) error {
	key, found := k.getCapability(ctx, portID)
	if !found {
		return errors.New("key not found") // TODO:
	}
	if !k.PortKeeper.Authenticate(key, portID) {
		return errors.New("authenticate failed") // TODO:
	}
	k.setCallbacks(ctx, portID, callbacks)
	return nil
}

func (k Keeper) ReleasePort(ctx sdk.Context, portID string) error {
	_, found := k.getCapability(ctx, portID)
	if !found {
		return errors.New("key not found") // TODO:
	}
	k.PortKeeper.ReleasePort(portID)
	k.deleteCapability(ctx, portID)
	k.deleteCallbacks(ctx, portID)
	return nil
}

func (k Keeper) LookupModule(ctx sdk.Context, portID string) (types.ModuleCallbacks, bool) {
	callbacks, found := k.getCallbacks(ctx, portID)
	if !found {
		return nil, false
	}
	return callbacks, true
}
