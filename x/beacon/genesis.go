package beacon

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestBeaconID(ctx, data.StartingBeaconID)
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	return GenesisState{
		Params:             params,
	}
}
