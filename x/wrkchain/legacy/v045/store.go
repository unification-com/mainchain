package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v040 "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

func migrateParams(ctx sdk.Context, paramsSubspace paramstypes.Subspace) {
	// Add the new module params
	var oldFeeReeRegParams uint64
	paramsSubspace.Get(ctx, types.KeyFeeRegister, &oldFeeReeRegParams)

	var oldFeeRecordParams uint64
	paramsSubspace.Get(ctx, types.KeyFeeRecord, &oldFeeRecordParams)

	var oldDenomParams string
	paramsSubspace.Get(ctx, types.KeyDenom, &oldDenomParams)

	params := types.NewParams(
		oldFeeReeRegParams,
		oldFeeRecordParams,
		types.PurchaseStorageFee,
		oldDenomParams,
		types.DefaultStorageLimit,
		types.DefaultMaxStorageLimit,
	)

	// save new paramset
	paramsSubspace.SetParamSet(ctx, &params)
}

func pruneWrkChainHeights(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec, wrkchainId uint64) (uint64, uint64) {
	store := ctx.KVStore(storeKey)
	// reverse order (desc) to get newest first. Only want to prune old heights
	// when default in-state storage limit is reached
	iterator := sdk.KVStoreReversePrefixIterator(store, types.WrkChainAllBlocksKey(wrkchainId))
	numInState := uint64(0)
	lowestHeightInState := uint64(0)
	limitCounter := uint64(0)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldWrkChainBlockHash v040.WrkChainBlock
		cdc.MustUnmarshal(iterator.Value(), &oldWrkChainBlockHash)
		limitCounter = limitCounter + 1

		if limitCounter > types.DefaultStorageLimit {
			store.Delete(types.WrkChainBlockKey(wrkchainId, oldWrkChainBlockHash.Height))
		} else {
			lowestHeightInState = oldWrkChainBlockHash.Height
			numInState = numInState + 1
		}
	}
	return numInState, lowestHeightInState
}

func pruneWrkChains(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) {
	// loop through WrkChains, then each WrkChain's hashes
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RegisteredWrkChainPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var oldWrkchain v040.WrkChain
		cdc.MustUnmarshal(iterator.Value(), &oldWrkchain)

		numInState, lowestHeightInState := pruneWrkChainHeights(ctx, storeKey, cdc, oldWrkchain.WrkchainId)

		newWrkchain := types.WrkChain{
			WrkchainId:       oldWrkchain.WrkchainId,
			Moniker:          oldWrkchain.Moniker,
			Name:             oldWrkchain.Name,
			Genesis:          oldWrkchain.Genesis,
			Type:             oldWrkchain.Type,
			Lastblock:        oldWrkchain.Lastblock,
			NumBlocksInState: numInState,
			LowestHeight:     lowestHeightInState,
			InStateLimit:     types.DefaultStorageLimit,
			RegTime:          oldWrkchain.RegTime,
			Owner:            oldWrkchain.Owner,
		}

		// save the updated WrkChain
		store.Set(types.WrkChainKey(oldWrkchain.WrkchainId), cdc.MustMarshal(&newWrkchain))
	}
}

// MigrateStore performs in-place store migrations from SDK v0.42 of the Wrkchain module to SDK v0.45.
// The migration includes:
//
// - Adding new params
// - Setting default in-store limit for Wrkchains
// - Pruning all old hashes exceeding in-state limit
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, paramsSubspace paramstypes.Subspace, cdc codec.BinaryCodec) error {

	// migrate Params
	migrateParams(ctx, paramsSubspace)

	// prune hashes
	pruneWrkChains(ctx, storeKey, cdc)

	return nil
}
