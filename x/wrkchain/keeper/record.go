package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// SetWrkChainBlock Sets the WrkChain Block struct for a wrkchainId & height
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainBlock types.WrkChainBlock) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainBlockKey(wrkchainBlock.WrkchainId, wrkchainBlock.Height), k.cdc.MustMarshalBinaryBare(&wrkchainBlock))

	return nil
}

// QuickCheckHeightIsRecorded Checks if the given height can be recorded
func (k Keeper) QuickCheckHeightIsRecorded(ctx sdk.Context, wrkchainId uint64, height uint64) bool {

	wrkchain, _ := k.GetWrkChain(ctx, wrkchainId)

	// only check if height being submitted is <= last recorded height.
	// Otherwise, no need to check entire db
	if height <= wrkchain.Lastblock && height > 0 {
		if height == wrkchain.Lastblock {
			return true
		} else {
			store := ctx.KVStore(k.storeKey)
			blockKey := types.WrkChainBlockKey(wrkchainId, height)
			return store.Has(blockKey)
		}
	}
	return false
}

// IsWrkChainBlockRecorded Check if the WrkChainBlock is present in the store or not
func (k Keeper) IsWrkChainBlockRecorded(ctx sdk.Context, wrkchainId uint64, height uint64) bool {
	store := ctx.KVStore(k.storeKey)
	blockKey := types.WrkChainBlockKey(wrkchainId, height)
	return store.Has(blockKey)
}

// IsAuthorisedToRecord ensures only the WRKChain owner is recording hashes
func (k Keeper) IsAuthorisedToRecord(ctx sdk.Context, wrkchainId uint64, recorder sdk.AccAddress) bool {
	return recorder.Equals(k.GetWrkChainOwner(ctx, wrkchainId))
}

// GetWrkChainBlock Gets the entire WRKChain metadata struct for a wrkchainId
func (k Keeper) GetWrkChainBlock(ctx sdk.Context, wrkchainId uint64, height uint64) types.WrkChainBlock {
	store := ctx.KVStore(k.storeKey)

	if !k.IsWrkChainBlockRecorded(ctx, wrkchainId, height) {
		// return a new empty WrkChainBlock struct
		return types.WrkChainBlock{}
	}

	blockKey := types.WrkChainBlockKey(wrkchainId, height)

	bz := store.Get(blockKey)
	var wrkchainBlock types.WrkChainBlock
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchainBlock)
	return wrkchainBlock
}

// GetWrkChainBlockHashesIterator Gets an iterator over all WrkChain hashess in
// which the keys are the WrkChain Ids and the values are the WrkChainBlocks
func (k Keeper) GetWrkChainBlockHashesIterator(ctx sdk.Context, wrkchainID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WrkChainAllBlocksKey(wrkchainID))
}

// IterateWrkChainBlockHashes iterates over the all the hashes for a wrkchain and performs a callback function
func (k Keeper) IterateWrkChainBlockHashes(ctx sdk.Context, wrkchainID uint64, cb func(wrkChain types.WrkChainBlock) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WrkChainAllBlocksKey(wrkchainID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var wcb types.WrkChainBlock
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &wcb)

		if cb(wcb) {
			break
		}
	}
}

// IterateWrkChainBlockHashesReverse iterates over the all the hashes for a wrkchain in reverse order
// and performs a callback function
func (k Keeper) IterateWrkChainBlockHashesReverse(ctx sdk.Context, wrkchainID uint64, cb func(wrkChain types.WrkChainBlock) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.WrkChainAllBlocksKey(wrkchainID))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var wcb types.WrkChainBlock
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &wcb)

		if cb(wcb) {
			break
		}
	}
}

// GetAllWrkChainBlockHashes returns all the wrkchain's hashes from store
func (k Keeper) GetAllWrkChainBlockHashes(ctx sdk.Context, wrkchainID uint64) (wrkChainBlocks []types.WrkChainBlock) {
	k.IterateWrkChainBlockHashes(ctx, wrkchainID, func(wcb types.WrkChainBlock) bool {
		wrkChainBlocks = append(wrkChainBlocks, wcb)
		return false
	})
	return
}

func prependBlock(x []types.WrkChainBlockGenesisExport, y types.WrkChainBlockGenesisExport) []types.WrkChainBlockGenesisExport {
	x = append(x, y)
	copy(x[1:], x)
	x[0] = y
	return x
}

// GetAllWrkChainBlockHashesForGenesisExport returns all the wrkchain's hashes from store for export in an optimised
// format ready for genesis
func (k Keeper) GetAllWrkChainBlockHashesForGenesisExport(ctx sdk.Context, wrkchainID uint64) (wrkChainBlocks []types.WrkChainBlockGenesisExport) {
	count := 0
	k.IterateWrkChainBlockHashesReverse(ctx, wrkchainID, func(wcb types.WrkChainBlock) bool {
		wcbExp := types.WrkChainBlockGenesisExport{
			He: wcb.Height,
			Bh: wcb.Blockhash,
			Ph: wcb.Parenthash,
			H1: wcb.Hash1,
			H2: wcb.Hash2,
			H3: wcb.Hash3,
			St: wcb.SubTime,
		}
		wrkChainBlocks = prependBlock(wrkChainBlocks, wcbExp) // append(wrkChainBlocks, wcbExp)
		count = count + 1
		if count == types.MaxBlockSubmissionsKeepInState {
			return true
		}
		return false
	})
	return
}

// GetWrkChainBlockHashesFiltered retrieves wrkchains filtered by a given set of params which
// include pagination parameters along a moniker and owner address.
//
// NOTE: If no filters are provided, all WRKChains will be returned in paginated
// form.
func (k Keeper) GetWrkChainBlockHashesFiltered(ctx sdk.Context, wrkchainID uint64, params types.QueryWrkChainBlockParams) []types.WrkChainBlock {
	wrkChainHashes := k.GetAllWrkChainBlockHashes(ctx, wrkchainID)
	filteredWrkChainHashes := make([]types.WrkChainBlock, 0, len(wrkChainHashes))

	for _, wcb := range wrkChainHashes {
		matchMinHeight, matchMaxHeight, matchMinDate, matchMaxDate, matchHash := true, true, true, true, true

		if params.MinHeight > 0 {
			matchMinHeight = wcb.Height >= params.MinHeight
		}

		if params.MaxHeight > 0 {
			matchMaxHeight = wcb.Height <= params.MaxHeight
		}

		if params.MinDate > 0 {
			matchMinDate = wcb.SubTime >= params.MinDate
		}

		if params.MaxDate > 0 {
			matchMaxDate = wcb.SubTime <= params.MaxDate
		}

		if len(params.BlockHash) > 0 {
			matchHash = wcb.Blockhash == params.BlockHash
		}

		if matchMinHeight && matchMaxHeight && matchMinDate && matchMaxDate && matchHash {
			filteredWrkChainHashes = append(filteredWrkChainHashes, wcb)
		}
	}

	start, end := client.Paginate(len(filteredWrkChainHashes), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredWrkChainHashes = []types.WrkChainBlock{}
	} else {
		filteredWrkChainHashes = filteredWrkChainHashes[start:end]
	}

	return filteredWrkChainHashes
}

// RecordNewWrkchainHashes records a WRKChain block has for a registered WRKChain
func (k Keeper) RecordNewWrkchainHashes(
	ctx sdk.Context,
	wrkchainId uint64,
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
	owner sdk.AccAddress) error {

	logger := k.Logger(ctx)

	// we're only ever adding new WRKChain data, never updating existing. Handler will have checked if height has
	// previously been recorded.
	wrkchainBlock := types.WrkChainBlock{}

	wrkchainBlock.WrkchainId = wrkchainId
	wrkchainBlock.Height = height
	wrkchainBlock.Blockhash = blockHash
	wrkchainBlock.Parenthash = parentHash
	wrkchainBlock.Hash1 = hash1
	wrkchainBlock.Hash2 = hash2
	wrkchainBlock.Hash3 = hash3
	wrkchainBlock.Owner = owner.String()
	wrkchainBlock.SubTime = uint64(ctx.BlockTime().Unix())

	err := k.SetWrkChainBlock(ctx, wrkchainBlock)

	if err != nil {
		return err
	}

	err = k.SetLastBlock(ctx, wrkchainId, height)

	if err != nil {
		return err
	}

	err = k.SetNumBlocks(ctx, wrkchainId)

	if err != nil {
		return err
	}

	if !ctx.IsCheckTx() {
		logger.Debug("wrkchain block recorded", "id", wrkchainId, "height", height, "hash", blockHash)
	}

	return nil
}
