package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

// SetWrkChainBlock Sets the WrkChain Block struct for a wrkchainId & height
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainBlock types.WrkChainBlock) error {
	// must have an owner, WRKChain ID, Height and BlockHash
	if wrkchainBlock.Owner.Empty() || wrkchainBlock.WrkChainID == 0 || wrkchainBlock.Height == 0 || len(wrkchainBlock.BlockHash) == 0 {
		return sdkerrors.Wrap(types.ErrMissingData, "must include owner, id, height and hash")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainBlockKey(wrkchainBlock.WrkChainID, wrkchainBlock.Height), k.cdc.MustMarshalBinaryLengthPrefixed(wrkchainBlock))

	return nil
}

// QuickCheckHeightIsRecorded Checks if the given height can be recorded
func (k Keeper) QuickCheckHeightIsRecorded(ctx sdk.Context, wrkchainId uint64, height uint64) bool {

	wrkchain := k.GetWrkChain(ctx, wrkchainId)

	// only check if height being submitted is <= last recorded height.
	// Otherwise, no need to check entire db
	if height <= wrkchain.LastBlock && height > 0 {
		if height == wrkchain.LastBlock {
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
		return types.NewWrkChainBlock()
	}

	blockKey := types.WrkChainBlockKey(wrkchainId, height)

	bz := store.Get(blockKey)
	var wrkchainBlock types.WrkChainBlock
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &wrkchainBlock)
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
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &wcb)

		if cb(wcb) {
			break
		}
	}
}

// GetAllWrkChainBlockHashes returns all the wrkchain's hashes from store
func (k Keeper) GetAllWrkChainBlockHashes(ctx sdk.Context, wrkchainID uint64) (wrkChainBlocks types.WrkChainBlocks) {
	k.IterateWrkChainBlockHashes(ctx, wrkchainID, func(wcb types.WrkChainBlock) bool {
		wrkChainBlocks = append(wrkChainBlocks, wcb)
		return false
	})
	return
}

// GetAllWrkChainBlockHashesForGenesisExport returns all the wrkchain's hashes from store for export in an optimised
// format ready for genesis
func (k Keeper) GetAllWrkChainBlockHashesForGenesisExport(ctx sdk.Context, wrkchainID uint64) (wrkChainBlocks types.WrkChainBlocksGenesisExport) {
	k.IterateWrkChainBlockHashes(ctx, wrkchainID, func(wcb types.WrkChainBlock) bool {
		wcbExp := types.WrkChainBlockGenesisExport {
			Height: wcb.Height,
			BlockHash: wcb.BlockHash,
			ParentHash: wcb.ParentHash,
			Hash1: wcb.Hash1,
			Hash2: wcb.Hash2,
			Hash3: wcb.Hash3,
			SubmitTime: wcb.SubmitTime,
		}
		wrkChainBlocks = append(wrkChainBlocks, wcbExp)
		return false
	})
	return
}

// GetWrkChainsFiltered retrieves wrkchains filtered by a given set of params which
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
			matchMinDate = wcb.SubmitTime >= int64(params.MinDate)
		}

		if params.MaxDate > 0 {
			matchMaxDate = wcb.SubmitTime <= int64(params.MaxDate)
		}

		if len(params.BlockHash) > 0 {
			matchHash = wcb.BlockHash == params.BlockHash
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

// RecordWrkchainHashes records a WRKChain block has for a registered WRKChain
func (k Keeper) RecordWrkchainHashes(
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
	wrkchainBlock := types.NewWrkChainBlock()

	wrkchainBlock.WrkChainID = wrkchainId
	wrkchainBlock.Height = height
	wrkchainBlock.BlockHash = blockHash
	wrkchainBlock.ParentHash = parentHash
	wrkchainBlock.Hash1 = hash1
	wrkchainBlock.Hash2 = hash2
	wrkchainBlock.Hash3 = hash3
	wrkchainBlock.Owner = owner
	wrkchainBlock.SubmitTime = ctx.BlockTime().Unix()

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
		logger.Debug("wrkchain hashes recorded", "wcid", wrkchainId, "height", height, "hash", blockHash, "owner", owner.String())
	}
	return nil
}
