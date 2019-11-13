package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	var purchaserAddr sdk.AccAddress
	queryParams := NewQueryPurchaseOrdersParams(1, 1000, StatusAccepted, purchaserAddr)
	acceptedPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)

	for _, po := range acceptedPurchaseOrders {
		if po.Status != StatusAccepted {
			panic("purchase order status is not accepted!")
		}

		// mark as completed
		po.Status = StatusCompleted
		err := k.SetPurchaseOrder(ctx, po)
		if err != nil {
			panic(err)
		}

		// Mint the Enterprise UND
		err = k.MintCoinsAndLock(ctx, po.Purchaser, po.Amount)
		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				EventTypeUndPurchaseComplete,
				sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
				sdk.NewAttribute(AttributeKeyPurchaser, po.Purchaser.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, po.Amount.String()),
			),
		)

	}
}
