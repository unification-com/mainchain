package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

// SetBeaconTimestamp Sets the Beacon timestamp struct for a beaconID + timestampID
func (k Keeper) SetBeaconTimestamp(ctx sdk.Context, beaconTimestamp types.BeaconTimestamp) sdk.Error {
	// must have an owner, Beacon ID, TimestampID and Hash
	if beaconTimestamp.Owner.Empty() || beaconTimestamp.BeaconID == 0 || beaconTimestamp.TimestampID == 0 || len(beaconTimestamp.Hash) == 0 || beaconTimestamp.SubmitTime == 0 {
		return sdk.ErrInternal("must include owner, id, submit time and hash")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.BeaconTimestampKey(beaconTimestamp.BeaconID, beaconTimestamp.TimestampID), k.cdc.MustMarshalBinaryLengthPrefixed(beaconTimestamp))

	return nil
}

// IsBeaconTimestampRecordedByID Check if the BEACON timestamp is present in the store or not, given
// the beaconID and timestampID
func (k Keeper) IsBeaconTimestampRecordedByID(ctx sdk.Context, beaconID uint64, timestampID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	timestampKey := types.BeaconTimestampKey(beaconID, timestampID)
	return store.Has(timestampKey)
}

func (k Keeper) IsBeaconTimestampRecordedByHashTime(ctx sdk.Context, beaconID uint64, hash string, subTime uint64) bool {
	params := types.NewQueryBeaconTimestampParams(1, 1, beaconID, hash, subTime)
	timestamps := k.GetBeaconTimestampsFiltered(ctx, params)
	return len(timestamps) > 0
}

// IsAuthorisedToRecord ensures only the BEACON owner is recording hashes
func (k Keeper) IsAuthorisedToRecord(ctx sdk.Context, beaconID uint64, recorder sdk.AccAddress) bool {
	return recorder.Equals(k.GetBeaconOwner(ctx, beaconID))
}

// GetBeaconTimestampByID Gets the beacon timestamp data for a beaconID and timestampID
func (k Keeper) GetBeaconTimestampByID(ctx sdk.Context, beaconID uint64, timestampID uint64) types.BeaconTimestamp {
	store := ctx.KVStore(k.storeKey)

	if !k.IsBeaconTimestampRecordedByID(ctx, beaconID, timestampID) {
		// return a new empty BeaconTimestamp struct
		return types.NewBeaconTimestamp()
	}

	timestampKey := types.BeaconTimestampKey(beaconID, timestampID)

	bz := store.Get(timestampKey)
	var beaconTimestamp types.BeaconTimestamp
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &beaconTimestamp)
	return beaconTimestamp
}

// IterateBeacons iterates over the all the BEACON's timestamps and performs a callback function
func (k Keeper) IterateBeaconTimestamps(ctx sdk.Context, beaconID uint64, cb func(beaconTimestamp types.BeaconTimestamp) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BeaconAllTimestampsKey(beaconID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var bts types.BeaconTimestamp
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &bts)

		if cb(bts) {
			break
		}
	}
}

// GetAllBeaconTimestamps Get an iterator over all a Beacon's timestamps in which the keys are the beaconID and the values are the BeaconTimestamps
func (k Keeper) GetAllBeaconTimestamps(ctx sdk.Context, beaconID uint64) (timestamps types.BeaconTimestamps) {

	k.IterateBeaconTimestamps(ctx, beaconID, func(bts types.BeaconTimestamp) bool {
		timestamps = append(timestamps, bts)
		return false
	})
	return
}

// GetBeaconTimestampsFiltered retrieves a BEACON's timestamps filtered by
// submit time, or the hash itself
func (k Keeper) GetBeaconTimestampsFiltered(ctx sdk.Context, params types.QueryBeaconTimestampParams) []types.BeaconTimestamp {
	timestamps := k.GetAllBeaconTimestamps(ctx, params.BeaconID)
	filteredTimestamps := make([]types.BeaconTimestamp, 0, len(timestamps))

	for _, bts := range timestamps {
		matchHash, matchSubmitTime := true, true

		if len(params.Hash) > 0 {
			matchHash = bts.Hash == params.Hash
		}

		if params.SubmitTime > 0 {
			matchSubmitTime = bts.SubmitTime == params.SubmitTime
		}

		if matchHash && matchSubmitTime {
			filteredTimestamps = append(filteredTimestamps, bts)
		}
	}

	start, end := client.Paginate(len(filteredTimestamps), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredTimestamps = []types.BeaconTimestamp{}
	} else {
		filteredTimestamps = filteredTimestamps[start:end]
	}

	return filteredTimestamps
}

// RecordBeaconTimestamp records a BEACON timestamp hash for a registered BEACON
func (k Keeper) RecordBeaconTimestamp(
	ctx sdk.Context,
	beaconID uint64,
	hash string,
	submitTime uint64,
	owner sdk.AccAddress) (uint64, sdk.Error) {

	logger := k.Logger(ctx)

	if !k.IsBeaconRegistered(ctx, beaconID) {
		// can't record hashes if BEACON isn't registered
		return 0, types.ErrBeaconDoesNotExist(k.codespace, "beacon does not exist")
	}

	beacon := k.GetBeacon(ctx, beaconID)

	if !k.IsAuthorisedToRecord(ctx, beacon.BeaconID, owner) {
		return 0, types.ErrNotBeaconOwner(k.codespace, "not authorised to record hashes for this beacon")
	}

	params := types.NewQueryBeaconTimestampParams(1, 1, beaconID, hash, submitTime)
	bts := k.GetBeaconTimestampsFiltered(ctx, params)

	if len(bts) > 0 {
		return 0, types.ErrBeaconTimestampAlreadyRecorded(k.codespace, fmt.Sprintf("timestamp hash %s already recorded at time %d", hash, submitTime))
	}

	timestampID := beacon.LastTimestampID + 1

	beaconTimestamp := k.GetBeaconTimestampByID(ctx, beacon.BeaconID, timestampID)

	beaconTimestamp.BeaconID = beacon.BeaconID
	beaconTimestamp.TimestampID = timestampID
	beaconTimestamp.Hash = hash
	beaconTimestamp.Owner = owner
	beaconTimestamp.SubmitTime = submitTime

	err := k.SetBeaconTimestamp(ctx, beaconTimestamp)

	if err != nil {
		return 0, err
	}

	err = k.SetLastTimestampID(ctx, beacon.BeaconID, timestampID)

	if err != nil {
		return 0, err
	}

	logger.Debug("beacon timestamp recorded", "bid", beacon.BeaconID, "tid", timestampID, "hash", hash, "owner", owner.String())

	return timestampID, nil
}
