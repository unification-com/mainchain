package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

// Tests for Highest Purchase Order ID

func TestSetGetHighestPurchaseOrderID(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	for i := uint64(1); i <= 1000; i++ {
		keeper.SetHighestPurchaseOrderID(ctx, i)
		poID, err := keeper.GetHighestPurchaseOrderID(ctx)
		require.NoError(t, err)
		require.True(t, poID == i)
	}
}

// Tests for Get/Set Purchase Order

func TestSetGetPurchaseOrder(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	status := types.StatusRaised

	for i := uint64(1); i <= 1000; i++ {
		po := types.NewEnterpriseUndPurchaseOrder()
		po.PurchaseOrderID = i
		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i))
		po.Purchaser = TestAddrs[1]
		po.Status = status

		err := keeper.SetPurchaseOrder(ctx, po)
		require.NoError(t, err)

		poExists := keeper.PurchaseOrderExists(ctx, i)
		require.True(t, poExists)

		poDb := keeper.GetPurchaseOrder(ctx, i)
		require.True(t, PurchaseOrderEqual(po, poDb))

		poStatus := keeper.GetPurchaseOrderStatus(ctx, i)
		require.True(t, poStatus == status)

		poFrom := keeper.GetPurchaseOrderPurchaser(ctx, i)
		require.True(t, poFrom.String() == TestAddrs[1].String())

		poAmount := keeper.GetPurchaseOrderAmount(ctx, i)
		require.True(t, poAmount.Denom == types.DefaultDenomination)
		require.True(t, poAmount.Amount.Int64() == int64(i))
	}

}

func TestSetEmptyPurchaseOrderValues(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	// Empty
	po1 := types.NewEnterpriseUndPurchaseOrder()

	// Only purchaser set
	po2 := po1
	po2.Purchaser = TestAddrs[1]

	// Only purchaser set, and Amount = 0
	po3 := po2
	po3.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 0)

	// Only purchaser and amount set, ID = 0
	po4 := po3
	po4.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 1000)

	po5 := po4
	po5.PurchaseOrderID = 1

	testCases := []struct {
		po          types.EnterpriseUndPurchaseOrder
		expectedErr sdk.Error
	}{
		{po1, sdk.ErrInternal("unable to raise purchase order - purchaser cannot be empty")},
		{po2, sdk.ErrInternal("unable to raise purchase order - amount not valid")},
		{po3, sdk.ErrInternal("unable to raise purchase order - amount must be positive")},
		{po4, sdk.ErrInternal("unable to raise purchase order - id must be positive non-zero")},
		{po5, nil},
	}

	for _, tc := range testCases {
		err := keeper.SetPurchaseOrder(ctx, tc.po)
		require.Equal(t, tc.expectedErr, err, "unexpected type of error: %s", err)
	}
}

// Tests for Raise new Purchase Order

func TestRaiseNewPurchaseOrder(t *testing.T) {

	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	testAddresses := GenerateRandomAddresses(100)

	i, _ := keeper.GetHighestPurchaseOrderID(ctx)

	for _, from := range testAddresses {
		amt := int64(RandInBetween(1, 10000))
		amount := sdk.NewInt64Coin(types.DefaultDenomination, amt)

		expectedPo := types.NewEnterpriseUndPurchaseOrder()
		expectedPo.PurchaseOrderID = i
		expectedPo.Amount = amount
		expectedPo.Purchaser = from
		expectedPo.Status = types.StatusRaised

		poID, err := keeper.RaiseNewPurchaseOrder(ctx, from, amount)
		require.NoError(t, err)
		require.True(t, poID == expectedPo.PurchaseOrderID)

		poExists := keeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		poDb := keeper.GetPurchaseOrder(ctx, poID)
		require.True(t, poDb.PurchaseOrderID == expectedPo.PurchaseOrderID)
		require.True(t, poDb.Status == types.StatusRaised)
		require.True(t, poDb.Purchaser.String() == from.String())
		require.True(t, poDb.Amount.Denom == types.DefaultDenomination)
		require.True(t, poDb.Amount.Amount.Int64() == amt)
		require.True(t, poDb.Amount.IsEqual(expectedPo.Amount))
		require.True(t, PurchaseOrderEqual(expectedPo, poDb))

		i = i + 1
	}
}

func TestFailRaiseNewPurchaseOrder(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	// Empty
	po1 := types.NewEnterpriseUndPurchaseOrder()

	// Only purchaser set
	po2 := po1
	po2.Purchaser = TestAddrs[1]

	// Only purchaser set, and Amount = 0
	po3 := po2
	po3.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 0)

	// Only purchaser set, and Amount = 0
	po4 := po3
	po4.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 1000)

	testCases := []struct {
		po           types.EnterpriseUndPurchaseOrder
		expectedErr  sdk.Error
		expectedPoID uint64
	}{
		{po1, sdk.ErrInternal("unable to raise purchase order - purchaser cannot be empty"), 0},
		{po2, sdk.ErrInternal("unable to raise purchase order - amount not valid"), 0},
		{po3, sdk.ErrInternal("unable to raise purchase order - amount must be positive"), 0},
		{po4, nil, 1},
	}

	for _, tc := range testCases {
		poID, err := keeper.RaiseNewPurchaseOrder(ctx, tc.po.Purchaser, tc.po.Amount)
		require.Equal(t, tc.expectedErr, err, "unexpected type of error: %s", err)
		require.True(t, poID == tc.expectedPoID)
	}
}

func TestHighestPurchaseOrderIdAfterRaise(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))
		from := TestAddrs[1]
		_, _ = keeper.RaiseNewPurchaseOrder(ctx, from, amount)

		nextID, _ := keeper.GetHighestPurchaseOrderID(ctx)
		expectedNextID := i+1
		require.True(t, nextID == expectedNextID)
	}
}

func TestPurchaseOrderExistsAfterRaise(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))
		from := TestAddrs[1]
		poID, _ := keeper.RaiseNewPurchaseOrder(ctx, from, amount)

		poExists := keeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		po := keeper.GetPurchaseOrder(ctx, poID)
		require.True(t, po.PurchaseOrderID == poID && po.PurchaseOrderID == i)
	}
}

// Tests for processing Purchase Orders

func TestProcessPurchaseOrderAfterRaise(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))
		from := TestAddrs[1]
		poID, _ := keeper.RaiseNewPurchaseOrder(ctx, from, amount)
		decision := RandomDecision()

		poExists := keeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		err := keeper.ProcessPurchaseOrder(ctx, poID, decision)
		require.NoError(t, err)

		po := keeper.GetPurchaseOrder(ctx, poID)
		require.True(t, po.Status == decision)
	}
}
