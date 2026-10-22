package keeper

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/wrkchain/types"
)

//__WRKCHAIN_ID_________________________________________________________

// GetHighestWrkChainID gets the highest WRKChain ID
func (k Keeper) GetHighestWrkChainID(ctx sdk.Context) (wrkChainID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestWrkChainIDKey)
	if bz == nil {
		return 0, errorsmod.Wrap(types.ErrInvalidGenesis, "initial wrkchain ID hasn't been set")
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
	store.Set(types.WrkChainKey(wrkchain.WrkchainId), k.cdc.MustMarshal(&wrkchain))

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
	if err := k.cdc.Unmarshal(bz, &wrkchain); err != nil {
		k.Logger(ctx).Error("corrupt wrkchain entry — treating as absent",
			"wrkchain_id", wrkchainId, "err", err)
		return types.WrkChain{}, false
	}
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

// IterateWrkChains iterates over the all the wrkchain metadata and performs a callback function
func (k Keeper) IterateWrkChains(ctx sdk.Context, cb func(wrkChain types.WrkChain) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.RegisteredWrkChainPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var wc types.WrkChain
		if err := k.cdc.Unmarshal(iterator.Value(), &wc); err != nil {
			k.Logger(ctx).Error("skipping corrupt wrkchain entry during iteration",
				"err", err)
			continue
		}

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

// CountWrkChainsForOwner counts the number of wrkchains registered to the
// given owner. Enforced against types.MaxWrkChainsPerOwner at registration
// time. O(total wrkchains) — acceptable while overall count stays modest.
func (k Keeper) CountWrkChainsForOwner(ctx sdk.Context, owner sdk.AccAddress) int {
	count := 0
	ownerStr := owner.String()
	k.IterateWrkChains(ctx, func(wc types.WrkChain) bool {
		if wc.Owner == ownerStr {
			count++
		}
		return false
	})
	return count
}

// Filtered wrkchain queries are served directly by the gRPC handler in
// grpc_query.go using query.FilteredPaginate against the store; no separate
// keeper helper is needed.

func (k Keeper) GetWrkChainStorageLimit(ctx sdk.Context, wrkchainId uint64) (types.WrkChainStorageLimit, bool) {
	store := ctx.KVStore(k.storeKey)
	if !k.HasWrkChainStorageLimit(ctx, wrkchainId) {
		return types.WrkChainStorageLimit{
			WrkchainId:   wrkchainId,
			InStateLimit: types.DefaultStorageLimit,
		}, false
	}

	storageKey := types.WrkChainStorageLimitKey(wrkchainId)
	bz := store.Get(storageKey)
	var storage types.WrkChainStorageLimit
	if err := k.cdc.Unmarshal(bz, &storage); err != nil {
		k.Logger(ctx).Error("corrupt wrkchain storage limit entry — treating as absent",
			"wrkchain_id", wrkchainId, "err", err)
		return types.WrkChainStorageLimit{
			WrkchainId:   wrkchainId,
			InStateLimit: types.DefaultStorageLimit,
		}, false
	}
	return storage, true
}

func (k Keeper) HasWrkChainStorageLimit(ctx sdk.Context, wrkchainId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	storageKey := types.WrkChainStorageLimitKey(wrkchainId)
	return store.Has(storageKey)
}

func (k Keeper) SetWrkChainStorageLimit(ctx sdk.Context, wrkchainId, limit uint64) error {

	store := ctx.KVStore(k.storeKey)
	storageLimit := types.WrkChainStorageLimit{
		WrkchainId:   wrkchainId,
		InStateLimit: limit,
	}
	store.Set(types.WrkChainStorageLimitKey(wrkchainId), k.cdc.MustMarshal(&storageLimit))

	return nil
}

func (k Keeper) IncreaseInStateStorage(ctx sdk.Context, wrkchainId, amount uint64) error {
	wrkchainStorage, _ := k.GetWrkChainStorageLimit(ctx, wrkchainId)
	newInStateLimit := wrkchainStorage.InStateLimit + amount
	err := k.SetWrkChainStorageLimit(ctx, wrkchainId, newInStateLimit)

	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetMaxPurchasableSlots(ctx sdk.Context, wrkchainId uint64) uint64 {
	wrkchainStorage, found := k.GetWrkChainStorageLimit(ctx, wrkchainId)
	if !found {
		return 0
	}

	maxStorageLimit := k.GetParamMaxStorageLimit(ctx)

	if wrkchainStorage.InStateLimit >= maxStorageLimit {
		return 0
	}

	return maxStorageLimit - wrkchainStorage.InStateLimit
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
	wrkchain.BaseType = baseType
	wrkchain.RegTime = uint64(ctx.BlockTime().Unix())

	err = k.SetWrkChain(ctx, wrkchain)
	if err != nil {
		return 0, err
	}

	err = k.SetWrkChainStorageLimit(ctx, wrkChainId, k.GetParamDefaultStorageLimit(ctx))

	if err != nil {
		return 0, err
	}

	k.SetHighestWrkChainID(ctx, wrkChainId+1)

	return wrkChainId, nil
}
