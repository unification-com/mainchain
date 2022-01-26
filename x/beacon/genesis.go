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

	for _, record := range data.RegisteredBeacons {
		beacon := types.Beacon{
			BeaconId:        record.Beacon.BeaconId,
			Moniker:         record.Beacon.Moniker,
			Name:            record.Beacon.Name,
			LastTimestampId: record.Beacon.LastTimestampId,
			Owner:           record.Beacon.Owner,
			NumInState:      record.Beacon.NumInState,
			FirstIdInState:  record.Beacon.FirstIdInState,
		}

		err := keeper.SetBeacon(ctx, beacon)
		if err != nil {
			panic(err)
		}

		for _, timestamp := range record.Timestamps {

			bts := types.BeaconTimestamp{
				TimestampId: timestamp.Id,
				SubmitTime:  timestamp.T,
				Hash:        timestamp.H,
			}

			err = keeper.SetBeaconTimestamp(ctx, beacon.BeaconId, bts)
			if err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	var records types.BeaconExports
	initialBeaconID, _ := k.GetHighestBeaconID(ctx)

	beacons := k.GetAllBeacons(ctx)

	if len(beacons) == 0 {
		return types.NewGenesisState(params, initialBeaconID, nil)
	}

	for _, b := range beacons {
		timestamps := k.GetAllBeaconTimestampsForExport(ctx, b.BeaconId)

		records = append(records, types.BeaconExport{
			Beacon: types.Beacon{
				BeaconId:        b.BeaconId,
				Moniker:         b.Moniker,
				Name:            b.Name,
				LastTimestampId: b.LastTimestampId,
				Owner:           b.Owner,
			},
			Timestamps: timestamps,
		})
	}
	return types.NewGenesisState(params, initialBeaconID, records)
}
