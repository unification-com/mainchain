package beacon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestBeaconID(ctx, data.StartingBeaconID)

	logger := keeper.Logger(ctx)

	for _, record := range data.Beacons {
		beacon := record.Beacon
		err := keeper.SetBeacon(ctx, beacon)
		if err != nil {
			panic(err)
		}

		if beacon.LastTimestampID > 0 {
			err = keeper.SetLastTimestampID(ctx, beacon.BeaconID, beacon.LastTimestampID)
			if err != nil {
				panic(err)
			}
		}

		logger.Info("setting beacon", "bid", beacon.BeaconID)

		for _, timestamp := range record.BeaconTimestamps {

			bts := BeaconTimestamp{
				BeaconID:    beacon.BeaconID,
				TimestampID: timestamp.TimestampID,
				SubmitTime:  timestamp.SubmitTime,
				Hash:        timestamp.Hash,
				Owner:       beacon.Owner,
			}

			err = keeper.SetBeaconTimestamp(ctx, bts)
			if err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	var records []BeaconExport
	initialBeaconID, _ := k.GetHighestBeaconID(ctx)

	beacons := k.GetAllBeacons(ctx)

	if len(beacons) == 0 {
		return GenesisState{
			Params:           params,
			StartingBeaconID: initialBeaconID,
			Beacons:          nil,
		}
	}

	for _, b := range beacons {
		beaconID := b.BeaconID
		timestamps := k.GetAllBeaconTimestampsForExport(ctx, beaconID)
		if timestamps == nil {
			timestamps = BeaconTimestampsGenesisExport{}
		}
		records = append(records, BeaconExport{Beacon: b, BeaconTimestamps: timestamps})
	}
	return GenesisState{
		Params:           params,
		StartingBeaconID: initialBeaconID,
		Beacons:          records,
	}
}
