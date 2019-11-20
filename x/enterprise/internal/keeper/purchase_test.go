package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

func TestSetGetPurchaseOrder(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	po := types.NewEnterpriseUndPurchaseOrder()
	po.PurchaseOrderID = 1
	po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 1000000)
	po.Purchaser = TestAddrs[1]
	po.Status = types.StatusRaised

	err := keeper.SetPurchaseOrder(ctx, po)
	require.NoError(t, err)

	poDb := keeper.GetPurchaseOrder(ctx, po.PurchaseOrderID)

	require.True(t, PurchaseOrderEqual(po, poDb))

}
