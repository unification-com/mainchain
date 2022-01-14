package wrkchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkchainId)

	for _, record := range data.RegisteredWrkchains {
		wrkChain := types.WrkChain{
			WrkchainId: record.Wrkchain.WrkchainId,
			Moniker:    record.Wrkchain.Moniker,
			Name:       record.Wrkchain.Name,
			Genesis:    record.Wrkchain.Genesis,
			Type:       record.Wrkchain.Type,
			Lastblock:  record.Wrkchain.Lastblock,
			NumBlocks:  record.Wrkchain.NumBlocks,
			RegTime:    record.Wrkchain.RegTime,
			Owner:      record.Wrkchain.Owner,
		}

		err := keeper.SetWrkChain(ctx, wrkChain)
		if err != nil {
			panic(err)
		}

		for _, block := range record.Blocks {
			blk := types.WrkChainBlock{
				WrkchainId: wrkChain.WrkchainId,
				Height:     block.He,
				Blockhash:  block.Bh,
				Parenthash: block.Ph,
				Hash1:      block.H1,
				Hash2:      block.H2,
				Hash3:      block.H3,
				SubTime:    block.St,
				Owner:      wrkChain.Owner,
			}

			err = keeper.SetWrkChainBlock(ctx, blk)
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

		records = append(records, types.WrkChainExport{
			Wrkchain: types.WrkChain{
				WrkchainId: wc.WrkchainId,
				Moniker:    wc.Moniker,
				Name:       wc.Name,
				Genesis:    wc.Genesis,
				Type:       wc.Type,
				Lastblock:  wc.Lastblock,
				NumBlocks:  wc.NumBlocks,
				RegTime:    wc.RegTime,
				Owner:      wc.Owner,
			},
			Blocks: blockHashList,
		})
	}

	return types.NewGenesisState(params, initialWrkChainID, records)
}
