package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"

)

// __TOTAL_LOCKED_UND___________________________________________________

// GetTotalLockedUnd returns the total locked UND
func (k Keeper) GetTotalLockedUnd(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.TotalLockedUndKey)

	if bz == nil {
		return sdk.NewInt64Coin(types.DefaultDenomination, 0)
	}

	var totalLocked sdk.Coin
	k.cdc.MustUnmarshalBinaryBare(bz, &totalLocked)
	return totalLocked
}

// SetTotalLockedUnd sets the total locked UND
func (k Keeper) SetTotalLockedUnd(ctx sdk.Context, totalLocked sdk.Coin) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TotalLockedUndKey, k.cdc.MustMarshalBinaryBare(totalLocked))
	return nil
}

// __MINTER_AND_UNLOCKER________________________________________________

// MintCoinsAndLock implements an alias call to the underlying supply keeper's
// MintCoinsAndLock to be used in BeginBlocker.
func (k Keeper) MintCoinsAndLock(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coin) sdk.Error {
	if amount.Amount.IsZero() {
		// skip as no coins need to be minted
		return nil
	}

	newCoins := sdk.NewCoins(amount)
	//first mint
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)

	// Send them to the purchaser's account
	err = k.sendCoinsFromModuleToAccount(ctx, recipient, newCoins)
	if err != nil {
		return err
	}

	// Delegate the Enterprise UND module so they can't be spent
	err = k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, recipient, types.ModuleName, newCoins)
	if err != nil {
		return err
	}

	// keep track of how much UND is locked for this account, and in total
	err = k.incrementLockedUnd(ctx, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UnlockCoinsForFees(ctx sdk.Context, feePayer sdk.AccAddress, feesToPay sdk.Coins) sdk.Error {

	lockedUnd := k.GetLockedUndForAccount(ctx, feePayer).Amount
	lockedUndCoins := sdk.NewCoins(lockedUnd)
	blockTime := ctx.BlockHeader().Time

	// calculate how much Locked UND would be left over after deducting Tx fees
	_, hasNeg := lockedUndCoins.SafeSub(feesToPay)

	if !hasNeg {
		// locked UND >= total fees
		// undelegate the fee amount to allow for payment
		err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, feesToPay)

		if err != nil {
			return err
		}

		// decrement the tracked locked UND
		feeNund := feesToPay.AmountOf(types.DefaultDenomination)
		feeNundCoin := sdk.NewCoin(types.DefaultDenomination, feeNund)
		err = k.DecrementLockedUnd(ctx, feePayer, feeNundCoin)
		if err != nil {
			return err
		}
	} else {
		// calculate how much can be undelegated, and if, by undelegating, the account
		// would have enough to pay for the fees. If not, don't undelegate
		feePayerAcc := k.accKeeper.GetAccount(ctx, feePayer)

		// How many spendable UND does the account have
		spendableCoins := feePayerAcc.SpendableCoins(blockTime)

		// calculate how much would be available if UND were unlocked
		potentiallyAvailable := spendableCoins.Add(lockedUndCoins)

		// is this enough to pay for the fees
		_, hasNeg := potentiallyAvailable.SafeSub(feesToPay)

		if !hasNeg {
			// undelegate the fee amount to allow for payment
			err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, lockedUndCoins)

			if err != nil {
				return err
			}

			err = k.DecrementLockedUnd(ctx, feePayer, lockedUnd)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// sendCoinsFromModuleToAccount implements an alias call to the underlying supply keeper's SendCoinsFromModuleToAccount
// Used in MintCoinsAndLock to send newly minted coins from enterprise module to recipient account
func (k Keeper) sendCoinsFromModuleToAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, newCoins sdk.Coins) sdk.Error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, newCoins)
}

//__LOCKED_UND__________________________________________________________

// Check if a record exists for locked UND given an account address
func (k Keeper) AccountHasLockedUnd(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.AddressStoreKey(address)
	return store.Has(addressKeyBz)
}

func (k Keeper) IsLocked(ctx sdk.Context, address sdk.AccAddress) bool {
	return k.GetLockedUndForAccount(ctx, address).Amount.IsPositive()
}

// Gets a record for Locked UND for a given address
func (k Keeper) GetLockedUndForAccount(ctx sdk.Context, address sdk.AccAddress) types.LockedUnd {
	store := ctx.KVStore(k.storeKey)

	if !k.AccountHasLockedUnd(ctx, address) {
		// return a new empty EnterpriseUndPurchaseOrder struct
		return types.NewLockedUnd(address)
	}

	bz := store.Get(types.AddressStoreKey(address))
	var lockedUnd types.LockedUnd
	k.cdc.MustUnmarshalBinaryBare(bz, &lockedUnd)
	return lockedUnd
}

func (k Keeper) GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin {
	return k.GetLockedUndForAccount(ctx, address).Amount
}

// Get an iterator over all accounts with Locked UND
func (k Keeper) GetAllLockedUndAccountsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LockedUndAddressKeyPrefix)
}

// Deletes the accepted purchase order once processed
func (k Keeper) DeleteLockedUndForAccount(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddressStoreKey(address))
}

// Sets the Locked UND data
func (k Keeper) SetLockedUndForAccount(ctx sdk.Context, lockedUnd types.LockedUnd) sdk.Error {
	// must have an owner
	if lockedUnd.Owner.Empty() {
		return sdk.ErrInternal("unable to set locked und - owner cannot be empty")
	}

	// must be a positive amount, or zero
	if lockedUnd.Amount.IsNegative() {
		return sdk.ErrInternal("unable to set locked und - amount must be positive")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.AddressStoreKey(lockedUnd.Owner), k.cdc.MustMarshalBinaryBare(lockedUnd))

	return nil
}

// incrementLockedUnd increments the amount of locked UND - used when purchase order is accepted
func (k Keeper) incrementLockedUnd(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin) sdk.Error {

	lockedUnd := k.GetLockedUndForAccount(ctx, address)
	lockedUnd.Amount = lockedUnd.Amount.Add(amount)

	err := k.SetLockedUndForAccount(ctx, lockedUnd)
	if err != nil {
		return err
	}

	totalLocked := k.GetTotalLockedUnd(ctx)
	totalLocked = totalLocked.Add(amount)
	err = k.SetTotalLockedUnd(ctx, totalLocked)

	if err != nil {
		return err
	}

	return nil
}

// DecrementLockedUnd decrements the amount of locked UND - used when purchase order is accepted
func (k Keeper) DecrementLockedUnd(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin) sdk.Error {

	lockedUnd := k.GetLockedUndForAccount(ctx, address)
	lockedCoins := sdk.NewCoins(lockedUnd.Amount)
	subAmountCoins := sdk.NewCoins(amount)

	_, hasNeg := lockedCoins.SafeSub(subAmountCoins)

	if hasNeg {
		// delete
		k.DeleteLockedUndForAccount(ctx, address)
		return nil
	}

	lockedUnd.Amount = lockedUnd.Amount.Sub(amount)

	err := k.SetLockedUndForAccount(ctx, lockedUnd)
	if err != nil {
		return err
	}

	// update total
	totalLocked := k.GetTotalLockedUnd(ctx)
	totalLockedCoins := sdk.NewCoins(totalLocked)
	_, hasNeg = totalLockedCoins.SafeSub(subAmountCoins)

	if hasNeg {
		err = k.SetTotalLockedUnd(ctx, sdk.NewInt64Coin(types.DefaultDenomination, 0))
		if err != nil {
			return err
		}
		return nil
	}

	totalLocked = totalLocked.Sub(amount)
	err = k.SetTotalLockedUnd(ctx, totalLocked)

	if err != nil {
		return err
	}

	return nil
}