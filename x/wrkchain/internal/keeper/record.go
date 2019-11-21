package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)


// SetWrkChainBlock Sets the WrkChain Block struct for a sha256(wrkchainId + height)
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainBlock types.WrkChainBlock) sdk.Error {
	// must have an owner, WRKChain ID, Height and BlockHash
	if wrkchainBlock.Owner.Empty() || wrkchainBlock.WrkChainID == 0 || wrkchainBlock.Height == 0 || len(wrkchainBlock.BlockHash) == 0 {
		return sdk.ErrInternal("must include owner, id, height and hash")
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

// GetWrkChainBlockHashes Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
func (k Keeper) GetWrkChainBlockHashes(ctx sdk.Context, wrkchainId uint64) []types.WrkChainBlock {
	store := ctx.KVStore(k.storeKey)
	var wrkchainBlocks []types.WrkChainBlock

	iterator := sdk.KVStorePrefixIterator(store, types.WrkChainAllBlocksKey(wrkchainId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var block types.WrkChainBlock
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &block)
		wrkchainBlocks = append(wrkchainBlocks, block)
	}

	return wrkchainBlocks
}

// RecordWrkchainHashes records a WRKChain block has for a registered wRKchain
func (k Keeper) RecordWrkchainHashes(
	ctx sdk.Context,
	wrkchainId uint64,
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
	owner sdk.AccAddress) sdk.Error {

	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// can't record hashes if WRKChain isn't registered
		return types.ErrWrkChainDoesNotExist(k.codespace, "WRKChain does not exist")
	}

	wrkchain := k.GetWrkChain(ctx, wrkchainId)

	if k.IsWrkChainBlockRecorded(ctx, wrkchain.WrkChainID, height) {
		return types.ErrWrkChainBlockAlreadyRecorded(k.codespace, "Block hashes already recorded for this height")
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
	wrkchainBlock.SubmitHeight = ctx.BlockHeight()

	err := k.SetWrkChainBlock(ctx, wrkchainBlock)

	if err != nil {
		return err
	}

	err = k.SetLastBlock(ctx, wrkchain.WrkChainID, height)

	if err != nil {
		return err
	}

	return nil
}
