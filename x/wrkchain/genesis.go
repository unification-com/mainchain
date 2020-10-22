package wrkchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkChainID)

	logger := keeper.Logger(ctx)

	for _, record := range data.WrkChains {
		wrkChain := record.WrkChain
		err := keeper.SetWrkChain(ctx, wrkChain)
		if err != nil {
			panic(err)
		}

		logger.Info("Registering WRKChain", "wc_id", wrkChain.WrkChainID)

		for _, block := range record.WrkChainBlocks {
			blk := WrkChainBlock{
				WrkChainID: wrkChain.WrkChainID,
				Height: block.Height,
				BlockHash: block.BlockHash,
				ParentHash: block.ParentHash,
				Hash1: block.Hash1,
				Hash2: block.Hash2,
				Hash3: block.Hash3,
				SubmitTime: block.SubmitTime,
				Owner: wrkChain.Owner,
			}
			//logger.Info("Registering Block for WRKChain", "wc_id", wrkChain.WrkChainID, "h", block.Height)
			err = keeper.SetWrkChainBlock(ctx, blk)
			if err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	var records []WrkChainExport
	initialWrkChainID, _ := k.GetHighestWrkChainID(ctx)

	wrkChains := k.GetAllWrkChains(ctx)

	if len(wrkChains) == 0 {
		return GenesisState{
			Params:             params,
			StartingWrkChainID: initialWrkChainID,
			WrkChains:          nil,
		}
	}

	for _, wc := range wrkChains {
		wrkchainId := wc.WrkChainID
		blockHashList := k.GetAllWrkChainBlockHashesForGenesisExport(ctx, wrkchainId)

		records = append(records, WrkChainExport{WrkChain: wc, WrkChainBlocks: blockHashList})
	}
	return GenesisState{
		Params:             params,
		StartingWrkChainID: initialWrkChainID,
		WrkChains:          records,
	}
}
