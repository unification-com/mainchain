package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

//__WRKCHAIN_ID_________________________________________________________

// GetHighestWrkChainID gets the highest WRKChain ID
func (k Keeper) GetHighestWrkChainID(ctx sdk.Context) (wrkChainID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestWrkChainIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial wrkchain ID hasn't been set")
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
func (k Keeper) SetWrkChain(ctx sdk.Context, wrkchain types.WrkChain) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.WrkChainKey(wrkchain.WrkchainId), k.cdc.MustMarshalBinaryBare(&wrkchain))

	return nil
}

func (k Keeper) SetNumBlocks(ctx sdk.Context, wrkchainId uint64) error {
	wrkchain, found := k.GetWrkChain(ctx, wrkchainId)
	if !found {
		// doesn't exist. Don't update
		return sdkerrors.Wrap(types.ErrWrkChainDoesNotExist, fmt.Sprintf("WRKChain %v does not exist", wrkchainId))
	}

	wrkchain.NumBlocks = wrkchain.NumBlocks + 1

	return k.SetWrkChain(ctx, wrkchain)
}

// SetLastBlock - sets the last block number submitted
func (k Keeper) SetLastBlock(ctx sdk.Context, wrkchainId uint64, blockNum uint64) error {
	wrkchain, found := k.GetWrkChain(ctx, wrkchainId)
	if !found {
		// doesn't exist. Don't update
		return sdkerrors.Wrap(types.ErrWrkChainDoesNotExist, fmt.Sprintf("WRKChain %v does not exist", wrkchainId))
	}

	if blockNum > wrkchain.Lastblock {
		wrkchain.Lastblock = blockNum
		return k.SetWrkChain(ctx, wrkchain)
	}

	return nil
}

// GetWrkChain Gets the entire WRKChain metadata struct for a wrkchainId
func (k Keeper) GetWrkChain(ctx sdk.Context, wrkchainId uint64) (types.WrkChain, bool) {
	store := ctx.KVStore(k.storeKey)
	if !k.IsWrkChainRegistered(ctx, wrkchainId) {
		// return a new empty WrkChain struct
		return types.WrkChain{}, false
	}
	bz := store.Get(types.WrkChainKey(wrkchainId))
	var wrkchain types.WrkChain
	k.cdc.MustUnmarshalBinaryBare(bz, &wrkchain)
	return wrkchain, true
}

// GetWrkChainOwner - get the current owner of a WRKChain
func (k Keeper) GetWrkChainOwner(ctx sdk.Context, wrkchainId uint64) sdk.AccAddress {
	wrkchain, found := k.GetWrkChain(ctx, wrkchainId)
	if !found {
		return sdk.AccAddress{}
	}
	accAddr, accErr := sdk.AccAddressFromBech32(wrkchain.Owner)
	if accErr != nil {
		return sdk.AccAddress{}
	}
	return accAddr
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
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &wc)

		if cb(wc) {
			break
		}
	}
}

// GetAllWrkChains returns all the registered wrkchain metadata from store
func (k Keeper) GetAllWrkChains(ctx sdk.Context) (wrkChains []types.WrkChain) {
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
func (k Keeper) GetWrkChainsFiltered(ctx sdk.Context, params types.QueryWrkChainsFilteredRequest) []types.WrkChain {
	wrkChains := k.GetAllWrkChains(ctx)
	filteredWrkChains := make([]types.WrkChain, 0, len(wrkChains))

	for _, wc := range wrkChains {
		matchMoniker, matchOwner := true, true

		if len(params.Moniker) > 0 {
			matchMoniker = wc.Moniker == params.Moniker
		}

		if len(params.Owner) > 0 {
			matchOwner = wc.Owner == params.Owner
		}

		if matchMoniker && matchOwner {
			filteredWrkChains = append(filteredWrkChains, wc)
		}
	}

	// todo - proper pagination
	start, end := client.Paginate(len(filteredWrkChains), 1, 100, 100)
	if start < 0 || end < 0 {
		filteredWrkChains = []types.WrkChain{}
	} else {
		filteredWrkChains = filteredWrkChains[start:end]
	}

	return filteredWrkChains
}

// RegisterNewWrkChain registers a WRKChain in the store
func (k Keeper) RegisterNewWrkChain(ctx sdk.Context, moniker string, wrkchainName string, genesisHash string, baseType string, owner sdk.AccAddress) (uint64, error) {

	wrkChainId, err := k.GetHighestWrkChainID(ctx)
	if err != nil {
		return 0, err
	}

	wrkchain := types.WrkChain{}

	wrkchain.WrkchainId = wrkChainId
	wrkchain.Moniker = moniker
	wrkchain.Lastblock = 0
	wrkchain.NumBlocks = 0
	wrkchain.Owner = owner.String()
	wrkchain.Name = wrkchainName
	wrkchain.Genesis = genesisHash
	wrkchain.Type = baseType
	wrkchain.RegTime = uint64(ctx.BlockTime().Unix())

	err = k.SetWrkChain(ctx, wrkchain)
	if err != nil {
		return 0, err
	}

	k.SetHighestWrkChainID(ctx, wrkChainId+1)

	return wrkChainId, nil
}
