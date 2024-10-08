package beacon_test

import (
	"errors"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/beacon/types"
	"strings"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon"
	"github.com/unification-com/mainchain/x/beacon/keeper"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := beacon.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognised beacon message type"))
}

func TestValidMsgRegisterBeacon(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := beacon.NewHandler(app.BeaconKeeper)

	msg := &types.MsgRegisterBeacon{
		Moniker: "moniker",
		Name:    "name",
		Owner:   testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgRegisterBeacon(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := beacon.NewHandler(app.BeaconKeeper)

	existsMoniker := simapp.GenerateRandomString(24)
	b := types.Beacon{
		Moniker: existsMoniker,
		Name:    "this exists",
		Owner:   testAddrs[1].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, b)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRegisterBeacon
	}{
		{
			name: "empty owner address",
			msg: &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name:    "name",
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name:    "name",
				Owner:   "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "name too big",
			msg: &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name:    simapp.GenerateRandomString(129),
				Owner:   testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit"),
		},
		{
			name: "moniker too big",
			msg: &types.MsgRegisterBeacon{
				Moniker: simapp.GenerateRandomString(65),
				Name:    "name",
				Owner:   testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit"),
		},
		{
			name: "zero length moniker",
			msg: &types.MsgRegisterBeacon{
				Moniker: "",
				Name:    "name",
				Owner:   testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker"),
		},
		{
			name: "successful",
			msg: &types.MsgRegisterBeacon{
				Moniker: simapp.GenerateRandomString(24),
				Name:    simapp.GenerateRandomString(24),
				Owner:   testAddrs[0].String(),
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

func TestValidMsgRecordBeaconTimestamp(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := beacon.NewHandler(app.BeaconKeeper)

	b := types.Beacon{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new beacon",
		Owner:   testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, b)
	require.Nil(t, err)

	msg := &types.MsgRecordBeaconTimestamp{
		BeaconId: 1,
		Hash:     simapp.GenerateRandomString(64),
		Owner:    testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgRecordBeaconTimestamp(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := beacon.NewHandler(app.BeaconKeeper)

	b := types.Beacon{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new beacon",
		Owner:   testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, b)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRecordBeaconTimestamp
	}{
		{
			name:          "empty owner address",
			msg:           &types.MsgRecordBeaconTimestamp{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgRecordBeaconTimestamp{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "hash too large",
			msg: &types.MsgRecordBeaconTimestamp{
				Owner: testAddrs[0].String(),
				Hash:  simapp.GenerateRandomString(67),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "hash too big. 66 character limit"),
		},
		{
			name: "beacon not registered",
			msg: &types.MsgRecordBeaconTimestamp{
				Owner:    testAddrs[0].String(),
				Hash:     simapp.GenerateRandomString(24),
				BeaconId: 2,
			},
			expectedError: sdkerrors.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet"),
		},
		{
			name: "not beacon owner",
			msg: &types.MsgRecordBeaconTimestamp{
				Owner:    testAddrs[1].String(),
				Hash:     simapp.GenerateRandomString(24),
				BeaconId: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon"),
		},
		{
			name: "successful",
			msg: &types.MsgRecordBeaconTimestamp{
				Owner:    testAddrs[0].String(),
				Hash:     simapp.GenerateRandomString(24),
				BeaconId: 1,
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

func TestValidMsgPurchaseBeaconStateStorage(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := beacon.NewHandler(app.BeaconKeeper)

	b := types.Beacon{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new beacon",
		Owner:   testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, b)
	require.Nil(t, err)

	msg := &types.MsgPurchaseBeaconStateStorage{
		BeaconId: 1,
		Number:   100,
		Owner:    testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgPurchaseBeaconStateStorage(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := beacon.NewHandler(app.BeaconKeeper)

	b := types.Beacon{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new beacon",
		Owner:   testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, b)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgPurchaseBeaconStateStorage
	}{
		{
			name:          "empty owner address",
			msg:           &types.MsgPurchaseBeaconStateStorage{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "cannot purchase zero",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner:  testAddrs[0].String(),
				Number: 0,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "cannot purchase zero"),
		},
		{
			name: "beacon not registered",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner:    testAddrs[0].String(),
				Number:   10,
				BeaconId: 2,
			},
			expectedError: sdkerrors.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet"),
		},
		{
			name: "not beacon owner",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner:    testAddrs[1].String(),
				Number:   10,
				BeaconId: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon"),
		},
		{
			name: "exceeds max storage",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner:    testAddrs[0].String(),
				Number:   simapp.TestMaxStorage,
				BeaconId: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrExceedsMaxStorage, fmt.Sprintf("%d will exceed max storage of %d", simapp.TestDefaultStorage+simapp.TestMaxStorage, simapp.TestMaxStorage)),
		},
		{
			name: "successful",
			msg: &types.MsgPurchaseBeaconStateStorage{
				Owner:    testAddrs[0].String(),
				Number:   10,
				BeaconId: 1,
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
