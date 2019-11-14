package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/tendermint/tendermint/libs/log"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	codespace sdk.CodespaceType
	cdc       *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the wrkchain Keeper
func NewKeeper(storeKey sdk.StoreKey, codespace sdk.CodespaceType, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:  storeKey,
		codespace: codespace,
		cdc:       cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) Cdc() *codec.Codec {
	return k.cdc
}

//__WRKCHAIN_ID_________________________________________________________

// GetHighestPurchaseOrderID gets the highest purchase order ID
func (k Keeper) GetHighestWrkChainID(ctx sdk.Context) (wrkChainID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestWrkChainIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(k.codespace, "initial wrkchain ID hasn't been set")
	}
	// convert from bytes to uint64
	wrkChainID = types.GetWrkChainIDFromBytes(bz)
	return wrkChainID, nil
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetHighestWrkChainID(ctx sdk.Context, wrkChainID uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	wrkChainIDbz := types.GetWrkChainIDBytes(wrkChainID)
	store.Set(types.HighestWrkChainIDKey, wrkChainIDbz)
}

// Sets the WrkChain metadata struct for a wrkchainId
func (k Keeper) SetWrkChain(ctx sdk.Context, wrkchain types.WrkChain) sdk.Error {
	// must have an owner
	if wrkchain.Owner.Empty() {
		return sdk.ErrInternal("unable to register WRKChain - must have an owner")
	}

	//must have an ID
	if wrkchain.WrkChainID == 0 {
		return sdk.ErrInternal("unable to register WRKChain - id must be positive non-zero")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainKey(wrkchain.WrkChainID), k.cdc.MustMarshalBinaryLengthPrefixed(wrkchain))

	return nil
}

// SetLastBlock - sets the last block number submitted
func (k Keeper) SetLastBlock(ctx sdk.Context, wrkchainId uint64, blockNum uint64) sdk.Error {
	wrkchain := k.GetWrkChain(ctx, wrkchainId)
	if wrkchain.Owner.Empty() {
		// doesn't exist. Don't update
		return types.ErrWrkChainDoesNotExist(k.codespace, "WRKChain does not exist")
	}
	if blockNum <= wrkchain.LastBlock {
		return sdk.ErrInternal("this block must be greater than last block")
	}
	wrkchain.LastBlock = blockNum
	return k.SetWrkChain(ctx, wrkchain)
}

// Gets the entire WRKChain metadata struct for a wrkchainId
func (k Keeper) GetWrkChain(ctx sdk.Context, wrkchainId uint64) types.WrkChain {
	store := ctx.KVStore(k.storeKey)
	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// return a new empty WrkChain struct
		return types.NewWrkChain()
	}
	bz := store.Get(types.WrkChainKey(wrkchainId))
	var wrkchain types.WrkChain
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &wrkchain)
	return wrkchain
}

// GetWrkChainOwner - get the current owner of a WRKChain
func (k Keeper) GetWrkChainOwner(ctx sdk.Context, wrkchainId uint64) sdk.AccAddress {
	return k.GetWrkChain(ctx, wrkchainId).Owner
}

// Check if the WrkChain is present in the store or not
func (k Keeper) IsWrkChainRegistered(ctx sdk.Context, wrkchainId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.WrkChainKey(wrkchainId))
}

// Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
func (k Keeper) GetWrkChainsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.RegisteredWrkChainPrefix)
}

// IterateWrkChains iterates over the all the wrkchain metadata and performs a callback function
func (k Keeper) IterateWrkChains(ctx sdk.Context, cb func(wrkChain types.WrkChain) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RegisteredWrkChainPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var wc types.WrkChain
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &wc)

		if cb(wc) {
			break
		}
	}
}

// GetAllWrkChains returns all the registered wrkchain metadata from store
func (keeper Keeper) GetAllWrkChains(ctx sdk.Context) (wrkChains types.WrkChains) {
	keeper.IterateWrkChains(ctx, func(wc types.WrkChain) bool {
		wrkChains = append(wrkChains, wc)
		return false
	})
	return
}

// GetPurchaseOrdersFiltered retrieves purchase orders filtered by a given set of params which
// include pagination parameters along a purchase order status.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
// form.
func (keeper Keeper) GetWrkChainsFiltered(ctx sdk.Context, params types.QueryWrkChainParams) []types.WrkChain {
	wrkChains := keeper.GetAllWrkChains(ctx)
	filteredWrkChains := make([]types.WrkChain, 0, len(wrkChains))

	for _, wc := range wrkChains {
		matchMoniker, matchOwner := true, true

		if len(params.Moniker) > 0 {
			matchMoniker = wc.Moniker == params.Moniker
		}

		if len(params.Owner) > 0 {
			matchOwner = wc.Owner.String() == params.Owner.String()
		}

		if matchMoniker && matchOwner {
			filteredWrkChains = append(filteredWrkChains, wc)
		}
	}

	start, end := client.Paginate(len(filteredWrkChains), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredWrkChains = []types.WrkChain{}
	} else {
		filteredWrkChains = filteredWrkChains[start:end]
	}

	return filteredWrkChains
}

func (k Keeper) RegisterWrkChain(ctx sdk.Context, moniker string, wrkchainName string, genesisHash string, owner sdk.AccAddress) (uint64, sdk.Error) {

	wrkChainId, err := k.GetHighestWrkChainID(ctx)
	if err != nil {
		return 0, err
	}

	params := types.NewQueryWrkChainParams(1, 1, moniker, sdk.AccAddress{})
	wrkChains := k.GetWrkChainsFiltered(ctx, params)

	if (len(wrkChains)) > 0 {
		errMsg := fmt.Sprintf("wrkchain already registered with moniker '%s' - id: %d, owner: %s", moniker, wrkChains[0].WrkChainID, wrkChains[0].Owner)
		return 0, types.ErrWrkChainAlreadyRegistered(k.codespace, errMsg)
	}

	wrkchain := k.GetWrkChain(ctx, wrkChainId)

	wrkchain.WrkChainID = wrkChainId
	wrkchain.Moniker = moniker
	wrkchain.LastBlock = 0
	wrkchain.Owner = owner
	wrkchain.Name = wrkchainName
	wrkchain.GenesisHash = genesisHash

	err = k.SetWrkChain(ctx, wrkchain)
	if err != nil {
		return 0, err
	}

	k.SetHighestWrkChainID(ctx, wrkChainId+1)

	return wrkChainId, nil
}

// Record hashes funcs

// Sets the WrkChain Block struct for a sha256(wrkchainId + height)
func (k Keeper) SetWrkChainBlock(ctx sdk.Context, wrkchainBlock types.WrkChainBlock) sdk.Error {
	// must have an owner, WRKChain ID, Height and BlockHash
	if wrkchainBlock.Owner.Empty() || wrkchainBlock.WrkChainID == 0 || wrkchainBlock.Height == 0 || len(wrkchainBlock.BlockHash) == 0 {
		return sdk.ErrInternal("must include owner, id, height and hash")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainBlockKey(wrkchainBlock.WrkChainID, wrkchainBlock.Height), k.cdc.MustMarshalBinaryLengthPrefixed(wrkchainBlock))

	return nil
}

// Check if the WrkChainBlock is present in the store or not
func (k Keeper) IsWrkChainBlockRecorded(ctx sdk.Context, wrkchainId uint64, height uint64) bool {
	store := ctx.KVStore(k.storeKey)
	blockKey := types.WrkChainBlockKey(wrkchainId, height)
	return store.Has(blockKey)
}

// IsAuthorisedToRecord ensures only the WRKChain owner is recording hashes
func (k Keeper) IsAuthorisedToRecord(ctx sdk.Context, wrkchainId uint64, recorder sdk.AccAddress) bool {
	return recorder.Equals(k.GetWrkChainOwner(ctx, wrkchainId))
}

// Gets the entire WRKChain metadata struct for a wrkchainId
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

// Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
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
