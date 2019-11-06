package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
	"strconv"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {

	acceptedPurchaseOrders := k.GetAllAcceptedPurchaseOrdersIterator(ctx)

	for ; acceptedPurchaseOrders.Valid(); acceptedPurchaseOrders.Next() {
		var po types.EnterpriseUndPurchaseOrder
		k.GetCdc().MustUnmarshalBinaryBare(acceptedPurchaseOrders.Value(), &po)
		// first delete
		k.DeleteAcceptedPurchaseOrder(ctx, po.PurchaseOrderID)

		if po.Status != types.StatusAccepted {
			panic("purchase order status is not accepted!")
		}

		// Mint the Enterprise UND
		mintedCoins := sdk.NewCoins(po.Amount)
		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}

		// Send them to the purchaser's account
		err = k.SendCoins(ctx, po.Purchaser, mintedCoins)
		if err != nil {
			panic(err)
		}

		// keep track of how much UND is locked
		err = k.IncrementLockedUnd(ctx, po.Purchaser, po.Amount)
		if err != nil {
			panic(err)
		}

		// Delegate the Enterprise UND so they can't be spent
		err = k.SupplyKeeper.DelegateCoinsFromAccountToModule(ctx, po.Purchaser, types.ModuleName, mintedCoins)
		if err != nil {
			panic(err)
		}

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
