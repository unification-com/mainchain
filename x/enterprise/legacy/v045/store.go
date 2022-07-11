package v045

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// setTotalUsed will calculate the total amount of eFUND used to pay towards fees to date, for all addresses.
// It will calculate the sum of eFUND purchases, based on purchase orders with the COMPLETE status,
// then subtract the total amount currently locked in order to get the amount spent on fees to date.
// It will then save this value to a new store entry for key []byte{0x98}
func setTotalUsed(ctx sdk.Context, storeKey sdk.StoreKey, paramsSubspace paramstypes.Subspace, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	var denom string
	paramsSubspace.Get(ctx, types.KeyDenom, &denom)

	// total completed POs. Completed POs are minted for the purchaser's use
	// so tally them for total purchased
	totalCompleted := sdk.NewInt64Coin(denom, int64(0))
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseOrderIDKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var po types.EnterpriseUndPurchaseOrder
		if err := cdc.Unmarshal(iterator.Value(), &po); err != nil {
			return sdkerrors.Wrap(err, "failed to unmarshal enterprise purchase order")
		}
		if po.Status == types.StatusCompleted {
			// PO was completed and purchased amount minted
			totalCompleted = totalCompleted.Add(po.Amount)
		}
	}

	// get how much is currently locked. Each time eFUND is spent, this value is
	// decremented by the amount of eFUND that was used to pay towards fees
	var totalLocked sdk.Coin
	totalLockedBz := store.Get(types.TotalLockedUndKey)

	if totalLockedBz == nil {
		totalLocked = sdk.NewInt64Coin(denom, 0)
	} else {
		if err := cdc.Unmarshal(totalLockedBz, &totalLocked); err != nil {
			return sdkerrors.Wrap(err, "failed to unmarshal total locked")
		}
	}

	// Total eFUND used to pay towards fees to date is (total purchased - total currently locked)
	totalUsed := totalCompleted.Sub(totalLocked)

	// save it to the new store key
	store.Set(types.TotalUsedUndKey, cdc.MustMarshal(&totalUsed))

	return nil
}

// getTotalPurchasedByAddress will get the total eFUND purchased for the given address.
// The value is the sum of COMPLETED purchase orders for the address. COMPLETED purchase orders
// have been processed and the purchased mount of eFUND has been minted.
func getTotalPurchasedByAddress(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec, denom, address string) (sdk.Coin, error) {
	store := ctx.KVStore(storeKey)

	// total completed POs. Completed POs are minted for the purchaser's use
	// so tally them for total purchased
	totalCompleted := sdk.NewInt64Coin(denom, int64(0))

	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseOrderIDKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var po types.EnterpriseUndPurchaseOrder
		if err := cdc.Unmarshal(iterator.Value(), &po); err != nil {
			return totalCompleted, sdkerrors.Wrap(err, "failed to unmarshal enterprise purchase order")
		}
		if po.Status == types.StatusCompleted && po.Purchaser == address {
			// PO was completed and purchased amount minted for this address
			totalCompleted = totalCompleted.Add(po.Amount)
		}
	}

	return totalCompleted, nil
}

// setTotalUsedByAddress will calculate the amount of eFUND used to pay towards fees to date.
// For each address in the LockedUnd keystore, it will calculate the sum of their eFUND purchases,
// then subtract the amount currently locked in order to get the amount spent on fees to date.
// It will then save this value to a new store entry for key ([]byte{0x06}, acc.Bytes()...)
func setTotalUsedByAddress(ctx sdk.Context, storeKey sdk.StoreKey, paramsSubspace paramstypes.Subspace, cdc codec.BinaryCodec) error {

	store := ctx.KVStore(storeKey)
	var denom string
	paramsSubspace.Get(ctx, types.KeyDenom, &denom)

	// LockedUndAddress is just a list of addresses and their current locked eFUND
	iterator := sdk.KVStorePrefixIterator(store, types.LockedUndAddressKeyPrefix)

	defer iterator.Close()
	// loop through each LockedUnd entry in the store
	for ; iterator.Valid(); iterator.Next() {
		var lockedUnd types.LockedUnd
		if err := cdc.Unmarshal(iterator.Value(), &lockedUnd); err != nil {
			return sdkerrors.Wrap(err, "failed to unmarshal locked und")
		}
		accAddr, accErr := sdk.AccAddressFromBech32(lockedUnd.Owner)
		if accErr != nil {
			return accErr
		}

		// get the total eFUND purchased for this address
		totalCompleted, err := getTotalPurchasedByAddress(ctx, storeKey, cdc, denom, lockedUnd.Owner)
		if err != nil {
			return err
		}

		// Total eFUND used to pay towards fees to date is (total purchased - total currently locked)
		totalUsed := totalCompleted.Sub(lockedUnd.Amount)

		// save it to the new store key
		store.Set(types.UsedUndAddressStoreKey(accAddr), cdc.MustMarshal(&totalUsed))
	}

	return nil
}

// MigrateStore performs in-place store migrations from SDK v0.42 of the Enterprise module to SDK v0.45.
// The migration includes:
//
// - Adding new total eFUND usage store
// - Adding new eFUND usage store for each account with eFUND
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, paramsSubspace paramstypes.Subspace, cdc codec.BinaryCodec) error {

	if err := setTotalUsed(ctx, storeKey, paramsSubspace, cdc); err != nil {
		return err
	}

	return setTotalUsedByAddress(ctx, storeKey, paramsSubspace, cdc)
}
