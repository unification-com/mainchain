package keeper

import (
	"fmt"
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

// GetWrkChainsFiltered retrieves wrkchains filtered by a given set of params which
// include pagination parameters along a moniker and owner address.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
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

	if len(blockHash) > 66 {
		return sdkerrors.Wrap(types.ErrContentTooLarge, "block hash too big. 66 character limit")
	}
	if len(parentHash) > 66 {
		return sdkerrors.Wrap(types.ErrContentTooLarge, "parent hash too big. 66 character limit")
	}
	if len(hash1) > 66 {
		return sdkerrors.Wrap(types.ErrContentTooLarge, "hash1 too big. 66 character limit")
	}
	if len(hash2) > 66 {
		return sdkerrors.Wrap(types.ErrContentTooLarge, "hash2 too big. 66 character limit")
	}
	if len(hash3) > 66 {
		return sdkerrors.Wrap(types.ErrContentTooLarge, "hash3 too big. 66 character limit")
	}

	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// can't record hashes if WRKChain isn't registered
		return sdkerrors.Wrap(types.ErrWrkChainDoesNotExist, fmt.Sprintf("WRKChain %v does not exist", wrkchainId))
	}

	wrkchain := k.GetWrkChain(ctx, wrkchainId)

	if !k.IsAuthorisedToRecord(ctx, wrkchain.WrkChainID, owner) {
		return sdkerrors.Wrap(types.ErrNotWrkChainOwner, "not authorised to record hashes for this wrkchain")
	}

	if k.IsWrkChainBlockRecorded(ctx, wrkchain.WrkChainID, height) {
		return sdkerrors.Wrap(types.ErrWrkChainBlockAlreadyRecorded, "Block hashes already recorded for this height")
	}

	wrkchainBlock := k.GetWrkChainBlock(ctx, wrkchain.WrkChainID, height)

	wrkchainBlock.WrkChainID = wrkchain.WrkChainID
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

	err = k.SetLastBlock(ctx, wrkchain.WrkChainID, height)

	if err != nil {
		return err
	}

	err = k.SetNumBlocks(ctx, wrkchain.WrkChainID)

	if err != nil {
		return err
	}

	if !ctx.IsCheckTx() {
		logger.Debug("wrkchain hashes recorded", "wcid", wrkchain.WrkChainID, "height", height, "hash", blockHash, "owner", owner.String())
	}
	return nil
}
