package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v040 "github.com/unification-com/mainchain/x/beacon/legacy/v040"
	"github.com/unification-com/mainchain/x/beacon/types"
)

func migrateParams(ctx sdk.Context, paramsSubspace paramstypes.Subspace) {
	// Add the new module params
	var oldFeeReeRegParams uint64
	paramsSubspace.Get(ctx, types.KeyFeeRegister, &oldFeeReeRegParams)

	var oldFeeRecordParams uint64
	paramsSubspace.Get(ctx, types.KeyFeeRecord, &oldFeeRecordParams)

	var oldDenomParams string
	paramsSubspace.Get(ctx, types.KeyDenom, &oldDenomParams)

	params := types.NewParams(
		oldFeeReeRegParams,
		oldFeeRecordParams,
		types.PurchaseStorageFee,
		oldDenomParams,
		types.DefaultStorageLimit,
		types.DefaultMaxStorageLimit,
	)

	// save new paramset
	paramsSubspace.SetParamSet(ctx, &params)
}

func pruneBeaconTimestamps(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec, beaconId uint64) (uint64, uint64, error) {
	store := ctx.KVStore(storeKey)
	// reverse order (desc) to get newest first. Only want to prune old timestamps
	// when default in-state storage limit is reached
	iterator := sdk.KVStoreReversePrefixIterator(store, types.BeaconAllTimestampsKey(beaconId))
	numInState := uint64(0)
	firstIdInState := uint64(0)
	limitCounter := uint64(0)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldBeaconTimestamp v040.BeaconTimestamp
		if err := cdc.Unmarshal(iterator.Value(), &oldBeaconTimestamp); err != nil {
			return 0, 0, sdkerrors.Wrap(err, "failed to unmarshal beacon timestamp")
		}

		limitCounter = limitCounter + 1

		if limitCounter > types.DefaultStorageLimit {
			store.Delete(types.BeaconTimestampKey(beaconId, oldBeaconTimestamp.TimestampId))
		} else {
			firstIdInState = oldBeaconTimestamp.TimestampId
			numInState = numInState + 1
		}
	}

	return numInState, firstIdInState, nil
}

func pruneBeacons(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	// loop through Beacons, then each beacon's timestamps
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RegisteredBeaconPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldBeacon v040.Beacon
		if err := cdc.Unmarshal(iterator.Value(), &oldBeacon); err != nil {
			return sdkerrors.Wrap(err, "failed to unmarshal beacon")
		}

		numInState, firstIdInState, err := pruneBeaconTimestamps(ctx, storeKey, cdc, oldBeacon.BeaconId)

		if err != nil {
			return err
		}

		newBeacon := types.Beacon{
			BeaconId:        oldBeacon.BeaconId,
			Moniker:         oldBeacon.Moniker,
			Name:            oldBeacon.Name,
			LastTimestampId: oldBeacon.LastTimestampId,
			FirstIdInState:  firstIdInState,
			NumInState:      numInState,
			RegTime:         oldBeacon.RegTime,
			Owner:           oldBeacon.Owner,
		}

		// save the updated BEACON
		newBeaconBz, err := cdc.Marshal(&newBeacon)

		if err != nil {
			return err
		}

		store.Set(types.BeaconKey(newBeacon.BeaconId), newBeaconBz)

		newBeaconStorage := types.BeaconStorageLimit{
			BeaconId:     oldBeacon.BeaconId,
			InStateLimit: types.DefaultStorageLimit,
		}

		newBeaconStorageBz, err := cdc.Marshal(&newBeaconStorage)

		if err != nil {
			return err
		}

		store.Set(types.BeaconStorageLimitKey(newBeacon.BeaconId), newBeaconStorageBz)
	}

	return nil
}

// MigrateStore performs in-place store migrations from SDK v0.42 of the BEACON module to SDK v0.45.
// The migration includes:
//
// - Adding new params
// - Setting default in-store limit for BEACONs
// - Pruning all old timestamps exceeding in-state limit
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, paramsSubspace paramstypes.Subspace, cdc codec.BinaryCodec) error {

	// migrate Params
	migrateParams(ctx, paramsSubspace)

	// prune timestamps
	return pruneBeacons(ctx, storeKey, cdc)

}
