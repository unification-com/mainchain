package wrkchain

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkchainId)

	for _, record := range data.RegisteredWrkchains {
		wrkChain := types.WrkChain{
			WrkchainId:   record.Wrkchain.WrkchainId,
			Moniker:      record.Wrkchain.Moniker,
			Name:         record.Wrkchain.Name,
			Genesis:      record.Wrkchain.Genesis,
			BaseType:     record.Wrkchain.BaseType,
			Lastblock:    record.Wrkchain.Lastblock,
			NumBlocks:    record.Wrkchain.NumBlocks,
			LowestHeight: record.Wrkchain.LowestHeight,
			RegTime:      record.Wrkchain.RegTime,
			Owner:        record.Wrkchain.Owner,
		}

		err := keeper.SetWrkChain(ctx, wrkChain)
		if err != nil {
			panic(err)
		}

		err = keeper.SetWrkChainStorageLimit(ctx, record.Wrkchain.WrkchainId, record.InStateLimit)
		if err != nil {
			panic(err)
		}

		for _, block := range record.Blocks {
			blk := types.WrkChainBlock{
				Height:     block.He,
				Blockhash:  block.Bh,
				Parenthash: block.Ph,
				Hash1:      block.H1,
				Hash2:      block.H2,
				Hash3:      block.H3,
				SubTime:    block.St,
			}

			err = keeper.SetWrkChainBlock(ctx, wrkChain.WrkchainId, blk)
			if err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	var records types.WrkChainExports
	initialWrkChainID, _ := k.GetHighestWrkChainID(ctx)

	wrkChains := k.GetAllWrkChains(ctx)

	if len(wrkChains) == 0 {
		return types.NewGenesisState(params, initialWrkChainID, nil)
	}

	for _, wc := range wrkChains {
		blockHashList := k.GetAllWrkChainBlockHashesForGenesisExport(ctx, wc.WrkchainId)
		lowestHeight := uint64(0)
		if len(blockHashList) > 0 {
			lowestHeight = blockHashList[0].He
		}

		wrkchainStorage, _ := k.GetWrkChainStorageLimit(ctx, wc.WrkchainId)

		records = append(records, types.WrkChainExport{
			Wrkchain: types.WrkChain{
				WrkchainId:   wc.WrkchainId,
				Moniker:      wc.Moniker,
				Name:         wc.Name,
				Genesis:      wc.Genesis,
				BaseType:     wc.BaseType,
				Lastblock:    wc.Lastblock,
				NumBlocks:    uint64(len(blockHashList)),
				LowestHeight: lowestHeight,
				RegTime:      wc.RegTime,
				Owner:        wc.Owner,
			},
			Blocks:       blockHashList,
			InStateLimit: wrkchainStorage.InStateLimit,
		})
	}

	return types.NewGenesisState(params, initialWrkChainID, records)
}
