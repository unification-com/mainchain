package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
)

// Tests for Highest Purchase Order ID

func TestSetGetHighestPurchaseOrderID(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	for i := uint64(1); i <= 1000; i++ {
		app.EnterpriseKeeper.SetHighestPurchaseOrderID(ctx, i)
		poID, err := app.EnterpriseKeeper.GetHighestPurchaseOrderID(ctx)
		require.NoError(t, err)
		require.True(t, poID == i)
	}
}

// Tests for Get/Set Purchase Order

func TestSetGetPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[1])

	status := types.StatusRaised

	for i := uint64(1); i <= 1000; i++ {
		po := types.EnterpriseUndPurchaseOrder{}
		po.Id = i
		po.Amount = sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i))

		purchaser := testAddrs[1]
		// should still be able to SetPurchaseOrder if address is not whitelisted.
		if i > 500 {
			purchaser = testAddrs[2]
		}
		po.Purchaser = purchaser.String()
		po.Status = status
		po.RaiseTime = uint64(ctx.BlockHeader().Time.Unix())

		err := app.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
		require.NoError(t, err)

		poExists := app.EnterpriseKeeper.PurchaseOrderExists(ctx, i)
		require.True(t, poExists)

		poDb, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, i)
		require.True(t, found)
		require.True(t, po.Id == poDb.Id)
		require.True(t, po.RaiseTime == poDb.RaiseTime)
		require.True(t, po.Amount.String() == poDb.Amount.String())
		require.True(t, po.Purchaser == poDb.Purchaser)

		poStatus := app.EnterpriseKeeper.GetPurchaseOrderStatus(ctx, i)
		require.True(t, poStatus == status)

		poFrom := app.EnterpriseKeeper.GetPurchaseOrderPurchaser(ctx, i)
		require.True(t, poFrom.String() == purchaser.String())

		poAmount := app.EnterpriseKeeper.GetPurchaseOrderAmount(ctx, i)
		require.True(t, poAmount.Denom == test_helpers.TestDenomination)
		require.True(t, poAmount.Amount.Int64() == int64(i))
	}

}

// Tests for Raise new Purchase Order

func TestRaiseNewPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(100)

	i, _ := app.EnterpriseKeeper.GetHighestPurchaseOrderID(ctx)

	for _, from := range testAddrs {
		amt := int64(test_helpers.RandInBetween(1, 10000))
		amount := sdk.NewInt64Coin(test_helpers.TestDenomination, amt)

		_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, from)

		expectedPo := types.EnterpriseUndPurchaseOrder{}
		expectedPo.Id = i
		expectedPo.Amount = amount
		expectedPo.Purchaser = from.String()
		expectedPo.Status = types.StatusRaised
		expectedPo.RaiseTime = uint64(ctx.BlockHeader().Time.Unix())

		poID, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, expectedPo)
		require.NoError(t, err)
		require.True(t, poID == expectedPo.Id)

		poExists := app.EnterpriseKeeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		poDb, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poID)
		require.True(t, found)

		require.True(t, poDb.Id == expectedPo.Id)
		require.True(t, poDb.Status == types.StatusRaised)
		require.True(t, poDb.Purchaser == from.String())
		require.True(t, poDb.Amount.Denom == test_helpers.TestDenomination)
		require.True(t, poDb.Amount.Amount.Int64() == amt)
		require.True(t, poDb.Amount.IsEqual(expectedPo.Amount))

		i = i + 1
	}
}

func TestHighestPurchaseOrderIdAfterRaise(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)
	from := testAddrs[0]

	_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i))
		po := types.EnterpriseUndPurchaseOrder{}
		po.Amount = amount
		po.Purchaser = from.String()

		_, _ = app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, po)

		nextID, _ := app.EnterpriseKeeper.GetHighestPurchaseOrderID(ctx)
		expectedNextID := i + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestPurchaseOrderExistsAfterRaise(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)
	from := testAddrs[0]
	_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i))
		po := types.EnterpriseUndPurchaseOrder{}
		po.Amount = amount
		po.Purchaser = from.String()

		poID, _ := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, po)

		poExists := app.EnterpriseKeeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		poDb, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poID)
		require.True(t, found)
		require.True(t, poDb.Id == poID && poDb.Id == i)
	}
}

// Tests for processing Purchase Orders

func TestProcessPurchaseOrderAfterRaise(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	entSigners := app.EnterpriseKeeper.GetParamEntSignersAsAddressArray(ctx)
	entSignerAddr := entSigners[0]

	from := testAddrs[0]
	_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i))
		po := types.EnterpriseUndPurchaseOrder{}
		po.Amount = amount
		po.Purchaser = from.String()

		poID, _ := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, po)
		decision := RandomDecision()

		poExists := app.EnterpriseKeeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		err := app.EnterpriseKeeper.ProcessPurchaseOrderDecision(ctx, poID, decision, entSignerAddr)
		require.NoError(t, err)

		poDb, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poID)
		require.True(t, found)

		require.True(t, AddressInDecisions(entSignerAddr, poDb.Decisions))

		for _, d := range poDb.Decisions {
			if d.Signer == entSignerAddr.String() {
				require.True(t, d.Decision == decision)
			}
		}
	}
}

func TestRaisedQueue(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	for i := uint64(1); i < 10000; i++ {
		isInQueue := app.EnterpriseKeeper.PurchaseOrderIsInRaisedQueue(ctx, i)
		require.False(t, isInQueue)

		app.EnterpriseKeeper.AddPoToRaisedQueue(ctx, i)
		isInQueue = app.EnterpriseKeeper.PurchaseOrderIsInRaisedQueue(ctx, i)
		require.True(t, isInQueue)

		app.EnterpriseKeeper.RemovePurchaseOrderFromRaisedQueue(ctx, i)
		isInQueue = app.EnterpriseKeeper.PurchaseOrderIsInRaisedQueue(ctx, i)
		require.False(t, isInQueue)
	}
}

func TestPurchaseOrderAddedToRaisedQueue(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	for _, from := range testAddrs {
		for i := uint64(1); i < 1000; i++ {
			amount := sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i))
			po := types.EnterpriseUndPurchaseOrder{}
			po.Amount = amount
			po.Purchaser = from.String()

			poID, _ := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, po)

			isInQueue := app.EnterpriseKeeper.PurchaseOrderIsInRaisedQueue(ctx, poID)
			require.True(t, isInQueue)
		}
	}
}

func TestRaisedQueueIterator(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	toAdd := []uint64{2, 23, 99, 12345, 6236285, 4, 7, 9, 24}

	for _, poId := range toAdd {
		app.EnterpriseKeeper.AddPoToRaisedQueue(ctx, poId)
		isInQueue := app.EnterpriseKeeper.PurchaseOrderIsInRaisedQueue(ctx, poId)
		require.True(t, isInQueue)
	}

	raisedPos := app.EnterpriseKeeper.GetAllRaisedPurchaseOrders(ctx)

	require.True(t, len(raisedPos) == len(toAdd))

	for _, poId := range toAdd {
		require.True(t, isInList(raisedPos, poId))
	}
}

func TestAcceptedQueue(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	for i := uint64(1); i < 10000; i++ {
		isInQueue := app.EnterpriseKeeper.PurchaseOrderIsInAcceptedQueue(ctx, i)
		require.False(t, isInQueue)

		app.EnterpriseKeeper.AddPoToAcceptedQueue(ctx, i)
		isInQueue = app.EnterpriseKeeper.PurchaseOrderIsInAcceptedQueue(ctx, i)
		require.True(t, isInQueue)

		app.EnterpriseKeeper.RemovePurchaseOrderFromAcceptedQueue(ctx, i)
		isInQueue = app.EnterpriseKeeper.PurchaseOrderIsInAcceptedQueue(ctx, i)
		require.False(t, isInQueue)
	}
}

func TestAcceptedQueueIterator(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	toAdd := []uint64{2, 23, 99, 12345, 6236285, 4, 7, 9, 24}

	for _, poId := range toAdd {
		app.EnterpriseKeeper.AddPoToAcceptedQueue(ctx, poId)
		isInQueue := app.EnterpriseKeeper.PurchaseOrderIsInAcceptedQueue(ctx, poId)
		require.True(t, isInQueue)
	}

	acceptedPos := app.EnterpriseKeeper.GetAllAcceptedPurchaseOrders(ctx)

	require.True(t, len(acceptedPos) == len(toAdd))

	for _, poId := range toAdd {
		require.True(t, isInList(acceptedPos, poId))
	}
}

func isInList(list []uint64, i uint64) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}
