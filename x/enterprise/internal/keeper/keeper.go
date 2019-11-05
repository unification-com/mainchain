package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace   params.Subspace
	codespace    sdk.CodespaceType
	supplyKeeper types.SupplyKeeper
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the enterprise Keeper
func NewKeeper(storeKey sdk.StoreKey, supplyKeeper types.SupplyKeeper, paramSpace params.Subspace, codespace sdk.CodespaceType, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:     storeKey,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:    codespace,
		supplyKeeper: supplyKeeper,
		cdc:          cdc,
	}
}

func (k Keeper) GetCodeSpace() sdk.CodespaceType {
	return k.codespace
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
	if purchaseOrder.PurchaseOrderID <= 0 {
		return sdk.ErrInternal("unable to raise purchase order - id must be positive non-zero")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.PurchaseOrderKey(purchaseOrder.PurchaseOrderID), k.cdc.MustMarshalBinaryBare(purchaseOrder))

	return nil
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

	// Todo - actually process based on decision...

	return nil
}
