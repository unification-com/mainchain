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

	logger := keeper.Logger(ctx)

	for _, record := range data.RegisteredWrkchains {
		wrkChain := types.WrkChain{
			WrkchainId: record.Wrkchain.WrkchainId,
			Moniker: record.Wrkchain.Moniker,
			Name: record.Wrkchain.Name,
			Genesis: record.Wrkchain.Genesis,
			Type: record.Wrkchain.Type,
			Lastblock: record.Wrkchain.Lastblock,
			NumBlocks: record.Wrkchain.NumBlocks,
			RegTime: record.Wrkchain.RegTime,
			Owner: record.Wrkchain.Owner,
		}

		err := keeper.SetWrkChain(ctx, wrkChain)
		if err != nil {
			panic(err)
		}

		logger.Debug("Registering WRKChain", "wc_id", wrkChain.WrkchainId)

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

			// also update NumBlocks!
			err = keeper.SetNumBlocks(ctx, wrkChain.WrkchainId)
			if err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	params := k.GetParams(ctx)
	var records []types.WrkChainExport
	initialWrkChainID, _ := k.GetHighestWrkChainID(ctx)

	wrkChains := k.GetAllWrkChains(ctx)

	if len(wrkChains) == 0 {
		return types.GenesisState{
			Params:              params,
			StartingWrkchainId:  initialWrkChainID,
			RegisteredWrkchains: records,
		}
	}

	for _, wc := range wrkChains {
		blockHashList := k.GetAllWrkChainBlockHashesForGenesisExport(ctx, wc.WrkchainId)
		if blockHashList == nil {
			blockHashList = []types.WrkChainBlockGenesisExport{}
		}
		records = append(records, types.WrkChainExport{Wrkchain: &wc, Blocks: blockHashList})
	}

	return types.GenesisState{
		Params:              params,
		StartingWrkchainId:  initialWrkChainID,
		RegisteredWrkchains: records,
	}
}