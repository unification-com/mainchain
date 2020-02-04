package enterprise

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.ProcessAcceptedPurchaseOrders(ctx)
	k.TallyPurchaseOrderDecisions(ctx)
}
