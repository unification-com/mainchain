package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/types"
)

// SetBeaconTimestamp Sets the Beacon timestamp struct for a beaconID + timestampID
func (k Keeper) SetBeaconTimestamp(ctx sdk.Context, beaconTimestamp types.BeaconTimestamp) error {

	store := ctx.KVStore(k.storeKey)
	store.Set(types.BeaconTimestampKey(beaconTimestamp.BeaconId, beaconTimestamp.TimestampId), k.cdc.MustMarshalBinaryBare(&beaconTimestamp))

	return nil
}

// IsBeaconTimestampRecordedByID Deep Check if the BEACON timestamp is present in the store or not, given
// the beaconID and timestampID
func (k Keeper) IsBeaconTimestampRecordedByID(ctx sdk.Context, beaconID uint64, timestampID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	timestampKey := types.BeaconTimestampKey(beaconID, timestampID)
	return store.Has(timestampKey)
}

// IsAuthorisedToRecord ensures only the BEACON owner is recording hashes
func (k Keeper) IsAuthorisedToRecord(ctx sdk.Context, beaconID uint64, recorder sdk.AccAddress) bool {
	return recorder.Equals(k.GetBeaconOwner(ctx, beaconID))
}

// GetBeaconTimestampByID Gets the beacon timestamp data for a beaconID and timestampID
func (k Keeper) GetBeaconTimestampByID(ctx sdk.Context, beaconID uint64, timestampID uint64) (types.BeaconTimestamp, bool) {
	store := ctx.KVStore(k.storeKey)

	if !k.IsBeaconTimestampRecordedByID(ctx, beaconID, timestampID) {
		// return a new empty BeaconTimestamp struct
		return types.BeaconTimestamp{}, false
	}

	timestampKey := types.BeaconTimestampKey(beaconID, timestampID)

	bz := store.Get(timestampKey)
	var beaconTimestamp types.BeaconTimestamp
	k.cdc.MustUnmarshalBinaryBare(bz, &beaconTimestamp)
	return beaconTimestamp, true
}

// IterateBeacons iterates over the all the BEACON's timestamps and performs a callback function
func (k Keeper) IterateBeaconTimestamps(ctx sdk.Context, beaconID uint64, cb func(beaconTimestamp types.BeaconTimestamp) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BeaconAllTimestampsKey(beaconID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bts types.BeaconTimestamp
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &bts)

		if cb(bts) {
			break
		}
	}
}

// IterateBeaconTimestampsReverse iterates over the all the BEACON's timestamps in reverse order and performs a callback function
func (k Keeper) IterateBeaconTimestampsReverse(ctx sdk.Context, beaconID uint64, cb func(beaconTimestamp types.BeaconTimestamp) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.BeaconAllTimestampsKey(beaconID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bts types.BeaconTimestamp
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &bts)

		if cb(bts) {
			break
		}
	}
}

// GetAllBeaconTimestamps Get an iterator over all a Beacon's timestamps in which the keys are the beaconID
// and the values are the BeaconTimestamps
func (k Keeper) GetAllBeaconTimestamps(ctx sdk.Context, beaconID uint64) (timestamps []types.BeaconTimestamp) {

	k.IterateBeaconTimestamps(ctx, beaconID, func(bts types.BeaconTimestamp) bool {
		timestamps = append(timestamps, bts)
		return false
	})
	return
}

func prependTimestamp(x []types.BeaconTimestampGenesisExport, y types.BeaconTimestampGenesisExport) []types.BeaconTimestampGenesisExport {
	x = append(x, y)
	copy(x[1:], x)
	x[0] = y
	return x
}

// GetAllBeaconTimestampsForExport Get an iterator over all a Beacon's timestamps in which an optimised version of
// the timestamp data is returned for genesis export
func (k Keeper) GetAllBeaconTimestampsForExport(ctx sdk.Context, beaconID uint64) (timestamps []types.BeaconTimestampGenesisExport) {

	count := 0
	k.IterateBeaconTimestampsReverse(ctx, beaconID, func(bts types.BeaconTimestamp) bool {
		btsExp := types.BeaconTimestampGenesisExport{
			Id: bts.TimestampId,
			T:  bts.SubmitTime,
			H:  bts.Hash,
		}
		timestamps = prependTimestamp(timestamps, btsExp) // append(timestamps, btsExp)
		count = count + 1

		return count == types.MaxHashSubmissionsKeepInState
	})
	return
}

// RecordBeaconTimestamp records a BEACON timestamp hash for a registered BEACON
func (k Keeper) RecordNewBeaconTimestamp(
	ctx sdk.Context,
	beaconId uint64,
	hash string,
	submitTime uint64,
	owner string) (uint64, error) {

	logger := k.Logger(ctx)

	beacon, _ := k.GetBeacon(ctx, beaconId)

	timestampId := beacon.LastTimestampId + 1

	// we're only ever recording new BEACON hashes, never updating existing. Handler has already run
	// checks for authorisation etc.
	beaconTimestamp := types.BeaconTimestamp{
		BeaconId:    beaconId,
		TimestampId: timestampId,
		SubmitTime:  submitTime,
		Hash:        hash,
		Owner:       owner,
	}

	err := k.SetBeaconTimestamp(ctx, beaconTimestamp)

	if err != nil {
		return 0, err
	}

	err = k.SetLastTimestampID(ctx, beacon.BeaconId, timestampId)

	if err != nil {
		return 0, err
	}

	if !ctx.IsCheckTx() {
		logger.Debug("beacon timestamp recorded", "id", beaconId, "hash", hash, "time", submitTime)
	}

	return timestampId, nil
}
