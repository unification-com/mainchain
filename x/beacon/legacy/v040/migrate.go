package v040

import (
	v038 "github.com/unification-com/mainchain/x/beacon/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/beacon/types"
)

func Migrate(oldBeaconState v038.GenesisState) *v040.GenesisState {

	newBeacons := make(v040.BeaconExports, len(oldBeaconState.Beacons))

	for i, oldBeacon := range oldBeaconState.Beacons {
		firstId := uint64(0)
		newBeaconTimestamps := make(v040.BeaconTimestampGenesisExports, len(oldBeacon.BeaconTimestamps))
		for j, oldBeaconTimestamp := range oldBeacon.BeaconTimestamps {

			if firstId == 0 || oldBeaconTimestamp.TimestampID < firstId {
				firstId = oldBeaconTimestamp.TimestampID
			}

			newBeaconTimestamps[j] = v040.BeaconTimestampGenesisExport{
				Id: oldBeaconTimestamp.TimestampID,
				T:  oldBeaconTimestamp.SubmitTime,
				H:  oldBeaconTimestamp.Hash,
			}
		}

		newBeacons[i] = v040.BeaconExport{
			Beacon: v040.Beacon{
				BeaconId:        oldBeacon.Beacon.BeaconID,
				Moniker:         oldBeacon.Beacon.Moniker,
				Name:            oldBeacon.Beacon.Name,
				LastTimestampId: oldBeacon.Beacon.LastTimestampID,
				Owner:           oldBeacon.Beacon.Owner.String(),
				NumInState:      uint64(len(oldBeacon.BeaconTimestamps)),
				FirstIdInState:  firstId,
			},
			Timestamps: newBeaconTimestamps,
		}
	}

	return &v040.GenesisState{
		Params: v040.Params{
			FeeRegister: oldBeaconState.Params.FeeRegister,
			FeeRecord:   oldBeaconState.Params.FeeRecord,
			Denom:       oldBeaconState.Params.Denom,
		},
		StartingBeaconId:  oldBeaconState.StartingBeaconID,
		RegisteredBeacons: newBeacons,
	}
}
