package enterprise

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

// NewHandler returns a handler for "enterprise" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPurchaseUnd:
			return handleMsgPurchaseUnd(ctx, keeper, msg)
		case MsgProcessUndPurchaseOrder:
			return handleMsgProcessPurchaseUnd(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized enterprise Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgPurchaseUnd(ctx sdk.Context, k Keeper, msg MsgPurchaseUnd) sdk.Result {

	purchaseOrderID, err := k.RaiseNewPurchaseOrder(ctx, msg.Purchaser, msg.Amount)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRaisePurchaseOrder,
			sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(purchaseOrderID, 10)),
			sdk.NewAttribute(AttributeKeyPurchaser, msg.Purchaser.String()),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetPurchaseOrderIDBytes(purchaseOrderID),
	}
}

func handleMsgProcessPurchaseUnd(ctx sdk.Context, k Keeper, msg MsgProcessUndPurchaseOrder) sdk.Result {

	params := k.GetParams(ctx)

	// check only the Enterprise account is signing
	if !msg.Signer.Equals(params.EntSource) {
		return sdk.ErrUnauthorized("unauthorised signer processing purchase order").Result()
	}

	if !ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return ErrInvalidDecision(k.GetCodeSpace(), "decision should be accept or reject").Result()
	}

	err := k.ProcessPurchaseOrder(ctx, msg.PurchaseOrderID, msg.Decision)

	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeProcessPurchaseOrder,
			sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(msg.PurchaseOrderID, 10)),
			sdk.NewAttribute(AttributeKeyDecision, msg.Decision.String()),
		),
	})

	retData := append(GetPurchaseOrderIDBytes(msg.PurchaseOrderID), []byte{byte(msg.Decision)}...)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   retData,
	}
}
