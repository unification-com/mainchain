package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
	"testing"
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

	_ = keeper.AddAddressToWhitelist(ctx, TestAddrs[1])

	status := types.StatusRaised

	for i := uint64(1); i <= 1000; i++ {
		po := types.NewEnterpriseUndPurchaseOrder()
		po.PurchaseOrderID = i
		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i))
		po.Purchaser = TestAddrs[1]
		po.Status = status
		po.RaisedTime = ctx.BlockHeader().Time.Unix()

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
	po1.Status = types.StatusRaised

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

	po6 := po5
	po6.Status = 0x05

	po7 := po6
	po7.Status = types.StatusNil

	_ = keeper.AddAddressToWhitelist(ctx, TestAddrs[1])

	testCases := []struct {
		po          types.EnterpriseUndPurchaseOrder
		expectedErr error
	}{
		{po1, sdkerrors.Wrap(types.ErrMissingData, "unable to set purchase order - purchaser cannot be empty")},
		{po2, sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount not valid")},
		{po3, sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount must be positive")},
		{po4, sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - id must be positive non-zero")},
		{po5, nil},
		{po6, sdkerrors.Wrap(types.ErrInvalidStatus, "unable to set purchase order - invalid status")},
		{po7, sdkerrors.Wrap(types.ErrInvalidStatus, "unable to set purchase order - invalid status")},
	}

	for _, tc := range testCases {
		err := keeper.SetPurchaseOrder(ctx, tc.po)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
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

		_ = keeper.AddAddressToWhitelist(ctx, from)

		expectedPo := types.NewEnterpriseUndPurchaseOrder()
		expectedPo.PurchaseOrderID = i
		expectedPo.Amount = amount
		expectedPo.Purchaser = from
		expectedPo.Status = types.StatusRaised
		expectedPo.RaisedTime = ctx.BlockHeader().Time.Unix()

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

	_ = keeper.AddAddressToWhitelist(ctx,  TestAddrs[1])

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
		remove       bool
		expectedErr  error
		expectedPoID uint64
	}{
		{po1, false, sdkerrors.Wrap(types.ErrMissingData, "unable to set purchase order - purchaser cannot be empty"), 0},
		{po2, false, sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount not valid"), 0},
		{po3, false, sdkerrors.Wrap(types.ErrInvalidData, "unable to set purchase order - amount must be positive"), 0},
		{po4, false, nil, 1},
		{po4, true, sdkerrors.Wrap(types.ErrNotAuthorisedToRaisePO,  fmt.Sprintf("%s is not whitelisted to raise purchase orders", TestAddrs[1])), 0},
	}

	for _, tc := range testCases {
		if tc.remove {
			_ = keeper.RemoveAddressFromWhitelist(ctx, TestAddrs[1])
		}
		poID, err := keeper.RaiseNewPurchaseOrder(ctx, tc.po.Purchaser, tc.po.Amount)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
		require.True(t, poID == tc.expectedPoID)
	}
}

func TestHighestPurchaseOrderIdAfterRaise(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	from := TestAddrs[1]

	_ = keeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))

		_, _ = keeper.RaiseNewPurchaseOrder(ctx, from, amount)

		nextID, _ := keeper.GetHighestPurchaseOrderID(ctx)
		expectedNextID := i + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestPurchaseOrderExistsAfterRaise(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	from := TestAddrs[1]
	_ = keeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))

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
	from := TestAddrs[1]
	_ = keeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))

		poID, _ := keeper.RaiseNewPurchaseOrder(ctx, from, amount)
		decision := RandomDecision()

		poExists := keeper.PurchaseOrderExists(ctx, poID)
		require.True(t, poExists)

		err := keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
		require.NoError(t, err)

		po := keeper.GetPurchaseOrder(ctx, poID)

		require.True(t, AddressInDecisions(EntSignerAddr, po.Decisions))

		for _, d := range po.Decisions {
			if d.Signer.Equals(EntSignerAddr) {
				require.True(t, d.Decision == decision)
			}
		}
	}
}

func TestProcessNotExistPurchaseOrder(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	for i := uint64(1); i < 1000; i++ {
		err := keeper.ProcessPurchaseOrderDecision(ctx, i, RandomDecision(), EntSignerAddr)
		errMsg := fmt.Sprintf("id: %d", i)
		expectedErr := sdkerrors.Wrap(types.ErrPurchaseOrderDoesNotExist, errMsg)
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestProcessingDuplicatePurchaseOrders(t *testing.T) {

	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	from := TestAddrs[1]
	_ = keeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))

		poID, _ := keeper.RaiseNewPurchaseOrder(ctx, from, amount)
		decision := RandomDecision()
		err := keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
		require.NoError(t, err)

		// reprocess
		errMsg := fmt.Sprintf("id %d already processed: %s", poID, decision)
		expectedErr := sdkerrors.Wrap(types.ErrPurchaseOrderAlreadyProcessed, errMsg)

		// mock blocker processing
		po := keeper.GetPurchaseOrder(ctx, poID)
		po.Status = decision
		_ = keeper.SetPurchaseOrder(ctx, po)

		err = keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

		// mock complete
		if decision == types.StatusAccepted {
			po := keeper.GetPurchaseOrder(ctx, poID)
			po.Status = types.StatusCompleted
			_ = keeper.SetPurchaseOrder(ctx, po)

			errMsg := fmt.Sprintf("id %d already processed: complete", poID)
			expectedErr := sdkerrors.Wrap(types.ErrPurchaseOrderAlreadyProcessed, errMsg)

			err = keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
			require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
		} else {
			errMsg := fmt.Sprintf("id %d already processed: reject", poID)
			expectedErr := sdkerrors.Wrap(types.ErrPurchaseOrderAlreadyProcessed, errMsg)

			err = keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
			require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
		}
	}
}

func TestProcessingDuplicateDecisions(t *testing.T) {

	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	from := TestAddrs[1]
	_ = keeper.AddAddressToWhitelist(ctx, from)

	for i := uint64(1); i < 1000; i++ {
		amount := sdk.NewInt64Coin(types.DefaultDenomination, int64(i))

		poID, _ := keeper.RaiseNewPurchaseOrder(ctx, from, amount)
		decision := RandomDecision()
		err := keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
		require.NoError(t, err)

		// reprocess
		errMsg := fmt.Sprintf("signer %s already decided: %s", EntSignerAddr, decision)
		expectedErr := sdkerrors.Wrap(types.ErrSignerAlreadyMadeDecision, errMsg)

		err = keeper.ProcessPurchaseOrderDecision(ctx, poID, decision, EntSignerAddr)
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestProcessPurchaseOrderInvalidDecision(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	_ = keeper.AddAddressToWhitelist(ctx, TestAddrs[0])

	po := types.NewEnterpriseUndPurchaseOrder()
	po.Status = types.StatusRaised
	po.PurchaseOrderID = 1
	po.Amount = sdk.NewInt64Coin(TestDenomination, 10000)
	po.Purchaser = TestAddrs[0]

	_ = keeper.SetPurchaseOrder(ctx, po)

	po.PurchaseOrderID = 2
	_ = keeper.SetPurchaseOrder(ctx, po)

	testCases := []struct {
		poId        uint64
		decision    types.PurchaseOrderStatus
		expectedErr error
	}{
		{1, types.StatusRaised, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, types.StatusCompleted, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, types.StatusNil, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, 0x05, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, 0x06, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, 0x07, sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")},
		{1, types.StatusAccepted, nil},
		{2, types.StatusRejected, nil},
	}

	for _, tc := range testCases {
		err := keeper.ProcessPurchaseOrderDecision(ctx, tc.poId, tc.decision, EntSignerAddr)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
	}
}

func TestUnauthorisedDecisionMaker(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	_ = keeper.AddAddressToWhitelist(ctx, TestAddrs[0])
	_ = keeper.AddAddressToWhitelist(ctx, TestAddrs[1])

	po := types.NewEnterpriseUndPurchaseOrder()
	po.Status = types.StatusRaised
	po.PurchaseOrderID = 1
	po.Amount = sdk.NewInt64Coin(TestDenomination, 10000)
	po.Purchaser = TestAddrs[0]

	_ = keeper.SetPurchaseOrder(ctx, po)

	po.PurchaseOrderID = 2
	po.Purchaser = TestAddrs[1]
	_ = keeper.SetPurchaseOrder(ctx, po)

	testCases := []struct {
		poId        uint64
		decision    types.PurchaseOrderStatus
		signer      sdk.AccAddress
		expectedErr error
	}{
		{1, types.StatusAccepted, TestAddrs[0], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{1, types.StatusAccepted, TestAddrs[1], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{1, types.StatusAccepted, TestAddrs[2], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{1, types.StatusRejected, TestAddrs[3], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{1, types.StatusRejected, TestAddrs[4], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{2, types.StatusAccepted, TestAddrs[0], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{2, types.StatusAccepted, TestAddrs[1], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{2, types.StatusAccepted, TestAddrs[2], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{2, types.StatusRejected, TestAddrs[3], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{2, types.StatusRejected, TestAddrs[4], sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order")},
		{1, types.StatusAccepted, EntSignerAddr, nil},
		{2, types.StatusRejected, EntSignerAddr, nil},
	}

	for _, tc := range testCases {
		err := keeper.ProcessPurchaseOrderDecision(ctx, tc.poId, tc.decision, tc.signer)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
	}
}
