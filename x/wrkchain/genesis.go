package wrkchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkChainID)
	for _, record := range data.WrkChains {
		keeper.SetWrkChain(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []WrkChain
	//iterator := k.GetWrkChainsIterator(ctx)
	//for ; iterator.Valid(); iterator.Next() {
	//	wrkchainId := string(iterator.Key())
	//	wrkChain := k.GetWrkChain(ctx, wrkchainId)
	//	records = append(records, wrkChain)
	//}
	return GenesisState{WrkChains: records}
}
