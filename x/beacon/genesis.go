package beacon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestBeaconID(ctx, data.StartingBeaconID)

	logger := ctx.Logger()

	for _, record := range data.Beacons {
		beacon := record.Beacon
		err := keeper.SetBeacon(ctx, beacon)
		if err != nil {
			panic(err)
		}

		logger.Info("setting beacon", beacon.BeaconID)

		for _, timestamp := range record.BeaconTimestamps {
			logger.Info("setting timestamp for beacon", timestamp.BeaconID, timestamp.TimestampID)
			err = keeper.SetBeaconTimestamp(ctx, timestamp)
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
		timestamps := k.GetAllBeaconTimestamps(ctx, beaconID)

		var tss []BeaconTimestamp

		for _, value := range timestamps {
			ts := BeaconTimestamp{
				BeaconID:    value.BeaconID,
				TimestampID: value.TimestampID,
				Hash:        value.Hash,
				SubmitTime:  value.SubmitTime,
				Owner:       value.Owner,
			}
			tss = append(tss, ts)
		}

		records = append(records, BeaconExport{Beacon: b, BeaconTimestamps: tss})
	}
	return GenesisState{
		Params:           params,
		StartingBeaconID: initialBeaconID,
		Beacons:          records,
	}
}
