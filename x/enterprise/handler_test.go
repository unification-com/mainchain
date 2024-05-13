package enterprise_test

import (
	"errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"strings"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := enterprise.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognised enterprise message type"))
}

func TestValidMsgUndPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[0])
	require.Nil(t, err)

	msg := &types.MsgUndPurchaseOrder{
		Purchaser: testAddrs[0].String(),
		Amount:    sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgUndPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[0])
	require.Nil(t, err)

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgUndPurchaseOrder
	}{
		{
			name:          "empty purchaser address",
			msg:           &types.MsgUndPurchaseOrder{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid purchaser address",
			msg: &types.MsgUndPurchaseOrder{
				Purchaser: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "invalid denomination",
			msg: &types.MsgUndPurchaseOrder{
				Purchaser: testAddrs[0].String(),
				Amount:    sdk.NewInt64Coin("rubbish", 100),
			},
			expectedError: sdkerrors.Wrap(types.ErrInvalidDenomination, fmt.Sprintf("denomination must be %s", test_helpers.TestDenomination)),
		},
		{
			name: "invalid amount",
			msg: &types.MsgUndPurchaseOrder{
				Purchaser: testAddrs[0].String(),
				Amount:    sdk.NewInt64Coin(test_helpers.TestDenomination, 0),
			},
			expectedError: sdkerrors.Wrap(types.ErrInvalidData, "amount must be > 0"),
		},
		{
			name: "purchaser not whitelisted",
			msg: &types.MsgUndPurchaseOrder{
				Purchaser: testAddrs[1].String(),
				Amount:    sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
			},
			expectedError: sdkerrors.Wrap(types.ErrNotAuthorisedToRaisePO, fmt.Sprintf("%s is not whitelisted to raise purchase orders", testAddrs[1].String())),
		},
		{
			name: "successful",
			msg: &types.MsgUndPurchaseOrder{
				Purchaser: testAddrs[0].String(),
				Amount:    sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h(ctx, tc.msg)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidMsgProcessUndPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	entSigners := app.EnterpriseKeeper.GetParamEntSignersAsAddressArray(ctx)
	entSignerAddr := entSigners[0]

	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[0])
	require.Nil(t, err)

	po := types.EnterpriseUndPurchaseOrder{
		Id:             1,
		Purchaser:      testAddrs[0].String(),
		Amount:         sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
		Status:         types.StatusRaised,
		RaiseTime:      uint64(time.Now().Unix()),
		CompletionTime: 0,
		Decisions:      nil,
	}
	err = app.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
	require.Nil(t, err)

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	msg := &types.MsgProcessUndPurchaseOrder{
		PurchaseOrderId: 1,
		Decision:        types.StatusAccepted,
		Signer:          entSignerAddr.String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestValidMsgProcessUndPurchaseOrderMultipleDecisions(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	entSignerAddr1 := testAddrs[0]
	entSignerAddr2 := testAddrs[1]

	signers := []string{entSignerAddr1.String(), entSignerAddr2.String()}

	app.EnterpriseKeeper.SetParams(ctx, types.Params{
		EntSigners:        strings.Join(signers, ","),
		Denom:             test_helpers.TestDenomination,
		MinAccepts:        2,
		DecisionTimeLimit: 1000,
	})

	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[0])
	require.Nil(t, err)

	po := types.EnterpriseUndPurchaseOrder{
		Id:             1,
		Purchaser:      testAddrs[2].String(),
		Amount:         sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
		Status:         types.StatusRaised,
		RaiseTime:      uint64(time.Now().Unix()),
		CompletionTime: 0,
		Decisions:      nil,
	}
	err = app.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
	require.Nil(t, err)

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	msg1 := &types.MsgProcessUndPurchaseOrder{
		PurchaseOrderId: 1,
		Decision:        types.StatusAccepted,
		Signer:          entSignerAddr1.String(),
	}

	res, err := h(ctx, msg1)

	require.NotNil(t, res)
	require.Nil(t, err)

	msg2 := &types.MsgProcessUndPurchaseOrder{
		PurchaseOrderId: 1,
		Decision:        types.StatusAccepted,
		Signer:          entSignerAddr2.String(),
	}

	res, err = h(ctx, msg2)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgProcessUndPurchaseOrder(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	entSigners := app.EnterpriseKeeper.GetParamEntSignersAsAddressArray(ctx)
	entSignerAddr := entSigners[0]

	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, testAddrs[0])
	require.Nil(t, err)

	po := types.EnterpriseUndPurchaseOrder{
		Id:             1,
		Purchaser:      testAddrs[0].String(),
		Amount:         sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
		Status:         types.StatusRaised,
		RaiseTime:      uint64(time.Now().Unix()),
		CompletionTime: 0,
		Decisions:      nil,
	}
	err = app.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
	require.Nil(t, err)

	// for testing purchaseOrder.Status != types.StatusRaised
	po1 := po
	po1.Id = 2
	po1.Status = types.StatusCompleted
	err = app.EnterpriseKeeper.SetPurchaseOrder(ctx, po1)
	require.Nil(t, err)

	// for testing repeated decisions
	po2 := po
	po2.Id = 3
	err = app.EnterpriseKeeper.SetPurchaseOrder(ctx, po2)
	require.Nil(t, err)

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgProcessUndPurchaseOrder
	}{
		{
			name:          "empty signer address",
			msg:           &types.MsgProcessUndPurchaseOrder{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid signer address",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "signer not authorised",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer processing purchase order"),
		},
		{
			name: "purchase order does not exist",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 99,
			},
			expectedError: sdkerrors.Wrap(types.ErrPurchaseOrderDoesNotExist, "id: 99"),
		},
		{
			name: "invalid decision",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 1,
				Decision:        types.StatusNil,
			},
			expectedError: sdkerrors.Wrap(types.ErrInvalidDecision, "decision should be accept or reject"),
		},
		{
			name: "po current status should only be raised",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 2,
				Decision:        types.StatusAccepted,
			},
			expectedError: sdkerrors.Wrapf(types.ErrPurchaseOrderAlreadyProcessed, "id %d already processed: %s", 2, types.StatusCompleted),
		},
		{
			name: "success - accept",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 1,
				Decision:        types.StatusAccepted,
			},
			expectedError: nil,
		},
		{
			name: "success - reject",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 3,
				Decision:        types.StatusRejected,
			},
			expectedError: nil,
		},

		{
			name: "signer cannot make more than one decision",
			msg: &types.MsgProcessUndPurchaseOrder{
				Signer:          entSignerAddr.String(),
				PurchaseOrderId: 3,
				Decision:        types.StatusAccepted,
			},
			expectedError: sdkerrors.Wrapf(types.ErrSignerAlreadyMadeDecision, "signer %s already decided: %s", entSignerAddr.String(), types.StatusRejected),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h(ctx, tc.msg)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidMsgWhitelistAddress(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	entSigners := app.EnterpriseKeeper.GetParamEntSignersAsAddressArray(ctx)
	entSignerAddr := entSigners[0]

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	msg := &types.MsgWhitelistAddress{
		Address: testAddrs[0].String(),
		Signer:  entSignerAddr.String(),
		Action:  types.WhitelistActionAdd,
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgWhitelistAddress(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)

	entSigners := app.EnterpriseKeeper.GetParamEntSignersAsAddressArray(ctx)
	entSignerAddr := entSigners[0]

	h := enterprise.NewHandler(app.EnterpriseKeeper)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgWhitelistAddress
	}{
		{
			name:          "empty signer address",
			msg:           &types.MsgWhitelistAddress{},
			expectedError: errors.New("signer address: empty address string is not allowed"),
		},
		{
			name: "invalid signer address",
			msg: &types.MsgWhitelistAddress{
				Signer: "rubbish",
			},
			expectedError: errors.New("signer address: decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "empty whitelist address",
			msg: &types.MsgWhitelistAddress{
				Signer: entSignerAddr.String(),
			},
			expectedError: errors.New("whitelist address: empty address string is not allowed"),
		},
		{
			name: "invalid whitelist address",
			msg: &types.MsgWhitelistAddress{
				Signer:  entSignerAddr.String(),
				Address: "rubbish",
			},
			expectedError: errors.New("whitelist address: decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "signer not authorised",
			msg: &types.MsgWhitelistAddress{
				Signer:  testAddrs[0].String(),
				Address: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist"),
		},
		{
			name: "invalid action",
			msg: &types.MsgWhitelistAddress{
				Signer:  entSignerAddr.String(),
				Address: testAddrs[0].String(),
				Action:  99,
			},
			expectedError: sdkerrors.Wrap(types.ErrInvalidDecision, "action should be add or remove"),
		},
		{
			name: "cannot remove non-existing address",
			msg: &types.MsgWhitelistAddress{
				Signer:  entSignerAddr.String(),
				Address: testAddrs[0].String(),
				Action:  types.WhitelistActionRemove,
			},
			expectedError: sdkerrors.Wrapf(types.ErrAddressNotWhitelisted, "%s not whitelisted", testAddrs[0].String()),
		},
		{
			name: "success",
			msg: &types.MsgWhitelistAddress{
				Signer:  entSignerAddr.String(),
				Address: testAddrs[0].String(),
				Action:  types.WhitelistActionAdd,
			},
			expectedError: nil,
		},
		{
			name: "cannot add address more than once",
			msg: &types.MsgWhitelistAddress{
				Signer:  entSignerAddr.String(),
				Address: testAddrs[0].String(),
				Action:  types.WhitelistActionAdd,
			},
			expectedError: sdkerrors.Wrapf(types.ErrAlreadyWhitelisted, "%s already whitelisted", testAddrs[0].String()),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h(ctx, tc.msg)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
