package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

//__BEACON_ID___________________________________________________________

// GetHighestBeaconID gets the highest BEACON ID
func (k Keeper) GetHighestBeaconID(ctx sdk.Context) (beaconID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestBeaconIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrapf(types.ErrInvalidGenesis, "initial beacon ID hasn't been set")
	}
	// convert from bytes to uint64
	beaconID = types.GetBeaconIDFromBytes(bz)
	return beaconID, nil
}

// SetHighestBeaconID sets the new highest BEACON ID to the store
func (k Keeper) SetHighestBeaconID(ctx sdk.Context, beaconID uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	beaconIDbz := types.GetBeaconIDBytes(beaconID)
	store.Set(types.HighestBeaconIDKey, beaconIDbz)
}

//__BEACONS_____________________________________________________________

// SetBeacon Sets the BEACON metadata struct for a beaconID
func (k Keeper) SetBeacon(ctx sdk.Context, beacon types.Beacon) error {
	// must have an owner
	if beacon.Owner.Empty() {
		return sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have owner")
	}

	//must have an ID
	if beacon.BeaconID == 0 {
		return sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - id must be positive non-zero")
	}

	//must have a moniker
	if len(beacon.Moniker) == 0 {
		return sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have a moniker")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.BeaconKey(beacon.BeaconID), k.cdc.MustMarshalBinaryLengthPrefixed(beacon))

	return nil
}

// SetLastTimestampID - sets the last timestamp ID submitted
func (k Keeper) SetLastTimestampID(ctx sdk.Context, beaconID uint64, timestampID uint64) error {
	beacon := k.GetBeacon(ctx, beaconID)
	if beacon.Owner.Empty() {
		// doesn't exist. Don't update
		return types.ErrBeaconDoesNotExist
	}
	if timestampID > beacon.LastTimestampID {
		beacon.LastTimestampID = timestampID
		return k.SetBeacon(ctx, beacon)
	}
	return nil
}

// GetBeacon Gets the entire BEACON metadata struct for a beaconID
func (k Keeper) GetBeacon(ctx sdk.Context, beaconID uint64) types.Beacon {
	store := ctx.KVStore(k.storeKey)
	if !k.IsBeaconRegistered(ctx, beaconID) {
		// return a new empty Beacon struct
		return types.NewBeacon()
	}
	bz := store.Get(types.BeaconKey(beaconID))
	var beacon types.Beacon
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &beacon)
	return beacon
}

// GetBeaconOwner - get the current owner of a BEACON
func (k Keeper) GetBeaconOwner(ctx sdk.Context, beaconID uint64) sdk.AccAddress {
	return k.GetBeacon(ctx, beaconID).Owner
}

// IsBeaconRegistered Checks if the BEACON is present in the store or not
func (k Keeper) IsBeaconRegistered(ctx sdk.Context, beaconID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BeaconKey(beaconID))
}

// GetBeaconsIterator Get an iterator over all BEACONs in which the keys are the BEACON Ids and the values are the BEACONs
func (k Keeper) GetBeaconsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RegisteredBeaconPrefix)
}

// IterateBeacons iterates over the all the BEACON metadata and performs a callback function
func (k Keeper) IterateBeacons(ctx sdk.Context, cb func(beacon types.Beacon) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RegisteredBeaconPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var b types.Beacon
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &b)

		if cb(b) {
			break
		}
	}
}

// GetAllBeacons returns all the registered BEACON metadata from store
func (k Keeper) GetAllBeacons(ctx sdk.Context) (beacons types.Beacons) {
	k.IterateBeacons(ctx, func(wc types.Beacon) bool {
		beacons = append(beacons, wc)
		return false
	})
	return
}

// GetBeaconsFiltered retrieves BEACONs filtered by a given set of params which
// include pagination parameters along a moniker and owner address.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
// form.
func (k Keeper) GetBeaconsFiltered(ctx sdk.Context, params types.QueryBeaconParams) []types.Beacon {
	beacons := k.GetAllBeacons(ctx)
	filteredBeacons := make([]types.Beacon, 0, len(beacons))

	for _, b := range beacons {
		matchMoniker, matchOwner := true, true

		if len(params.Moniker) > 0 {
			matchMoniker = b.Moniker == params.Moniker
		}

		if len(params.Owner) > 0 {
			matchOwner = b.Owner.String() == params.Owner.String()
		}

		if matchMoniker && matchOwner {
			filteredBeacons = append(filteredBeacons, b)
		}
	}

	start, end := client.Paginate(len(filteredBeacons), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredBeacons = []types.Beacon{}
	} else {
		filteredBeacons = filteredBeacons[start:end]
	}

	return filteredBeacons
}

// RegisterBeacon registers a BEACON in the store
func (k Keeper) RegisterBeacon(ctx sdk.Context, moniker string, beaconName string, owner sdk.AccAddress) (uint64, error) {

	logger := k.Logger(ctx)

	//must have a moniker
	if len(moniker) == 0 {
		return 0, sdkerrors.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker")
	}
	if len(beaconName) > 128 {
		return 0, sdkerrors.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(moniker) > 64 {
		return 0, sdkerrors.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	beaconID, err := k.GetHighestBeaconID(ctx)
	if err != nil {
		return 0, err
	}

	params := types.NewQueryBeaconParams(1, 1, moniker, sdk.AccAddress{})
	beacons := k.GetBeaconsFiltered(ctx, params)

	if (len(beacons)) > 0 {
		errMsg := fmt.Sprintf("beacon already registered with moniker '%s' - id: %d, owner: %s", moniker, beacons[0].BeaconID, beacons[0].Owner)
		return 0, sdkerrors.Wrap(types.ErrBeaconAlreadyRegistered, errMsg)
	}

	beacon := k.GetBeacon(ctx, beaconID)

	beacon.BeaconID = beaconID
	beacon.Moniker = moniker
	beacon.LastTimestampID = 0
	beacon.Owner = owner
	beacon.Name = beaconName

	err = k.SetBeacon(ctx, beacon)
	if err != nil {
		return 0, err
	}

	k.SetHighestBeaconID(ctx, beaconID+1)

	if !ctx.IsCheckTx() {
		logger.Info("beacon registered", "id", beaconID, "moniker", moniker, "owner", owner.String())
	}

	return beaconID, nil
}
