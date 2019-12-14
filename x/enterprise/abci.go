package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	processAcceptedPurchaseOrders(ctx, k)
	tallyPurchaseOrderDecisions(ctx, k)
}

func tallyPurchaseOrderDecisions(ctx sdk.Context, k Keeper) {
	queryParams := NewQueryPurchaseOrdersParams(1, 1000, StatusRaised, sdk.AccAddress{})
	raisedPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)
	entParams := k.GetParams(ctx)
	timeNow := ctx.BlockHeader().Time.Unix()

	rejectThreshold := len(entParams.EntSigners) - int(entParams.MinAccepts)

	logger := k.Logger(ctx)

	for _, po := range raisedPurchaseOrders {
		if po.Status != StatusRaised {
			panic("purchase order status is not raised!")
		}
		numAccepts := 0
		numRejects := 0
		for _, d := range po.Decisions {
			if d.Decision == StatusAccepted {
				numAccepts = numAccepts + 1
			}
			if d.Decision == StatusRejected {
				numRejects = numRejects + 1
			}
		}

		// first check if it's a stale PO
		timeDiff := timeNow - po.RaisedTime
		if timeDiff >= int64(entParams.DecisionLimit) && numAccepts < int(entParams.MinAccepts) {
			po.Status = StatusRejected
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("auto reject stale purchase order", "poid", po.PurchaseOrderID)
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					EventTypeAutoRejectStalePurchaseOrder,
					sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(AttributeKeyPurchaser, po.Purchaser.String()),
				),
			)
			continue
		}

		// check rejects
		if numRejects > rejectThreshold {
			po.Status = StatusRejected
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("purchase order rejected", "poid", po.PurchaseOrderID, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(AttributeKeyPurchaser, po.Purchaser.String()),
					sdk.NewAttribute(AttributeKeyDecision, po.Status.String()),
					sdk.NewAttribute(AttributeKeyNumAccepts, strconv.FormatUint(uint64(numAccepts), 10)),
					sdk.NewAttribute(AttributeKeyNumRejects, strconv.FormatUint(uint64(numRejects), 10)),
				),
			)
			continue
		}

		// check if there are enough accepts
		if numAccepts >= int(entParams.MinAccepts) {
			po.Status = StatusAccepted
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("purchase order accepted", "poid", po.PurchaseOrderID, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(AttributeKeyPurchaser, po.Purchaser.String()),
					sdk.NewAttribute(AttributeKeyDecision, po.Status.String()),
					sdk.NewAttribute(AttributeKeyNumAccepts, strconv.FormatUint(uint64(numAccepts), 10)),
					sdk.NewAttribute(AttributeKeyNumRejects, strconv.FormatUint(uint64(numRejects), 10)),
				),
			)
		}

	}
}

func processAcceptedPurchaseOrders(ctx sdk.Context, k Keeper) {
	queryParams := NewQueryPurchaseOrdersParams(1, 1000, StatusAccepted, sdk.AccAddress{})
	acceptedPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)
	logger := k.Logger(ctx)

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

		logger.Info("purchase order complete", "poid", po.PurchaseOrderID, "status", po.Status)

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
