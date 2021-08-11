package beacon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/unification-com/mainchain/x/beacon/keeper"
	"github.com/unification-com/mainchain/x/beacon/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestBeaconID(ctx, data.StartingBeaconId)

	logger := keeper.Logger(ctx)

	for _, record := range data.RegisteredBeacons {
		beacon := types.Beacon{
			BeaconId: record.Beacon.BeaconId,
			Moniker: record.Beacon.Moniker,
			Name: record.Beacon.Name,
			LastTimestampId: record.Beacon.LastTimestampId,
			Owner: record.Beacon.Owner,
		}

		err := keeper.SetBeacon(ctx, beacon)
		if err != nil {
			panic(err)
		}

		if beacon.LastTimestampId > 0 {
			err = keeper.SetLastTimestampID(ctx, beacon.BeaconId, beacon.LastTimestampId)
			if err != nil {
				panic(err)
			}
		}

		logger.Debug("setting beacon", "bid", beacon.BeaconId)

		for _, timestamp := range record.Timestamps {

			bts := types.BeaconTimestamp{
				BeaconId:    beacon.BeaconId,
				TimestampId: timestamp.Id,
				SubmitTime:  timestamp.T,
				Hash:        timestamp.H,
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

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	params := k.GetParams(ctx)
	var records []types.BeaconExport
	initialBeaconID, _ := k.GetHighestBeaconID(ctx)

	beacons := k.GetAllBeacons(ctx)

	if len(beacons) == 0 {
		return types.GenesisState{
			Params:            params,
			StartingBeaconId:  initialBeaconID,
			RegisteredBeacons: nil,
		}
	}

	for _, b := range beacons {
		beaconID := b.BeaconId
		timestamps := k.GetAllBeaconTimestampsForExport(ctx, beaconID)
		if timestamps == nil {
			timestamps = []types.BeaconTimestampGenesisExport{}
		}
		records = append(records, types.BeaconExport{Beacon: &b, Timestamps: timestamps})
	}
	return types.GenesisState{
		Params:            params,
		StartingBeaconId:  initialBeaconID,
		RegisteredBeacons: records,
	}
}
