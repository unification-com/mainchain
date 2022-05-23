package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// __ACCOUNT_QUERIES_____________________________________________________

func (k Keeper) GetEnterpriseUserAccount(ctx sdk.Context, owner sdk.AccAddress) types.EnterpriseUserAccount {
	locked := k.GetLockedUndForAccount(ctx, owner)
	unlocked := k.bankKeeper.GetBalance(ctx, owner, k.GetParamDenom(ctx))
	lockedCoin := locked.Amount

	total := unlocked.Add(lockedCoin)

	userAccount := types.EnterpriseUserAccount{
		Owner:    owner.String(),
		Locked:   lockedCoin,
		Unlocked: unlocked,
		Total:    total,
	}

	return userAccount
}

// __TOTAL_LOCKED_FUND___________________________________________________

// GetTotalLockedUnd returns the total locked FUND
func (k Keeper) GetTotalLockedUnd(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.TotalLockedUndKey)

	if bz == nil {
		return sdk.NewInt64Coin(k.GetParamDenom(ctx), 0)
	}

	var totalLocked sdk.Coin
	k.cdc.MustUnmarshalBinaryBare(bz, &totalLocked)
	return totalLocked
}

// GetTotalUnLockedUnd returns the amount of unlocked FUND - i.e. in active
// circulation (totalSupply - locked)
func (k Keeper) GetTotalUnLockedUnd(ctx sdk.Context) sdk.Coin {
	supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(k.GetParamDenom(ctx))
	total := sdk.NewCoin(k.GetParamDenom(ctx), supply)
	locked := k.GetTotalLockedUnd(ctx)

	unlocked := total.Sub(locked)

	return unlocked
}

// GetTotalUndSupply returns the total FUND in supply, obtained from the supply module's keeper
func (k Keeper) GetTotalUndSupply(ctx sdk.Context) sdk.Coin {
	supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(k.GetParamDenom(ctx))
	total := sdk.NewCoin(k.GetParamDenom(ctx), supply)
	return total
}

// GetEnterpriseSupplyIncludingLockedUnd returns information including total FUND supply, total locked and unlocked
func (k Keeper) GetEnterpriseSupplyIncludingLockedUnd(ctx sdk.Context) types.UndSupply {
	supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(k.GetParamDenom(ctx))
	total := sdk.NewCoin(k.GetParamDenom(ctx), supply)
	locked := k.GetTotalLockedUnd(ctx)

	unlocked := total.Sub(locked)

	denom := k.GetParamDenom(ctx)

	totalSupply := types.UndSupply{
		Denom:  denom,
		Locked: locked.Amount.Uint64(),
		Amount: unlocked.Amount.Uint64(),
		Total:  total.Amount.Uint64(),
	}

	return totalSupply
}

func (k Keeper) GetTotalSupplyWithLockedNundRemoved(ctx sdk.Context) sdk.Coins {
	supplyCoins := k.bankKeeper.GetSupply(ctx).GetTotal()
	locked := k.GetTotalLockedUnd(ctx)

	for i, c := range supplyCoins {
		if c.Denom == k.GetParamDenom(ctx) {
			unlocked := c.Sub(locked)
			supplyCoins[i] = unlocked
		}
	}

	return supplyCoins
}

func (k Keeper) GetSupplyOfWithLockedNundRemoved(ctx sdk.Context, denom string) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)

	if denom == k.GetParamDenom(ctx) {
		locked := k.GetTotalLockedUnd(ctx)
		unlocked := supply.Sub(locked.Amount)
		return unlocked
	} else {
		return supply
	}
}

// SetTotalLockedUnd sets the total locked FUND
func (k Keeper) SetTotalLockedUnd(ctx sdk.Context, totalLocked sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TotalLockedUndKey, k.cdc.MustMarshalBinaryBare(&totalLocked))
	return nil
}

// __MINTER_AND_UNLOCKER________________________________________________

// MintCoinsAndLock implements an alias call to the underlying bank keeper's MintCoins
// MintCoinsAndLock to be used in BeginBlocker.
func (k Keeper) MintCoinsAndLock(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coin) error {
	if amount.Amount.IsZero() {
		// skip as no coins need to be minted
		return nil
	}

	newCoins := sdk.NewCoins(amount)
	//first mint
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
	if err != nil {
		return err
	}

	// Send them to the purchaser's account
	err = k.sendCoinsFromModuleToAccount(ctx, recipient, newCoins)
	if err != nil {
		return err
	}

	// Delegate the Enterprise FUND module so they can't be spent
	err = k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, recipient, types.ModuleName, newCoins)
	if err != nil {
		return err
	}

	// keep track of how much FUND is locked for this account, and in total
	err = k.incrementLockedUnd(ctx, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UnlockCoinsForFees(ctx sdk.Context, feePayer sdk.AccAddress, feesToPay sdk.Coins) error {

	logger := k.Logger(ctx)
	lockedUnd := k.GetLockedUndForAccount(ctx, feePayer).Amount
	lockedUndCoins := sdk.NewCoins(lockedUnd)
	//blockTime := uint64(ctx.BlockHeader().Time.Unix())

	// calculate how much Locked FUND would be left over after deducting Tx fees
	_, hasNeg := lockedUndCoins.SafeSub(feesToPay)

	if !hasNeg {
		// locked FUND >= total fees
		// undelegate the fee amount to allow for payment
		err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, feesToPay)

		if err != nil {
			return err
		}

		// decrement the tracked locked FUND
		feeNund := feesToPay.AmountOf(k.GetParamDenom(ctx))
		feeNundCoin := sdk.NewCoin(k.GetParamDenom(ctx), feeNund)
		err = k.decrementLockedUnd(ctx, feePayer, feeNundCoin)
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
		// calculate how much can be undelegated, and if, by undelegating, the account
		// would have enough to pay for the fees. If not, don't undelegate
		//feePayerAcc := k.accKeeper.GetAccount(ctx, feePayer)

		// How many spendable FUND does the account have
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, feePayer)

		// calculate how much would be available if FUND were unlocked
		potentiallyAvailable := spendableCoins.Add(lockedUndCoins...)

		// is this enough to pay for the fees
		_, hasNeg := potentiallyAvailable.SafeSub(feesToPay)

		// only undelegate & unlock if the resulting unlock will be enough to pay for the fees.
		if !hasNeg {
			// undelegate the fee amount to allow for payment
			err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.ModuleName, feePayer, lockedUndCoins)

			if err != nil {
				return err
			}

			err = k.decrementLockedUnd(ctx, feePayer, lockedUnd)
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

// sendCoinsFromModuleToAccount implements an alias call to the underlying supply keeper's SendCoinsFromModuleToAccount
// Used in MintCoinsAndLock to send newly minted coins from enterprise module to recipient account
func (k Keeper) sendCoinsFromModuleToAccount(ctx sdk.Context, recipientAddr sdk.AccAddress, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, newCoins)
}

//__LOCKED_FUND__________________________________________________________

// Check if a record exists for locked FUND given an account address
func (k Keeper) AccountHasLockedUnd(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.AddressStoreKey(address)
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

	bz := store.Get(types.AddressStoreKey(address))
	var lockedUnd types.LockedUnd
	k.cdc.MustUnmarshalBinaryBare(bz, &lockedUnd)
	return lockedUnd
}

func (k Keeper) GetLockedUndAmountForAccount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin {
	return k.GetLockedUndForAccount(ctx, address).Amount
}

// Get an iterator over all accounts with Locked FUND
func (k Keeper) GetAllLockedUndAccountsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LockedUndAddressKeyPrefix)
}

func (k Keeper) GetAllLockedUnds(ctx sdk.Context) (lockedUnds []types.LockedUnd) {
	lockedIterator := k.GetAllLockedUndAccountsIterator(ctx)

	for ; lockedIterator.Valid(); lockedIterator.Next() {
		var l types.LockedUnd
		k.cdc.MustUnmarshalBinaryBare(lockedIterator.Value(), &l)
		lockedUnds = append(lockedUnds, l)
	}

	return lockedUnds
}

// Deletes the accepted purchase order once processed
func (k Keeper) DeleteLockedUndForAccount(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddressStoreKey(address))
}

// Sets the Locked FUND data
func (k Keeper) SetLockedUndForAccount(ctx sdk.Context, lockedUnd types.LockedUnd) error {
	// must have an owner
	//if lockedUnd.Owner.Empty() {
	//	return sdkerrors.Wrap(types.ErrMissingData, "unable to set locked und - owner cannot be empty")
	//}
	owner, accErr := sdk.AccAddressFromBech32(lockedUnd.Owner)
	if accErr != nil {
		return accErr
	}

	// must be a positive amount, or zero
	if lockedUnd.Amount.IsNegative() {
		return sdkerrors.Wrap(types.ErrMissingData, "unable to set locked und - amount must be positive")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.AddressStoreKey(owner), k.cdc.MustMarshalBinaryBare(&lockedUnd))

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
	subAmountCoins := sdk.NewCoins(amount)

	_, hasNeg := lockedCoins.SafeSub(subAmountCoins)

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
	_, hasNeg = totalLockedCoins.SafeSub(subAmountCoins)

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
