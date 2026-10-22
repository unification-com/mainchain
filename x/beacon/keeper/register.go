package keeper

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/beacon/types"
)

//__BEACON_ID___________________________________________________________

// GetHighestBeaconID gets the highest BEACON ID
func (k Keeper) GetHighestBeaconID(ctx sdk.Context) (beaconID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestBeaconIDKey)
	if bz == nil {
		return 0, errorsmod.Wrapf(types.ErrInvalidGenesis, "initial beacon ID hasn't been set")
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
	if err := k.cdc.Unmarshal(bz, &beacon); err != nil {
		k.Logger(ctx).Error("corrupt beacon entry — treating as absent",
			"beacon_id", beaconID, "err", err)
		return types.Beacon{}, false
	}
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

// IterateBeacons iterates over the all the BEACON metadata and performs a callback function
func (k Keeper) IterateBeacons(ctx sdk.Context, cb func(beacon types.Beacon) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.RegisteredBeaconPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var b types.Beacon
		if err := k.cdc.Unmarshal(iterator.Value(), &b); err != nil {
			k.Logger(ctx).Error("skipping corrupt beacon entry during iteration",
				"err", err)
			continue
		}

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

// CountBeaconsForOwner counts the number of beacons registered to the given
// owner. Enforced against types.MaxBeaconsPerOwner at registration time.
// O(total beacons) — acceptable while overall beacon count stays modest;
// if growth becomes a concern, add a per-owner secondary index.
func (k Keeper) CountBeaconsForOwner(ctx sdk.Context, owner sdk.AccAddress) int {
	count := 0
	ownerStr := owner.String()
	k.IterateBeacons(ctx, func(b types.Beacon) bool {
		if b.Owner == ownerStr {
			count++
		}
		return false
	})
	return count
}

// Filtered beacon queries are served directly by the gRPC handler in
// grpc_query.go using query.FilteredPaginate against the store; no separate
// keeper helper is needed.

func (k Keeper) GetBeaconStorageLimit(ctx sdk.Context, beaconID uint64) (types.BeaconStorageLimit, bool) {
	store := ctx.KVStore(k.storeKey)
	if !k.HasBeaconStorageLimit(ctx, beaconID) {
		return types.BeaconStorageLimit{
			BeaconId:     beaconID,
			InStateLimit: types.DefaultStorageLimit,
		}, false
	}

	storageKey := types.BeaconStorageLimitKey(beaconID)
	bz := store.Get(storageKey)
	var storage types.BeaconStorageLimit
	if err := k.cdc.Unmarshal(bz, &storage); err != nil {
		k.Logger(ctx).Error("corrupt beacon storage limit entry — treating as absent",
			"beacon_id", beaconID, "err", err)
		return types.BeaconStorageLimit{
			BeaconId:     beaconID,
			InStateLimit: types.DefaultStorageLimit,
		}, false
	}
	return storage, true
}

func (k Keeper) HasBeaconStorageLimit(ctx sdk.Context, beaconID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	storageKey := types.BeaconStorageLimitKey(beaconID)
	return store.Has(storageKey)
}

func (k Keeper) SetBeaconStorageLimit(ctx sdk.Context, beaconId, limit uint64) error {

	store := ctx.KVStore(k.storeKey)
	storageLimit := types.BeaconStorageLimit{
		BeaconId:     beaconId,
		InStateLimit: limit,
	}
	store.Set(types.BeaconStorageLimitKey(beaconId), k.cdc.MustMarshal(&storageLimit))

	return nil
}

func (k Keeper) IncreaseInStateStorage(ctx sdk.Context, beaconId, amount uint64) error {
	beaconStorage, _ := k.GetBeaconStorageLimit(ctx, beaconId)
	newInStateLimit := beaconStorage.InStateLimit + amount
	err := k.SetBeaconStorageLimit(ctx, beaconId, newInStateLimit)

	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetMaxPurchasableSlots(ctx sdk.Context, beaconId uint64) uint64 {
	beaconStorage, found := k.GetBeaconStorageLimit(ctx, beaconId)
	if !found {
		return 0
	}

	maxStorageLimit := k.GetParamMaxStorageLimit(ctx)

	if beaconStorage.InStateLimit >= maxStorageLimit {
		return 0
	}

	return maxStorageLimit - beaconStorage.InStateLimit
}

// RegisterBeacon registers a BEACON in the store
func (k Keeper) RegisterNewBeacon(ctx sdk.Context, beacon types.Beacon) (uint64, error) {

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

	err = k.SetBeaconStorageLimit(ctx, beaconId, k.GetParamDefaultStorageLimit(ctx))
	if err != nil {
		return 0, err
	}

	k.SetHighestBeaconID(ctx, beaconId+1)

	return beaconId, nil
}
