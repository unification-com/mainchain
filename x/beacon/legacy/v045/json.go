package v045

import (
	v040 "github.com/unification-com/mainchain/x/beacon/legacy/v040"
	"github.com/unification-com/mainchain/x/beacon/types"
)

// MigrateJSON accepts exported v0.40 x/beacon genesis state and migrates it to
// v0.45 (0.43) x/beacon genesis state.
func MigrateJSON(oldBeaconState *v040.GenesisState) *types.GenesisState {
	newBeacons := make(types.BeaconExports, len(oldBeaconState.RegisteredBeacons))

	for i, oldBeacon := range oldBeaconState.RegisteredBeacons {
		firstId := uint64(0)
		newBeaconTimestamps := make(types.BeaconTimestampGenesisExports, len(oldBeacon.Timestamps))
		for j, oldBeaconTimestamp := range oldBeacon.Timestamps {
			if firstId == 0 || oldBeaconTimestamp.Id < firstId {
				firstId = oldBeaconTimestamp.Id
			}
			newBeaconTimestamps[j] = types.BeaconTimestampGenesisExport{
				Id: oldBeaconTimestamp.Id,
				T:  oldBeaconTimestamp.T,
				H:  oldBeaconTimestamp.H,
			}
		}
		newBeacons[i] = types.BeaconExport{
			Beacon: types.Beacon{
				BeaconId:        oldBeacon.Beacon.BeaconId,
				Moniker:         oldBeacon.Beacon.Moniker,
				Name:            oldBeacon.Beacon.Name,
				LastTimestampId: oldBeacon.Beacon.LastTimestampId,
				Owner:           oldBeacon.Beacon.Owner,
				NumInState:      uint64(len(oldBeacon.Timestamps)),
				FirstIdInState:  firstId,
				RegTime:         oldBeacon.Beacon.RegTime, // reg time was never recorded originally
			},
			Timestamps:   newBeaconTimestamps,
			InStateLimit: types.DefaultStorageLimit,
		}
	}

	return &types.GenesisState{
		Params: types.Params{
			FeeRegister:         oldBeaconState.Params.FeeRegister,
			FeeRecord:           oldBeaconState.Params.FeeRecord,
			FeePurchaseStorage:  types.PurchaseStorageFee,
			Denom:               oldBeaconState.Params.Denom,
			DefaultStorageLimit: types.DefaultStorageLimit,
			MaxStorageLimit:     types.DefaultMaxStorageLimit,
		},
		StartingBeaconId:  oldBeaconState.StartingBeaconId,
		RegisteredBeacons: newBeacons,
	}
}
