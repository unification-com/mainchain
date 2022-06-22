package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/beacon/types"
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
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BeaconKey(beacon.BeaconId), k.cdc.MustMarshal(&beacon))

	return nil
}

// GetBeacon Gets the entire BEACON metadata struct for a beaconID
func (k Keeper) GetBeacon(ctx sdk.Context, beaconID uint64) (types.Beacon, bool) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsBeaconRegistered(ctx, beaconID) {
		// return a new empty Beacon struct
		return types.Beacon{}, false
	}
	bz := store.Get(types.BeaconKey(beaconID))
	var beacon types.Beacon
	k.cdc.MustUnmarshal(bz, &beacon)
	return beacon, true
}

// GetBeaconOwner - get the current owner of a BEACON
func (k Keeper) GetBeaconOwner(ctx sdk.Context, beaconID uint64) sdk.AccAddress {
	beacon, found := k.GetBeacon(ctx, beaconID)

	if !found {
		return sdk.AccAddress{}
	}

	accAddr, accErr := sdk.AccAddressFromBech32(beacon.Owner)
	if accErr != nil {
		return sdk.AccAddress{}
	}
	return accAddr
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
		k.cdc.MustUnmarshal(iterator.Value(), &b)

		if cb(b) {
			break
		}
	}
}

// GetAllBeacons returns all the registered BEACON metadata from store
func (k Keeper) GetAllBeacons(ctx sdk.Context) (beacons []types.Beacon) {
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
func (k Keeper) GetBeaconsFiltered(ctx sdk.Context, params types.QueryBeaconsFilteredRequest) []types.Beacon {
	beacons := k.GetAllBeacons(ctx)
	filteredBeacons := make([]types.Beacon, 0, len(beacons))

	for _, b := range beacons {
		matchMoniker, matchOwner := true, true

		if len(params.Moniker) > 0 {
			matchMoniker = b.Moniker == params.Moniker
		}

		if len(params.Owner) > 0 {
			matchOwner = b.Owner == params.Owner
		}

		if matchMoniker && matchOwner {
			filteredBeacons = append(filteredBeacons, b)
		}
	}

	// Todo - need to migrate this to proper pagination
	start, end := client.Paginate(len(filteredBeacons), 1, 100, 100)
	if start < 0 || end < 0 {
		filteredBeacons = []types.Beacon{}
	} else {
		filteredBeacons = filteredBeacons[start:end]
	}

	return filteredBeacons
}

// RegisterBeacon registers a BEACON in the store
func (k Keeper) RegisterNewBeacon(ctx sdk.Context, beacon types.Beacon) (uint64, error) {

	logger := k.Logger(ctx)

	beaconId, err := k.GetHighestBeaconID(ctx)
	if err != nil {
		return 0, err
	}

	beacon.BeaconId = beaconId
	beacon.LastTimestampId = 0
	beacon.FirstIdInState = 0
	beacon.NumInState = 0
	beacon.RegTime = uint64(ctx.BlockTime().Unix())

	err = k.SetBeacon(ctx, beacon)
	if err != nil {
		return 0, err
	}

	k.SetHighestBeaconID(ctx, beaconId+1)

	if !ctx.IsCheckTx() {
		logger.Debug("beacon registered", "id", beaconId, "moniker", beacon.Moniker, "owner", beacon.Owner)
	}

	return beaconId, nil
}
