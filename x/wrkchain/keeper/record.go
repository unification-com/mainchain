package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// SetWrkChainBlock Sets the WrkChain Block struct for a wrkchainId & height
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainId uint64, wrkchainBlock types.WrkChainBlock) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainBlockKey(wrkchainId, wrkchainBlock.Height), k.cdc.MustMarshalBinaryBare(&wrkchainBlock))

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
func (k Keeper) GetWrkChainBlock(ctx sdk.Context, wrkchainId uint64, height uint64) (types.WrkChainBlock, bool) {
	store := ctx.KVStore(k.storeKey)

	if !k.IsWrkChainBlockRecorded(ctx, wrkchainId, height) {
		// return a new empty WrkChainBlock struct
		return types.WrkChainBlock{}, false
	}

	blockKey := types.WrkChainBlockKey(wrkchainId, height)

	bz := store.Get(blockKey)
	var wrkchainBlock types.WrkChainBlock
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchainBlock)
	return wrkchainBlock, true
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

func prependBlock(x types.WrkChainBlockGenesisExports, y types.WrkChainBlockGenesisExport) types.WrkChainBlockGenesisExports {
	x = append(x, y)
	copy(x[1:], x)
	x[0] = y
	return x
}

// GetAllWrkChainBlockHashesForGenesisExport returns all the wrkchain's hashes from store for export in an optimised
// format ready for genesis
func (k Keeper) GetAllWrkChainBlockHashesForGenesisExport(ctx sdk.Context, wrkchainID uint64) (wrkChainBlocks types.WrkChainBlockGenesisExports) {
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
		return count == types.MaxBlockSubmissionsKeepInState
	})
	return
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
	hash3 string) error {

	logger := k.Logger(ctx)

	// we're only ever adding new WRKChain data, never updating existing. Handler will have checked if height has
	// previously been recorded.
	wrkchainBlock := types.WrkChainBlock{
		Height:     height,
		Blockhash:  blockHash,
		Parenthash: parentHash,
		Hash1:      hash1,
		Hash2:      hash2,
		Hash3:      hash3,
		SubTime:    uint64(ctx.BlockTime().Unix()),
	}

	err := k.SetWrkChainBlock(ctx, wrkchainId, wrkchainBlock)

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
