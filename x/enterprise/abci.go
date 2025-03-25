package enterprise

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	err := k.ProcessAcceptedPurchaseOrders(sdkCtx)

	if err != nil {
		return err
	}

	err = k.TallyPurchaseOrderDecisions(sdkCtx)

	if err != nil {
		return err
	}

	return nil
}
