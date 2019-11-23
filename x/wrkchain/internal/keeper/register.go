package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

//__WRKCHAIN_ID_________________________________________________________

// GetHighestWrkChainID gets the highest WRKChain ID
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

// SetHighestWrkChainID sets the new highest WRKChain ID to the store
func (k Keeper) SetHighestWrkChainID(ctx sdk.Context, wrkChainID uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	wrkChainIDbz := types.GetWrkChainIDBytes(wrkChainID)
	store.Set(types.HighestWrkChainIDKey, wrkChainIDbz)
}

//__WRKCHAINS___________________________________________________________

// SetWrkChain Sets the WrkChain metadata struct for a wrkchainId
func (k Keeper) SetWrkChain(ctx sdk.Context, wrkchain types.WrkChain) sdk.Error {
	// must have an owner
	if wrkchain.Owner.Empty() {
		return sdk.ErrInternal("unable to set WRKChain - must have an owner")
	}

	//must have an ID
	if wrkchain.WrkChainID == 0 {
		return sdk.ErrInternal("unable to set WRKChain - id must be positive non-zero")
	}

	//must have a moniker
	if len(wrkchain.Moniker) == 0 {
		return sdk.ErrInternal("unable to set WRKChain - must have a moniker")
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
	if blockNum > wrkchain.LastBlock {
		wrkchain.LastBlock = blockNum
		return k.SetWrkChain(ctx, wrkchain)
	}
	return nil
}

// GetWrkChain Gets the entire WRKChain metadata struct for a wrkchainId
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

// IsWrkChainRegistered Checks if the WrkChain is present in the store or not
func (k Keeper) IsWrkChainRegistered(ctx sdk.Context, wrkchainId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.WrkChainKey(wrkchainId))
}

// GetWrkChainsIterator Get an iterator over all WrkChains in which the keys are the WrkChain Ids and the values are the WrkChains
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
func (k Keeper) GetAllWrkChains(ctx sdk.Context) (wrkChains types.WrkChains) {
	k.IterateWrkChains(ctx, func(wc types.WrkChain) bool {
		wrkChains = append(wrkChains, wc)
		return false
	})
	return
}

// GetWrkChainsFiltered retrieves wrkchains filtered by a given set of params which
// include pagination parameters along a moniker and owner address.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
// form.
func (k Keeper) GetWrkChainsFiltered(ctx sdk.Context, params types.QueryWrkChainParams) []types.WrkChain {
	wrkChains := k.GetAllWrkChains(ctx)
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

// RegisterWrkChain registers a WRKChain in the store
func (k Keeper) RegisterWrkChain(ctx sdk.Context, moniker string, wrkchainName string, genesisHash string, owner sdk.AccAddress) (uint64, sdk.Error) {

	//must have a moniker
	if len(moniker) == 0 {
		return 0, sdk.ErrInternal("unable to set WRKChain - must have a moniker")
	}

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
	wrkchain.RegisterHeight = ctx.BlockHeight()
	wrkchain.RegisterTime = ctx.BlockTime().Unix()

	err = k.SetWrkChain(ctx, wrkchain)
	if err != nil {
		return 0, err
	}

	k.SetHighestWrkChainID(ctx, wrkChainId+1)

	return wrkChainId, nil
}
