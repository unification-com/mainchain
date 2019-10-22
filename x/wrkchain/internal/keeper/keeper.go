package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Sets the WrkChain metadata struct for a wrkchainId
func (k Keeper) SetWrkChain(ctx sdk.Context, wrkchainId string, wrkchain types.WrkChain) {
	// must have an owner
	if wrkchain.Owner.Empty() {
		return
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(wrkchainId), k.cdc.MustMarshalBinaryBare(wrkchain))
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

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWrkChain(ctx sdk.Context, wrkchainId string) types.WrkChain {
	store := ctx.KVStore(k.storeKey)
	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// return a new empty WrkChain struct
		return types.NewWrkChain()
	}
	bz := store.Get([]byte(wrkchainId))
	var wrkchain types.WrkChain
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchain)
	return wrkchain
}

// Check if the WrkChain is present in the store or not
func (k Keeper) IsWrkChainRegistered(ctx sdk.Context, wrkchainId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(wrkchainId))
}

// Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
func (k Keeper) GetWrkChainsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
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
