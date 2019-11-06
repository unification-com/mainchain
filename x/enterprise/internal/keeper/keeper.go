package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace   params.Subspace
	codespace    sdk.CodespaceType
	supplyKeeper types.SupplyKeeper
	accKeeper    auth.AccountKeeper
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the enterprise Keeper
func NewKeeper(storeKey sdk.StoreKey, supplyKeeper types.SupplyKeeper, accKeeper auth.AccountKeeper, paramSpace params.Subspace, codespace sdk.CodespaceType, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:     storeKey,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:    codespace,
		supplyKeeper: supplyKeeper,
		accKeeper:    accKeeper,
		cdc:          cdc,
	}
}

func (k Keeper) GetCodeSpace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

//__PARAMS______________________________________________________________

// GetParams returns the total set of Enterprise UND parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of Enterprise UND parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//__PURCHASE_ORDER_ID___________________________________________________

// GetHighestPurchaseOrderID gets the highest purchase order ID
func (k Keeper) GetHighestPurchaseOrderID(ctx sdk.Context) (purchaseOrderID uint64, err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestPurchaseOrderIDKey)
	if bz == nil {
		return 0, types.ErrInvalidGenesis(k.codespace, "initial proposal ID hasn't been set")
	}
	// convert from bytes to uint64
	purchaseOrderID = types.GetPurchaseOrderIDFromBytes(bz)
	return purchaseOrderID, nil
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetHighestPurchaseOrderID(ctx sdk.Context, purchaseOrderID uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert from uint64 to bytes for storage
	purchaseOrderIDbz := types.GetPurchaseOrderIDBytes(purchaseOrderID)
	store.Set(types.HighestPurchaseOrderIDKey, purchaseOrderIDbz)
}

// __TOTAL_LOCKED_UND___________________________________________________

// GetTotalLockedUnd retuns the total locked UND
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

//__PURCHASE_ORDERS______________________________________________

// Check if a raised purchase order for a given purchaseOrderID is in the store or not
func (k Keeper) PurchaseOrderExists(ctx sdk.Context, purchaseOrderID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	purchaseOrderIDbz := types.PurchaseOrderKey(purchaseOrderID)
	return store.Has(purchaseOrderIDbz)
}

// Gets a purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrder(ctx sdk.Context, purchaseOrderID uint64) types.EnterpriseUndPurchaseOrder {
	store := ctx.KVStore(k.storeKey)

	if !k.PurchaseOrderExists(ctx, purchaseOrderID) {
		// return a new empty EnterpriseUndPurchaseOrder struct
		return types.NewEnterpriseUnd()
	}

	bz := store.Get(types.PurchaseOrderKey(purchaseOrderID))
	var enterpriseUndPurchaseOrder types.EnterpriseUndPurchaseOrder
	k.cdc.MustUnmarshalBinaryBare(bz, &enterpriseUndPurchaseOrder)
	return enterpriseUndPurchaseOrder
}

// GetPurchaseOrderPurchaser - get the Purchaser address of a purchase order
// should be the same as the search term!
func (k Keeper) GetPurchaseOrderPurchaser(ctx sdk.Context, purchaseOrderID uint64) sdk.AccAddress {
	return k.GetPurchaseOrder(ctx, purchaseOrderID).Purchaser
}

// GetPurchaseOrderAmount - get the Amount of a raised purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrderAmount(ctx sdk.Context, purchaseOrderID uint64) sdk.Coin {
	return k.GetPurchaseOrder(ctx, purchaseOrderID).Amount
}

// GetPurchaseOrderStatus - get the Decision of a purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrderStatus(ctx sdk.Context, purchaseOrderID uint64) types.PurchaseOrderStatus {
	return k.GetPurchaseOrder(ctx, purchaseOrderID).Status
}

// Get an iterator over all Purchase Orders
func (k Keeper) GetAllPurchaseOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.PurchaseOrderIDKeyPrefix)
}

// Get an iterator over all accepted Purchase Orders - used by BeginBlocker to process minting
func (k Keeper) GetAllAcceptedPurchaseOrdersIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.AcceptedPurchaseOrderIDKeyPrefix)
}

// Sets the Purchase Order data
func (k Keeper) SetPurchaseOrder(ctx sdk.Context, purchaseOrder types.EnterpriseUndPurchaseOrder) sdk.Error {
	// must have a purchaser
	if purchaseOrder.Purchaser.Empty() {
		return sdk.ErrInternal("unable to raise purchase order - purchaser cannot be empty")
	}

	// must be a positive amount
	if purchaseOrder.Amount.IsZero() || purchaseOrder.Amount.IsNegative() {
		return sdk.ErrInternal("unable to raise purchase order - amount must be positive")
	}

	//must have an ID
	if purchaseOrder.PurchaseOrderID == 0 {
		return sdk.ErrInternal("unable to raise purchase order - id must be positive non-zero")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.PurchaseOrderKey(purchaseOrder.PurchaseOrderID), k.cdc.MustMarshalBinaryBare(purchaseOrder))

	return nil
}

// SetAcceptedPurchaseOrder creates a temporary store for an accepted purchase order. This will be queried by
// BeginBlocker to mint coins
// Todo - merge with standard PO and create new "complete" status
func (k Keeper) SetAcceptedPurchaseOrder(ctx sdk.Context, purchaseOrder types.EnterpriseUndPurchaseOrder) sdk.Error {

	//must have accepted status
	if purchaseOrder.Status != types.StatusAccepted {
		return sdk.ErrInternal("can only add accepted purchase orders")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.AcceptedPurchaseOrderKey(purchaseOrder.PurchaseOrderID), k.cdc.MustMarshalBinaryBare(purchaseOrder))

	return nil
}

// Deletes the accepted purchase order once processed
func (k Keeper) DeleteAcceptedPurchaseOrder(ctx sdk.Context, purchaseOrderID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AcceptedPurchaseOrderKey(purchaseOrderID))
}

func (k Keeper) RaiseNewPurchaseOrder(ctx sdk.Context, purchaser sdk.AccAddress, amount sdk.Coin) (uint64, sdk.Error) {

	purchaseOrderID, err := k.GetHighestPurchaseOrderID(ctx)
	if err != nil {
		return 0, err
	}

	purchaseOrder := k.GetPurchaseOrder(ctx, purchaseOrderID)
	purchaseOrder.PurchaseOrderID = purchaseOrderID
	purchaseOrder.Purchaser = purchaser
	purchaseOrder.Amount = amount
	purchaseOrder.Status = types.StatusRaised

	err = k.SetPurchaseOrder(ctx, purchaseOrder)
	if err != nil {
		return 0, err
	}

	k.SetHighestPurchaseOrderID(ctx, purchaseOrderID+1)

	return purchaseOrderID, nil
}

func (k Keeper) ProcessPurchaseOrder(ctx sdk.Context, purchaseOrderID uint64, decision types.PurchaseOrderStatus) sdk.Error {

	if !k.PurchaseOrderExists(ctx, purchaseOrderID) {
		errMsg := fmt.Sprintf("purchase order id does not exist: %d", purchaseOrderID)
		return types.ErrPurchaseOrderDoesNotExist(k.codespace, errMsg)
	}

	purchaseOrder := k.GetPurchaseOrder(ctx, purchaseOrderID)

	if purchaseOrder.Status == types.StatusAccepted || purchaseOrder.Status == types.StatusRejected {
		errMsg := fmt.Sprintf("purchase order %d already processed: %s", purchaseOrderID, purchaseOrder.Status.String())
		return types.ErrPurchaseOrderAlreadyProcessed(k.codespace, errMsg)
	}

	purchaseOrder.Status = decision

	// update the status
	err := k.SetPurchaseOrder(ctx, purchaseOrder)

	if err != nil {
		return err
	}

	if decision == types.StatusAccepted {
		err := k.SetAcceptedPurchaseOrder(ctx, purchaseOrder)

		if err != nil {
			return err
		}
	}

	return nil
}

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

	// Delegate the Enterprise UND so they can't be spent
	err = k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, recipient, types.ModuleName, newCoins)
	if err != nil {
		return err
	}

	// keep track of how much UND is locked for this account
	err = k.incrementLockedUnd(ctx, recipient, amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UnlockCoinsForFees(ctx sdk.Context, feePayer sdk.AccAddress, feesToPay sdk.Coins) sdk.Error {

	lockedUnd := k.GetLockedUnd(ctx, feePayer).Amount
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
func (k Keeper) LockedUndExists(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	addressKeyBz := types.AddressStoreKey(address)
	return store.Has(addressKeyBz)
}

func (k Keeper) IsLocked(ctx sdk.Context, address sdk.AccAddress) bool {
	return k.GetLockedUnd(ctx, address).Amount.IsPositive()
}

// Gets a record for Locked UND for a given address
func (k Keeper) GetLockedUnd(ctx sdk.Context, address sdk.AccAddress) types.LockedUnd {
	store := ctx.KVStore(k.storeKey)

	if !k.LockedUndExists(ctx, address) {
		// return a new empty EnterpriseUndPurchaseOrder struct
		return types.NewLockedUnd(address)
	}

	bz := store.Get(types.AddressStoreKey(address))
	var lockedUnd types.LockedUnd
	k.cdc.MustUnmarshalBinaryBare(bz, &lockedUnd)
	return lockedUnd
}

func (k Keeper) GetLockedUndAmount(ctx sdk.Context, address sdk.AccAddress) sdk.Coin {
	return k.GetLockedUnd(ctx, address).Amount
}

// Get an iterator over all accounts with Locked UND
func (k Keeper) GetAllLockedUndIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.LockedUndAddressKeyPrefix)
}

// Deletes the accepted purchase order once processed
func (k Keeper) DeleteLockedUnd(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddressStoreKey(address))
}

// Sets the Locked UND data
func (k Keeper) SetLockedUnd(ctx sdk.Context, lockedUnd types.LockedUnd) sdk.Error {
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

	lockedUnd := k.GetLockedUnd(ctx, address)
	lockedUnd.Amount = lockedUnd.Amount.Add(amount)

	err := k.SetLockedUnd(ctx, lockedUnd)
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

	lockedUnd := k.GetLockedUnd(ctx, address)
	lockedCoins := sdk.NewCoins(lockedUnd.Amount)
	subAmountCoins := sdk.NewCoins(amount)

	_, hasNeg := lockedCoins.SafeSub(subAmountCoins)

	if hasNeg {
		// delete
		k.DeleteLockedUnd(ctx, address)
		return nil
	}

	lockedUnd.Amount = lockedUnd.Amount.Sub(amount)

	err := k.SetLockedUnd(ctx, lockedUnd)
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

// ToDo - temporarily track unlocked/undelegated UND with separate keyPrefix. If WRKChain registration fails
// (e.g. WRKChain exists etc. - in WRKChain's handler.go) and source was Enterprise UND, then re-lock, and
// re-delegate the correct amount
