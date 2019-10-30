package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
	}
}

// Sets the WrkChain metadata struct for a wrkchainId
func (k Keeper) SetWrkChain(ctx sdk.Context, wrkchainId string, wrkchain types.WrkChain) {
	// must have an owner
	if wrkchain.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWrkChainStoreKey(wrkchainId), k.cdc.MustMarshalBinaryBare(wrkchain))
}

// SetLastBlock - sets the last block number submitted
func (k Keeper) SetLastBlock(ctx sdk.Context, wrkchainId string, blockNum uint64) {
	wrkchain := k.GetWrkChain(ctx, wrkchainId)
	if wrkchain.Owner.Empty() {
		// doesn't exist. Don't update
		return
	}
	if blockNum <= wrkchain.LastBlock {
		return
	}
	wrkchain.LastBlock = blockNum
	k.SetWrkChain(ctx, wrkchainId, wrkchain)
}

// Gets the entire WRKChain metadata struct for a wrkchainId
func (k Keeper) GetWrkChain(ctx sdk.Context, wrkchainId string) types.WrkChain {
	store := ctx.KVStore(k.storeKey)
	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// return a new empty WrkChain struct
		return types.NewWrkChain()
	}
	bz := store.Get(types.GetWrkChainStoreKey(wrkchainId))
	var wrkchain types.WrkChain
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchain)
	return wrkchain
}

// GetWrkChainOwner - get the current owner of a WRKChain
func (k Keeper) GetWrkChainOwner(ctx sdk.Context, wrkchainId string) sdk.AccAddress {
	return k.GetWrkChain(ctx, wrkchainId).Owner
}

// Check if the WrkChain is present in the store or not
func (k Keeper) IsWrkChainRegistered(ctx sdk.Context, wrkchainId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWrkChainStoreKey(wrkchainId))
}

// Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
func (k Keeper) GetWrkChainsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RegisteredWrkChainPrefix)
}

func (k Keeper) RegisterWrkChain(ctx sdk.Context, wrkchainId string, wrkchainName string, genesisHash string, owner sdk.AccAddress) {
	wrkchain := k.GetWrkChain(ctx, wrkchainId)

	wrkchain.WrkChainID = wrkchainId
	wrkchain.LastBlock = 0
	wrkchain.Owner = owner
	wrkchain.WrkChainName = wrkchainName
	wrkchain.GenesisHash = genesisHash
	k.SetWrkChain(ctx, wrkchainId, wrkchain)
}

// Record hashes funcs

// Sets the WrkChain Block struct for a sha256(wrkchainId + height)
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainBlockKey []byte, wrkchainBlock types.WrkChainBlock) {
	// must have an owner, WRKChain ID, Height and BlockHash
	if wrkchainBlock.Owner.Empty() || len(wrkchainBlock.WrkChainID) == 0 || wrkchainBlock.Height == 0 || len(wrkchainBlock.BlockHash) == 0 {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(wrkchainBlockKey, k.cdc.MustMarshalBinaryBare(wrkchainBlock))
}

// Check if the WrkChainBlock is present in the store or not
func (k Keeper) IsWrkChainBlockRecorded(ctx sdk.Context, wrkchainId string, height uint64) bool {
	store := ctx.KVStore(k.storeKey)
	blockKey := types.GetWrkChainBlockHashStoreKey(wrkchainId, height)
	return store.Has(blockKey)
}

// IsAuthorisedToRecord ensures only the WRKChain owner is recording hashes
func (k Keeper) IsAuthorisedToRecord(ctx sdk.Context, wrkchainId string, recorder sdk.AccAddress) bool {
	return recorder.Equals(k.GetWrkChainOwner(ctx, wrkchainId))
}

// Gets the entire WRKChain metadata struct for a wrkchainId
func (k Keeper) GetWrkChainBlock(ctx sdk.Context, wrkchainId string, height uint64) types.WrkChainBlock {
	store := ctx.KVStore(k.storeKey)

	if !k.IsWrkChainBlockRecorded(ctx, wrkchainId, height) {
		// return a new empty WrkChainBlock struct
		return types.NewWrkChainBlock()
	}

	blockKey := types.GetWrkChainBlockHashStoreKey(wrkchainId, height)

	bz := store.Get(blockKey)
	var wrkchainBlock types.WrkChainBlock
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchainBlock)
	return wrkchainBlock
}

// Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
func (k Keeper) GetWrkChainBlockHashes(ctx sdk.Context, wrkchainId string) []types.WrkChainBlock {
	store := ctx.KVStore(k.storeKey)
	var wrkchainBlocks []types.WrkChainBlock

	iterator := sdk.KVStorePrefixIterator(store, types.GetWrkChainBlockHashStoreKeyPrefix(wrkchainId))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var block types.WrkChainBlock
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &block)
		wrkchainBlocks = append(wrkchainBlocks, block)
	}

	return wrkchainBlocks
}

func (k Keeper) RecordWrkchainHashes(
	ctx sdk.Context,
	wrkchainId string,
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
	owner sdk.AccAddress) {

	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// can't record hashes if WRKChain isn't registered
		return
	}

	blockKey := types.GetWrkChainBlockHashStoreKey(wrkchainId, height)
	wrkchain := k.GetWrkChain(ctx, wrkchainId)
	wrkchainBlock := k.GetWrkChainBlock(ctx, wrkchain.WrkChainID, height)

	wrkchainBlock.WrkChainID = wrkchain.WrkChainID
	wrkchainBlock.Height = height
	wrkchainBlock.BlockHash = blockHash
	wrkchainBlock.ParentHash = parentHash
	wrkchainBlock.Hash1 = hash1
	wrkchainBlock.Hash2 = hash2
	wrkchainBlock.Hash3 = hash3
	wrkchainBlock.Owner = owner
	wrkchainBlock.SubmitTime = uint64(time.Now().Unix()) // todo - change to block time?

	k.SetWrkChainBlock(ctx, blockKey, wrkchainBlock)
	k.SetLastBlock(ctx, wrkchain.WrkChainID, height)

}
