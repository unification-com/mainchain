package wrkchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkChainID)

	logger := ctx.Logger()

	for _, record := range data.WrkChains {
		wrkChain := record.WrkChain
		_ = keeper.SetWrkChain(ctx, wrkChain)
		_, _ = keeper.RegisterWrkChain(ctx, wrkChain.Moniker, wrkChain.Name, wrkChain.GenesisHash, wrkChain.Owner)

		logger.Info("Registering WRKChain", wrkChain.WrkChainID)

		for _, block := range record.WrkChainBlocks {
			logger.Info("Registering Block for WRKChain", wrkChain.WrkChainID, block.Height)
			_ = keeper.RecordWrkchainHashes(ctx, block.WrkChainID, block.Height, block.BlockHash, block.ParentHash, block.Hash1, block.Hash2, block.Hash3, block.Owner)
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
		blockHashList := k.GetWrkChainBlockHashes(ctx, wrkchainId)

		var hashes []types.WrkChainBlock

		for _, value := range blockHashList {
			hash := types.WrkChainBlock{
				WrkChainID:   value.WrkChainID,
				Height:       value.Height,
				BlockHash:    value.BlockHash,
				ParentHash:   value.ParentHash,
				Hash1:        value.Hash1,
				Hash2:        value.Hash2,
				Hash3:        value.Hash3,
				SubmitTime:   value.SubmitTime,
				SubmitHeight: value.SubmitHeight,
				Owner:        value.Owner,
			}
			hashes = append(hashes, hash)
		}

		records = append(records, WrkChainExport{WrkChain: wc, WrkChainBlocks: hashes})
	}
	return GenesisState{
		Params:             params,
		StartingWrkChainID: initialWrkChainID,
		WrkChains:          records,
	}
}
