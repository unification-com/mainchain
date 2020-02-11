package enterprise

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// NewHandler returns a handler for "enterprise" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgPurchaseUnd:
			return handleMsgPurchaseUnd(ctx, keeper, msg)
		case MsgProcessUndPurchaseOrder:
			return handleMsgProcessPurchaseUnd(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgPurchaseUnd(ctx sdk.Context, k Keeper, msg MsgPurchaseUnd) (*sdk.Result, error) {

	if msg.Amount.Denom != k.GetParamDenom(ctx) {
		return nil, sdkerrors.Wrap(ErrInvalidDenomination, fmt.Sprintf("denomination must be %s", k.GetParamDenom(ctx)))
	}

	purchaseOrderID, err := k.RaiseNewPurchaseOrder(ctx, msg.Purchaser, msg.Amount)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeRaisePurchaseOrder,
			sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(purchaseOrderID, 10)),
			sdk.NewAttribute(AttributeKeyPurchaser, msg.Purchaser.String()),
			sdk.NewAttribute(AttributeKeyAmount, msg.Amount.String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   GetPurchaseOrderIDBytes(purchaseOrderID),
	}, nil
}

func handleMsgProcessPurchaseUnd(ctx sdk.Context, k Keeper, msg MsgProcessUndPurchaseOrder) (*sdk.Result, error) {

	// check only authorised Enterprise account is signing
	if !k.IsAuthorisedToDecide(ctx, msg.Signer) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")
	}

	if !ValidPurchaseOrderAcceptRejectStatus(msg.Decision) {
		return nil, sdkerrors.Wrap(ErrInvalidDecision, "decision should be accept or reject")
	}

	err := k.ProcessPurchaseOrderDecision(ctx, msg.PurchaseOrderID, msg.Decision, msg.Signer)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeProcessPurchaseOrderDecision,
			sdk.NewAttribute(AttributeKeyPurchaseOrderID, strconv.FormatUint(msg.PurchaseOrderID, 10)),
			sdk.NewAttribute(AttributeKeySigner, msg.Signer.String()),
			sdk.NewAttribute(AttributeKeyDecision, msg.Decision.String()),
		),
	})

	retData := append(GetPurchaseOrderIDBytes(msg.PurchaseOrderID), []byte{byte(msg.Decision)}...)

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
		Data:   retData,
	}, nil
}
