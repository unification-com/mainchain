package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

//__PURCHASE_ORDER_ID___________________________________________________

// GetHighestPurchaseOrderID gets the highest purchase order ID
func (k Keeper) GetHighestPurchaseOrderID(ctx sdk.Context) (purchaseOrderID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestPurchaseOrderIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}
	// convert from bytes to uint64
	purchaseOrderID = types.GetPurchaseOrderIDFromBytes(bz)
	return purchaseOrderID, nil
}

// SetHighestPurchaseOrderID sets the new proposal ID to the store
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
		return types.NewEnterpriseUndPurchaseOrder()
	}

	bz := store.Get(types.PurchaseOrderKey(purchaseOrderID))
	var enterpriseUndPurchaseOrder types.EnterpriseUndPurchaseOrder
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &enterpriseUndPurchaseOrder)
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

// IteratePurchaseOrders iterates over the all the purchase orders and performs a callback function
func (k Keeper) IteratePurchaseOrders(ctx sdk.Context, cb func(purchaseOrder types.EnterpriseUndPurchaseOrder) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseOrderIDKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var po types.EnterpriseUndPurchaseOrder
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &po)

		if cb(po) {
			break
		}
	}
}

// GetAllPurchaseOrders returns all the purchase orders from store
func (k Keeper) GetAllPurchaseOrders(ctx sdk.Context) (purchaseOrders types.PurchaseOrders) {
	k.IteratePurchaseOrders(ctx, func(po types.EnterpriseUndPurchaseOrder) bool {
		purchaseOrders = append(purchaseOrders, po)
		return false
	})
	return
}

// GetPurchaseOrdersFiltered retrieves purchase orders filtered by a given set of params which
// include pagination parameters along a purchase order status.
//
// NOTE: If no filters are provided, all proposals will be returned in paginated
// form.
func (k Keeper) GetPurchaseOrdersFiltered(ctx sdk.Context, params types.QueryPurchaseOrdersParams) []types.EnterpriseUndPurchaseOrder {
	purchaseOrders := k.GetAllPurchaseOrders(ctx)
	filteredPurchaseOrders := make([]types.EnterpriseUndPurchaseOrder, 0, len(purchaseOrders))

	for _, po := range purchaseOrders {
		matchStatus, matchPurchaser := true, true

		// match status (if supplied/valid)
		if types.ValidPurchaseOrderStatus(params.PurchaseOrderStatus) {
			matchStatus = po.Status == params.PurchaseOrderStatus
		}

		if len(params.Purchaser) > 0 {
			matchPurchaser = po.Purchaser.String() == params.Purchaser.String()
		}

		if matchStatus && matchPurchaser {
			filteredPurchaseOrders = append(filteredPurchaseOrders, po)
		}
	}

	start, end := client.Paginate(len(filteredPurchaseOrders), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		filteredPurchaseOrders = []types.EnterpriseUndPurchaseOrder{}
	} else {
		filteredPurchaseOrders = filteredPurchaseOrders[start:end]
	}

	return filteredPurchaseOrders
}

// Sets the Purchase Order data
func (k Keeper) SetPurchaseOrder(ctx sdk.Context, purchaseOrder types.EnterpriseUndPurchaseOrder) error {
	// must have a purchaser
	if purchaseOrder.Purchaser.Empty() {
		return sdkerrors.Wrap(types.ErrMissingData, "unable to set purchase order - purchaser cannot be empty")
	}

	if !purchaseOrder.Amount.IsValid() {
		return sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount not valid")
	}

	// must be a positive amount
	if purchaseOrder.Amount.IsZero() || purchaseOrder.Amount.IsNegative() {
		return sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount must be positive")
	}

	//must have an ID
	if purchaseOrder.PurchaseOrderID == 0 {
		return sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - id must be positive non-zero")
	}

	if !types.ValidPurchaseOrderStatus(purchaseOrder.Status) {
		return sdkerrors.Wrap(types.ErrInvalidStatus, "unable to set purchase order - invalid status")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.PurchaseOrderKey(purchaseOrder.PurchaseOrderID), k.cdc.MustMarshalBinaryLengthPrefixed(purchaseOrder))

	return nil
}

func (k Keeper) RaiseNewPurchaseOrder(ctx sdk.Context, purchaser sdk.AccAddress, amount sdk.Coin) (uint64, error) {

	logger := k.Logger(ctx)

	if !k.AddressIsWhitelisted(ctx, purchaser) {
		return 0, sdkerrors.Wrap(types.ErrNotAuthorisedToRaisePO,  fmt.Sprintf("%s is not whitelisted", purchaser))
	}

	purchaseOrderID, err := k.GetHighestPurchaseOrderID(ctx)
	if err != nil {
		return 0, err
	}

	purchaseOrder := k.GetPurchaseOrder(ctx, purchaseOrderID)
	purchaseOrder.PurchaseOrderID = purchaseOrderID
	purchaseOrder.Purchaser = purchaser
	purchaseOrder.Amount = amount
	purchaseOrder.Status = types.StatusRaised
	purchaseOrder.RaisedTime = ctx.BlockHeader().Time.Unix()

	err = k.SetPurchaseOrder(ctx, purchaseOrder)
	if err != nil {
		return 0, err
	}

	k.SetHighestPurchaseOrderID(ctx, purchaseOrderID+1)

	logger.Info("enterprise und purchase order raised", "id", purchaseOrderID, "from", purchaser.String(), "amt", amount.String())

	return purchaseOrderID, nil
}

func (k Keeper) IsAuthorisedToDecide(ctx sdk.Context, signer sdk.AccAddress) bool {
	entSigners := k.GetParamEntSignersAsAddressArray(ctx)
	isAuthorised := false
	for _, authAddr := range entSigners {
		if signer.Equals(authAddr) {
			isAuthorised = true
		}
	}
	return isAuthorised
}

func (k Keeper) ProcessPurchaseOrderDecision(ctx sdk.Context, purchaseOrderID uint64, decision types.PurchaseOrderStatus, signer sdk.AccAddress) error {

	logger := k.Logger(ctx)

	if !k.PurchaseOrderExists(ctx, purchaseOrderID) {
		errMsg := fmt.Sprintf("id: %d", purchaseOrderID)
		return sdkerrors.Wrap(types.ErrPurchaseOrderDoesNotExist, errMsg)
	}

	if !types.ValidPurchaseOrderAcceptRejectStatus(decision) {
		return sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")
	}

	if !k.IsAuthorisedToDecide(ctx, signer) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")
	}

	purchaseOrder := k.GetPurchaseOrder(ctx, purchaseOrderID)

	if purchaseOrder.Status == types.StatusNil {
		errMsg := fmt.Sprintf("purchase order %d not raised!", purchaseOrderID)
		return sdkerrors.Wrap(types.ErrPurchaseOrderNotRaised, errMsg)
	}

	if purchaseOrder.Status != types.StatusRaised {
		errMsg := fmt.Sprintf("id %d already processed: %s", purchaseOrderID, purchaseOrder.Status.String())
		return sdkerrors.Wrap(types.ErrPurchaseOrderAlreadyProcessed, errMsg)
	}

	currentDecisions := purchaseOrder.Decisions
	for _, d := range currentDecisions {
		if d.Signer.Equals(signer) {
			errMsg := fmt.Sprintf("signer %s already decided: %s", signer.String(), d.Decision.String())
			return sdkerrors.Wrap(types.ErrSignerAlreadyMadeDecision, errMsg)
		}
	}

	poDecision := types.NewPurchaseOrderDecision(signer, decision)
	poDecision.DecisionTime = ctx.BlockHeader().Time.Unix()
	purchaseOrder.Decisions = append(purchaseOrder.Decisions, poDecision)

	// update the status
	err := k.SetPurchaseOrder(ctx, purchaseOrder)

	if err != nil {
		return err
	}

	logger.Info("enterprise und purchase order decision made", "id", purchaseOrderID, "signer", signer, "decision", decision)

	return nil
}
