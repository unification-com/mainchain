package wrkchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	undtypes "github.com/unification-com/mainchain/types"
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

			// also update NumBlocks!
			err = keeper.SetNumBlocks(ctx, wrkChain.WrkChainID)
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
	exportWrkChainDataIds := viper.GetIntSlice(undtypes.FlagExportIncludeWrkchainData)

	if len(wrkChains) == 0 {
		return GenesisState{
			Params:             params,
			StartingWrkChainID: initialWrkChainID,
			WrkChains:          nil,
		}
	}

	for _, wc := range wrkChains {
		exportData := false
		for _, expWrkChainId := range exportWrkChainDataIds {
			if uint64(expWrkChainId) == wc.WrkChainID {
				exportData = true
			}
		}

		if exportData {
			wrkchainId := wc.WrkChainID
			blockHashList := k.GetAllWrkChainBlockHashesForGenesisExport(ctx, wrkchainId)
			if blockHashList == nil {
				blockHashList = WrkChainBlocksGenesisExport{}
			}
			records = append(records, WrkChainExport{WrkChain: wc, WrkChainBlocks: blockHashList})
		} else {
			wc.LastBlock = 0
			wc.NumberBlocks = 0
			records = append(records, WrkChainExport{WrkChain: wc, WrkChainBlocks: WrkChainBlocksGenesisExport{}})
		}
	}
	return GenesisState{
		Params:             params,
		StartingWrkChainID: initialWrkChainID,
		WrkChains:          records,
	}
}
