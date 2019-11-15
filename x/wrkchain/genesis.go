package wrkchain

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	keeper.SetHighestWrkChainID(ctx, data.StartingWrkChainID)
	for _, record := range data.WrkChains {
		keeper.SetWrkChain(ctx, record.WrkChain)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []WrkChainExport

	iterator := k.GetWrkChainsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		wrkchainId := iterator.Key()
		num := binary.LittleEndian.Uint64(wrkchainId)
		blockHashList := k.GetWrkChainBlockHashes(ctx, num)

		var hashes []types.WrkChainBlock

		for _, value := range blockHashList {
			hash := types.WrkChainBlock{
				WrkChainID:   num,
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

		wrkChain := k.GetWrkChain(ctx, num)
		records = append(records, WrkChainExport{WrkChain: wrkChain, WrkChainBlocks: hashes})
	}
	return GenesisState{WrkChains: records}
}
