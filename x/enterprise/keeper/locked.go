package keeper

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// __ACCOUNT_QUERIES_____________________________________________________

func (k Keeper) GetEnterpriseUserAccount(ctx sdk.Context, owner sdk.AccAddress) types.EnterpriseUserAccount {
	locked := k.GetLockedUndForAccount(ctx, owner)
	bankSupply := k.bankKeeper.GetBalance(ctx, owner, k.GetParamDenom(ctx))
	lockedEFUND := locked.Amount
	spent := k.GetSpentEFUNDForAccount(ctx, owner)

	spendable := bankSupply.Add(lockedEFUND)

	userAccount := types.EnterpriseUserAccount{
		Owner:         owner.String(),
		LockedEfund:   lockedEFUND,
		GeneralSupply: bankSupply,
		SpentEfund:    spent.Amount,
		Spendable:     spendable,
	}

	return userAccount
}

// __TOTAL_LOCKED_FUND___________________________________________________

// GetTotalLockedUnd returns the total locked eFUND
func (k Keeper) GetTotalLockedUnd(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.TotalLockedUndKey)

	if bz == nil {
		return sdk.NewInt64Coin(k.GetParamDenom(ctx), 0)
	}

	var totalLocked sdk.Coin
	k.cdc.MustUnmarshal(bz, &totalLocked)
	return totalLocked
}

// SetTotalLockedUnd sets the total locked FUND
func (k Keeper) SetTotalLockedUnd(ctx sdk.Context, totalLocked sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TotalLockedUndKey, k.cdc.MustMarshal(&totalLocked))
	return nil
}

//__SPENT_eFUND__________________________________________________________

// GetTotalSpentEFUND returns the total spent eFUND
func (k Keeper) GetTotalSpentEFUND(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.TotalSpentEFUNDKey)

	if bz == nil {
		return sdk.NewInt64Coin(k.GetParamDenom(ctx), 0)
	}

	var totalUsed sdk.Coin
	k.cdc.MustUnmarshal(bz, &totalUsed)
	return totalUsed
}

// SetTotalSpentEFUND sets the total used eFUND
func (k Keeper) SetTotalSpentEFUND(ctx sdk.Context, totalUsed sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TotalSpentEFUNDKey, k.cdc.MustMarshal(&totalUsed))
	return nil
}

// AccountHasSpentEFUND checks if given account has an entry for spent eFUND
func (k Keeper) AccountHasSpentEFUND(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.SpentEFUNDAddressStoreKey(address)
	return store.Has(addressKeyBz)
}

// GetSpentEFUNDForAccount Gets a record for used eFUND for a given address
func (k Keeper) GetSpentEFUNDForAccount(ctx sdk.Context, address sdk.AccAddress) types.SpentEFUND {
	store := ctx.KVStore(k.storeKey)

	if !k.AccountHasSpentEFUND(ctx, address) {
		// return a new zero coin
		return types.SpentEFUND{
			Owner:  address.String(),
			Amount: sdk.NewInt64Coin(k.GetParamDenom(ctx), 0),
		}
	}

	bz := store.Get(types.SpentEFUNDAddressStoreKey(address))
	var spentEFUND types.SpentEFUND
	k.cdc.MustUnmarshal(bz, &spentEFUND)
	return spentEFUND
}

// SetSpentEFUNDForAccount Sets the spent eFUND data
func (k Keeper) SetSpentEFUNDForAccount(ctx sdk.Context, spent types.SpentEFUND) error {
	store := ctx.KVStore(k.storeKey)
	owner, accErr := sdk.AccAddressFromBech32(spent.Owner)
	if accErr != nil {
		return accErr
	}
	store.Set(types.SpentEFUNDAddressStoreKey(owner), k.cdc.MustMarshal(&spent))
	return nil
}

func (k Keeper) GetSpentEFUNDAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin {
	return k.GetSpentEFUNDForAccount(ctx, address).Amount
}

// Get an iterator over all accounts with spent eFUND
func (k Keeper) GetAllSpentEFUNDAccountsIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return storetypes.KVStorePrefixIterator(store, types.SpentEFUNDAddressKeyPrefix)
}

func (k Keeper) GetAllSpentEFUNDs(ctx sdk.Context) (spentEFUNDs []types.SpentEFUND) {
	spentIterator := k.GetAllSpentEFUNDAccountsIterator(ctx)

	for ; spentIterator.Valid(); spentIterator.Next() {
		var spent types.SpentEFUND
		k.cdc.MustUnmarshal(spentIterator.Value(), &spent)
		spentEFUNDs = append(spentEFUNDs, spent)
	}

	return spentEFUNDs
}

// __MINTER_AND_UNLOCKER________________________________________________

// CreateAndLockEFUND creates and locks eFUND
// CreateAndLockEFUND to be used in BeginBlocker.
// Minting will be handled in UnlockCoinsForFees
func (k Keeper) CreateAndLockEFUND(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coin) error {
	if amount.Amount.IsZero() {
		// skip as no coins need to be minted
		return nil
	}
	// keep track of how much FUND is locked for this account, and in total
	err := k.incrementLockedUnd(ctx, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}

// UnlockAndMintCoinsForFees unlocks any locked eFUND and mints them in the bank keeper
// to pay for fees
func (k Keeper) UnlockAndMintCoinsForFees(ctx sdk.Context, feePayer sdk.AccAddress, feesToPay sdk.Coins) error {

	logger := k.Logger(ctx)
	lockedUnd := k.GetLockedUndForAccount(ctx, feePayer).Amount
	lockedUndCoins := sdk.NewCoins(lockedUnd)
	feeNund := feesToPay.AmountOf(k.GetParamDenom(ctx))
	feeNundCoin := sdk.NewCoin(k.GetParamDenom(ctx), feeNund)
	_, feeToPay := feesToPay.Find(k.GetParamDenom(ctx))
	//blockTime := uint64(ctx.BlockHeader().Time.Unix())

	// calculate how much Locked FUND would be left over after deducting Tx fees
	_, hasNeg := lockedUndCoins.SafeSub(feeToPay)

	if !hasNeg {
		// locked FUND >= total fees
		// mint the fee amount to allow for payment
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, feesToPay)
		if err != nil {
			return err
		}

		// Send them to the purchaser's account
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, feesToPay)
		if err != nil {
			return err
		}

		// decrement the tracked locked eFUND
		err = k.decrementLockedUnd(ctx, feePayer, feeNundCoin)
		if err != nil {
			return err
		}

		// increment the used eFUND
		err = k.incrementSpentEFUND(ctx, feePayer, feeNundCoin)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeUndUnlocked,
				sdk.NewAttribute(types.AttributeKeyPurchaser, feePayer.String()),
				sdk.NewAttribute(types.AttributeKeyAmount, feeNundCoin.String()),
			),
		)

		logger.Debug("enterprise unlocking und", "for", feePayer.String(), "amt", feeNundCoin.Amount)

	} else {
		// calculate how much can be unlocked and minted, and if, by minting, the account
		// would have enough to pay for the fees. If not, don't mint/unlock

		// How many spendable FUND does the account have
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, feePayer)

		// calculate how much would be available if FUND were unlocked
		potentiallyAvailable := spendableCoins.Add(lockedUndCoins...)

		// is this enough to pay for the fees?
		_, hasNeg2 := potentiallyAvailable.SafeSub(feeToPay)

		// only undelegate & unlock if the resulting unlock will be enough to pay for the fees.
		if !hasNeg2 {
			// undelegate the fee amount to allow for payment
			//err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, lockedUndCoins)

			//first mint
			err := k.bankKeeper.MintCoins(ctx, types.ModuleName, lockedUndCoins)
			if err != nil {
				return err
			}

			// Send them to the purchaser's account
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, lockedUndCoins)
			if err != nil {
				return err
			}

			// decrement the tracked locked eFUND
			err = k.decrementLockedUnd(ctx, feePayer, lockedUnd)
			if err != nil {
				return err
			}

			// increment the used eFUND
			err = k.incrementSpentEFUND(ctx, feePayer, lockedUnd)
			if err != nil {
				return err
			}

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUndUnlocked,
					sdk.NewAttribute(types.AttributeKeyPurchaser, feePayer.String()),
					sdk.NewAttribute(types.AttributeKeyAmount, lockedUnd.String()),
				),
			)

			logger.Debug("enterprise unlocking und", "for", feePayer.String(), "amt", lockedUnd.Amount)

		}
	}

	return nil
}

//__LOCKED_FUND__________________________________________________________

// Check if a record exists for locked FUND given an account address
func (k Keeper) AccountHasLockedUnd(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.LockedUndAddressStoreKey(address)
	return store.Has(addressKeyBz)
}

func (k Keeper) IsLocked(ctx sdk.Context, address sdk.AccAddress) bool {
	return k.GetLockedUndForAccount(ctx, address).Amount.IsPositive()
}

// Gets a record for Locked FUND for a given address
func (k Keeper) GetLockedUndForAccount(ctx sdk.Context, address sdk.AccAddress) types.LockedUnd {
	store := ctx.KVStore(k.storeKey)

	if !k.AccountHasLockedUnd(ctx, address) {
		// return a new empty EnterpriseUndPurchaseOrder struct
		return types.LockedUnd{
			Owner:  address.String(),
			Amount: sdk.NewInt64Coin(k.GetParamDenom(ctx), 0),
		}
	}

	bz := store.Get(types.LockedUndAddressStoreKey(address))
	var lockedUnd types.LockedUnd
	k.cdc.MustUnmarshal(bz, &lockedUnd)
	return lockedUnd
}

func (k Keeper) GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin {
	return k.GetLockedUndForAccount(ctx, address).Amount
}

// Get an iterator over all accounts with Locked FUND
func (k Keeper) GetAllLockedUndAccountsIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(k.storeKey)
	return storetypes.KVStorePrefixIterator(store, types.LockedUndAddressKeyPrefix)
}

func (k Keeper) GetAllLockedUnds(ctx sdk.Context) (lockedUnds []types.LockedUnd) {
	lockedIterator := k.GetAllLockedUndAccountsIterator(ctx)

	for ; lockedIterator.Valid(); lockedIterator.Next() {
		var l types.LockedUnd
		k.cdc.MustUnmarshal(lockedIterator.Value(), &l)
		lockedUnds = append(lockedUnds, l)
	}

	return lockedUnds
}

// Sets the Locked FUND data
func (k Keeper) SetLockedUndForAccount(ctx sdk.Context, lockedUnd types.LockedUnd) error {
	// must have an owner
	//if lockedUnd.Owner.Empty() {
	//	return errorsmod.Wrap(types.ErrMissingData, "unable to set locked und - owner cannot be empty")
	//}
	owner, accErr := sdk.AccAddressFromBech32(lockedUnd.Owner)
	if accErr != nil {
		return accErr
	}

	// must be a positive amount, or zero
	if lockedUnd.Amount.IsNegative() {
		return errorsmod.Wrap(types.ErrMissingData, "unable to set locked und - amount must be positive")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.LockedUndAddressStoreKey(owner), k.cdc.MustMarshal(&lockedUnd))

	return nil
}

func (k Keeper) incrementSpentEFUND(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin) error {
	spentEFUND := k.GetSpentEFUNDForAccount(ctx, address)
	spentEFUND.Amount = spentEFUND.Amount.Add(amount)
	err := k.SetSpentEFUNDForAccount(ctx, spentEFUND)

	if err != nil {
		return err
	}

	totalSpent := k.GetTotalSpentEFUND(ctx)
	newTotalUsed := totalSpent.Add(amount)
	err = k.SetTotalSpentEFUND(ctx, newTotalUsed)

	if err != nil {
		return err
	}
	return nil
}

// incrementLockedUnd increments the amount of locked FUND - used when purchase order is accepted
func (k Keeper) incrementLockedUnd(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin) error {

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

// decrementLockedUnd decrements the amount of locked FUND - used when purchase order is accepted
func (k Keeper) decrementLockedUnd(ctx sdk.Context, address sdk.AccAddress, amount sdk.Coin) error {

	lockedUnd := k.GetLockedUndForAccount(ctx, address)
	lockedCoins := sdk.NewCoins(lockedUnd.Amount)
	//subAmountCoins := sdk.NewCoins(amount)

	_, hasNeg := lockedCoins.SafeSub(amount)

	if hasNeg {
		lockedUnd.Amount = sdk.NewInt64Coin(k.GetParamDenom(ctx), 0)
	} else {
		lockedUnd.Amount = lockedUnd.Amount.Sub(amount)
	}

	err := k.SetLockedUndForAccount(ctx, lockedUnd)
	if err != nil {
		return err
	}

	// update total
	totalLocked := k.GetTotalLockedUnd(ctx)
	totalLockedCoins := sdk.NewCoins(totalLocked)
	_, hasNeg = totalLockedCoins.SafeSub(amount)

	if hasNeg {
		err = k.SetTotalLockedUnd(ctx, sdk.NewInt64Coin(k.GetParamDenom(ctx), 0))
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
