package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

//__PURCHASE_ORDER_ID___________________________________________________

// GetHighestPurchaseOrderID gets the highest purchase order ID
func (k Keeper) GetHighestPurchaseOrderID(ctx sdk.Context) (purchaseOrderID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.HighestPurchaseOrderIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial purchase order ID hasn't been set")
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

//__RAISED_PO_QUEUE
// Raised POs are added to the queue for processing in the ABCI Blocker
// Once processed (either accepted or rejected), they are removed from the
// queue

func (k Keeper) AddPoToRaisedQueue(ctx sdk.Context, purchaseOrderId uint64) {
	store := ctx.KVStore(k.storeKey)
	purchaseOrderIDbz := types.GetPurchaseOrderIDBytes(purchaseOrderId)
	queueKey := types.RaisedQueueStoreKey(purchaseOrderId)
	store.Set(queueKey, purchaseOrderIDbz)
}

func (k Keeper) PurchaseOrderIsInRaisedQueue(ctx sdk.Context, purchaseOrderId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	queueKey := types.RaisedQueueStoreKey(purchaseOrderId)
	return store.Has(queueKey)
}

func (k Keeper) RemovePurchaseOrderFromRaisedQueue(ctx sdk.Context, purchaseOrderId uint64) {
	if k.PurchaseOrderIsInRaisedQueue(ctx, purchaseOrderId) {
		store := ctx.KVStore(k.storeKey)
		queueKey := types.RaisedQueueStoreKey(purchaseOrderId)
		store.Delete(queueKey)
	}
}

// IterateRaisedQueue iterates over the all the raised purchase orders and performs a callback function
func (k Keeper) IterateRaisedQueue(ctx sdk.Context, cb func(purchaseOrderId uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.RaisedPoPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poId := types.GetPurchaseOrderIDFromBytes(iterator.Value())

		if cb(poId) {
			break
		}
	}
}

// GetAllRaisedPurchaseOrders returns all the purchase orders from store
func (k Keeper) GetAllRaisedPurchaseOrders(ctx sdk.Context) (purchaseOrderIds []uint64) {
	k.IterateRaisedQueue(ctx, func(poId uint64) bool {
		purchaseOrderIds = append(purchaseOrderIds, poId)
		return false
	})
	return
}

//__ACCEPTED_PO_QUEUE
// Accepted POs are added to the queue for processing in the ABCI Blocker
// Once processed and eFUNC minted, they are removed from the queue

func (k Keeper) AddPoToAcceptedQueue(ctx sdk.Context, purchaseOrderId uint64) {
	store := ctx.KVStore(k.storeKey)
	purchaseOrderIDbz := types.GetPurchaseOrderIDBytes(purchaseOrderId)
	queueKey := types.AcceptedQueueStoreKey(purchaseOrderId)
	store.Set(queueKey, purchaseOrderIDbz)
}

func (k Keeper) PurchaseOrderIsInAcceptedQueue(ctx sdk.Context, purchaseOrderId uint64) bool {
	store := ctx.KVStore(k.storeKey)
	queueKey := types.AcceptedQueueStoreKey(purchaseOrderId)
	return store.Has(queueKey)
}

func (k Keeper) RemovePurchaseOrderFromAcceptedQueue(ctx sdk.Context, purchaseOrderId uint64) {
	if k.PurchaseOrderIsInAcceptedQueue(ctx, purchaseOrderId) {
		store := ctx.KVStore(k.storeKey)
		queueKey := types.AcceptedQueueStoreKey(purchaseOrderId)
		store.Delete(queueKey)
	}
}

// IterateAcceptedQueue iterates over the all the Accepted purchase orders and performs a callback function
func (k Keeper) IterateAcceptedQueue(ctx sdk.Context, cb func(purchaseOrderId uint64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AcceptedPoPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poId := types.GetPurchaseOrderIDFromBytes(iterator.Value())

		if cb(poId) {
			break
		}
	}
}

// GetAllAcceptedPurchaseOrders returns all the purchase orders from store
func (k Keeper) GetAllAcceptedPurchaseOrders(ctx sdk.Context) (purchaseOrderIds []uint64) {
	k.IterateAcceptedQueue(ctx, func(poId uint64) bool {
		purchaseOrderIds = append(purchaseOrderIds, poId)
		return false
	})
	return
}

//__PURCHASE_ORDERS______________________________________________

// Check if a raised purchase order for a given purchaseOrderID is in the store or not
func (k Keeper) PurchaseOrderExists(ctx sdk.Context, purchaseOrderID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	purchaseOrderIDbz := types.PurchaseOrderKey(purchaseOrderID)
	return store.Has(purchaseOrderIDbz)
}

// Gets a purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrder(ctx sdk.Context, purchaseOrderID uint64) (types.EnterpriseUndPurchaseOrder, bool) {
	store := ctx.KVStore(k.storeKey)

	if !k.PurchaseOrderExists(ctx, purchaseOrderID) {
		// return a new empty EnterpriseUndPurchaseOrder struct
		return types.EnterpriseUndPurchaseOrder{}, false
	}

	bz := store.Get(types.PurchaseOrderKey(purchaseOrderID))
	var enterpriseUndPurchaseOrder types.EnterpriseUndPurchaseOrder
	k.cdc.MustUnmarshalBinaryBare(bz, &enterpriseUndPurchaseOrder)
	return enterpriseUndPurchaseOrder, true
}

// GetPurchaseOrderPurchaser - get the Purchaser address of a purchase order
// should be the same as the search term!
func (k Keeper) GetPurchaseOrderPurchaser(ctx sdk.Context, purchaseOrderID uint64) sdk.AccAddress {
	po, found := k.GetPurchaseOrder(ctx, purchaseOrderID)
	if !found {
		return sdk.AccAddress{}
	}
	accAddr, accErr := sdk.AccAddressFromBech32(po.Purchaser)

	if accErr != nil {
		return sdk.AccAddress{}
	}
	return accAddr
}

// GetPurchaseOrderAmount - get the Amount of a raised purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrderAmount(ctx sdk.Context, purchaseOrderID uint64) sdk.Coin {
	po, found := k.GetPurchaseOrder(ctx, purchaseOrderID)
	if !found {
		return sdk.Coin{}
	}
	return po.Amount
}

// GetPurchaseOrderStatus - get the Decision of a purchase order for a given purchaseOrderID
func (k Keeper) GetPurchaseOrderStatus(ctx sdk.Context, purchaseOrderID uint64) types.PurchaseOrderStatus {
	po, found := k.GetPurchaseOrder(ctx, purchaseOrderID)
	if !found {
		return types.StatusNil
	}
	return po.Status
}

// IteratePurchaseOrders iterates over the all the purchase orders and performs a callback function
func (k Keeper) IteratePurchaseOrders(ctx sdk.Context, cb func(purchaseOrder types.EnterpriseUndPurchaseOrder) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PurchaseOrderIDKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var po types.EnterpriseUndPurchaseOrder
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &po)

		if cb(po) {
			break
		}
	}
}

// GetAllPurchaseOrders returns all the purchase orders from store
func (k Keeper) GetAllPurchaseOrders(ctx sdk.Context) (purchaseOrders []types.EnterpriseUndPurchaseOrder) {
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
			matchPurchaser = po.Purchaser == params.Purchaser.String()
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
	if !types.ValidPurchaseOrderStatus(purchaseOrder.Status) {
		return sdkerrors.Wrap(types.ErrInvalidStatus, "unable to set purchase order - invalid status")
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.PurchaseOrderKey(purchaseOrder.Id), k.cdc.MustMarshalBinaryBare(&purchaseOrder))

	return nil
}

func (k Keeper) RaiseNewPurchaseOrder(ctx sdk.Context, purchaseOrder types.EnterpriseUndPurchaseOrder) (uint64, error) {

	logger := k.Logger(ctx)

	purchaseOrderId, err := k.GetHighestPurchaseOrderID(ctx)
	if err != nil {
		return 0, err
	}

	purchaseOrder.Id = purchaseOrderId
	purchaseOrder.Status = types.StatusRaised
	purchaseOrder.RaiseTime = uint64(ctx.BlockHeader().Time.Unix())

	err = k.SetPurchaseOrder(ctx, purchaseOrder)
	if err != nil {
		return 0, err
	}

	k.AddPoToRaisedQueue(ctx, purchaseOrderId)

	k.SetHighestPurchaseOrderID(ctx, purchaseOrderId+1)

	if !ctx.IsCheckTx() {
		logger.Debug("enterprise und purchase order raised", "id", purchaseOrderId, "from", purchaseOrder.Purchaser, "amt", purchaseOrder.Amount.String())
	}

	return purchaseOrderId, nil
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

	// todo - check found
	purchaseOrder, _ := k.GetPurchaseOrder(ctx, purchaseOrderID)

	poDecision := types.PurchaseOrderDecision{
		Signer:       signer.String(),
		Decision:     decision,
		DecisionTime: uint64(ctx.BlockHeader().Time.Unix()),
	}
	purchaseOrder.Decisions = append(purchaseOrder.Decisions, &poDecision)

	// update the status
	err := k.SetPurchaseOrder(ctx, purchaseOrder)

	if err != nil {
		return err
	}

	if !ctx.IsCheckTx() {
		logger.Debug("enterprise und purchase order decision made", "id", purchaseOrderID, "signer", signer, "decision", decision)
	}
	return nil
}
