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
			return handleMMsgPurchaseUnd(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized enterprise Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMMsgPurchaseUnd(ctx sdk.Context, keeper Keeper, msg MsgPurchaseUnd) sdk.Result {

	purchaseOrderID, err := keeper.RaiseNewPurchaseOrder(ctx, msg.Purchaser, msg.Amount)

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
