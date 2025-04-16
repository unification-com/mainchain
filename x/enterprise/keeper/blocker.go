package keeper

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

func (k Keeper) TallyPurchaseOrderDecisions(ctx sdk.Context) error {

	raisedPurchaseOrderIds := k.GetAllRaisedPurchaseOrders(ctx)
	entParams := k.GetParams(ctx)
	timeNow := uint64(ctx.BlockHeader().Time.Unix())

	entSigners := strings.Split(entParams.EntSigners, ",")

	rejectThreshold := len(entSigners) - int(entParams.MinAccepts)

	logger := k.Logger(ctx)

	for _, poId := range raisedPurchaseOrderIds {
		po, found := k.GetPurchaseOrder(ctx, poId)
		if !found {
			logger.Warn("purchase order not found in abci method TallyPurchaseOrderDecisions", "poid", poId)
			continue
		}
		if po.Status != types.StatusRaised {
			logger.Warn("purchase order status is not raised in abci method TallyPurchaseOrderDecisions", "poid", poId)
			continue
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
		timeDiff := timeNow - po.RaiseTime
		if timeDiff >= entParams.DecisionTimeLimit && numAccepts < int(entParams.MinAccepts) {
			po.Status = types.StatusRejected
			po.CompletionTime = timeNow
			err := k.SetPurchaseOrder(ctx, po)
			if err != nil {
				return err
			}
			if !ctx.IsCheckTx() {
				logger.Debug("auto reject stale purchase order", "poid", po.Id)
			}

			// remove from raised queue
			k.RemovePurchaseOrderFromRaisedQueue(ctx, poId)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeAutoRejectStalePurchaseOrder,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.Id, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser),
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
				return err
			}
			if !ctx.IsCheckTx() {
				logger.Debug("purchase order rejected", "poid", po.Id, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())
			}

			// remove from raised queue
			k.RemovePurchaseOrderFromRaisedQueue(ctx, poId)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.Id, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser),
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
				return err
			}
			if !ctx.IsCheckTx() {
				logger.Debug("purchase order accepted", "poid", po.Id, "accepts", numAccepts, "rejects", numRejects, "decision", po.Status.String())
			}

			// remove from raised queue
			k.RemovePurchaseOrderFromRaisedQueue(ctx, poId)

			// add to accepted queue
			k.AddPoToAcceptedQueue(ctx, poId)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeTallyPurchaseOrderDecisions,
					sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.Id, 10)),
					sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser),
					sdk.NewAttribute(types.AttributeKeyDecision, po.Status.String()),
					sdk.NewAttribute(types.AttributeKeyNumAccepts, strconv.FormatUint(uint64(numAccepts), 10)),
					sdk.NewAttribute(types.AttributeKeyNumRejects, strconv.FormatUint(uint64(numRejects), 10)),
				),
			)
		}
	}

	return nil
}

func (k Keeper) ProcessAcceptedPurchaseOrders(ctx sdk.Context) error {
	acceptedPurchaseOrderIds := k.GetAllAcceptedPurchaseOrders(ctx)
	logger := k.Logger(ctx)

	for _, poId := range acceptedPurchaseOrderIds {
		po, found := k.GetPurchaseOrder(ctx, poId)
		if !found {
			logger.Warn("purchase order not found in abci method ProcessAcceptedPurchaseOrders", "poid", poId)
			continue
		}
		if po.Status != types.StatusAccepted {
			logger.Warn("purchase order status is not accepted in abci method ProcessAcceptedPurchaseOrders", "poid", poId)
			continue
		}

		// mark as completed
		po.Status = types.StatusCompleted
		err := k.SetPurchaseOrder(ctx, po)
		if err != nil {
			return err
		}

		purchaser, err := sdk.AccAddressFromBech32(po.Purchaser)
		if err != nil {
			return err
		}

		// Mint the Enterprise FUND
		err = k.CreateAndLockEFUND(ctx, purchaser, po.Amount)
		if err != nil {
			return err
		}

		if !ctx.IsCheckTx() {
			logger.Debug("purchase order complete", "poid", po.Id, "status", po.Status)
		}

		// remove from the accepted queue
		k.RemovePurchaseOrderFromAcceptedQueue(ctx, poId)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeUndPurchaseComplete,
				sdk.NewAttribute(types.AttributeKeyPurchaseOrderID, strconv.FormatUint(po.Id, 10)),
				sdk.NewAttribute(types.AttributeKeyPurchaser, po.Purchaser),
				sdk.NewAttribute(sdk.AttributeKeyAmount, po.Amount.String()),
			),
		)
	}

	return nil
}
