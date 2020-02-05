package keeper

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

func (k Keeper) TallyPurchaseOrderDecisions(ctx sdk.Context) {
	queryParams := types.NewQueryPurchaseOrdersParams(1, 1000, types.StatusRaised, sdk.AccAddress{})
	raisedPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)
	entParams := k.GetParams(ctx)
	timeNow := ctx.BlockHeader().Time.Unix()

	entSigners := strings.Split(entParams.EntSigners, ",")

	rejectThreshold := len(entSigners) - int(entParams.MinAccepts)

	logger := k.Logger(ctx)

	for _, po := range raisedPurchaseOrders {
		if po.Status != types.StatusRaised {
			panic("purchase order status is not raised!")
		}
		numAccepts := 0
		numRejects := 0
		for _, d := range po.Decisions {
			if d.Decision == types.StatusAccepted {
				numAccepts = numAccepts + 1
			}
			if d.Decision == types.StatusRejected {
				numRejects = numRejects + 1
			}
		}

		// first check if it's a stale PO
		timeDiff := timeNow - po.RaisedTime
		if timeDiff >= int64(entParams.DecisionLimit) && numAccepts < int(entParams.MinAccepts) {
			po.Status = types.StatusRejected
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("auto reject stale purchase order", "poid", po.PurchaseOrderID)
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeAutoRejectStalePurchaseOrder,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser.String()),
				),
			)
			continue
		}

		// check rejects
		if numRejects > rejectThreshold {
			po.Status = types.StatusRejected
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("purchase order rejected", "poid", po.PurchaseOrderID, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser.String()),
					sdk.NewAttribute(types.AttributeKeyDecision, po.Status.String()),
					sdk.NewAttribute(types.AttributeKeyNumAccepts, strconv.FormatUint(uint64(numAccepts), 10)),
					sdk.NewAttribute(types.AttributeKeyNumRejects, strconv.FormatUint(uint64(numRejects), 10)),
				),
			)
			continue
		}

		// check if there are enough accepts
		if numAccepts >= int(entParams.MinAccepts) {
			po.Status = types.StatusAccepted
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				panic(err)
			}
			logger.Info("purchase order accepted", "poid", po.PurchaseOrderID, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser.String()),
					sdk.NewAttribute(types.AttributeKeyDecision, po.Status.String()),
					sdk.NewAttribute(types.AttributeKeyNumAccepts, strconv.FormatUint(uint64(numAccepts), 10)),
					sdk.NewAttribute(types.AttributeKeyNumRejects, strconv.FormatUint(uint64(numRejects), 10)),
				),
			)
		}
	}
}

func (k Keeper) ProcessAcceptedPurchaseOrders(ctx sdk.Context) {
	queryParams := types.NewQueryPurchaseOrdersParams(1, 1000, types.StatusAccepted, sdk.AccAddress{})
	acceptedPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)
	logger := k.Logger(ctx)

	for _, po := range acceptedPurchaseOrders {
		if po.Status != types.StatusAccepted {
			panic("purchase order status is not accepted!")
		}

		// mark as completed
		po.Status = types.StatusCompleted
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
				types.EventTypeUndPurchaseComplete,
				sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.PurchaseOrderID, 10)),
				sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, po.Amount.String()),
			),
		)
	}
}
