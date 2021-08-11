package beacon_test

import (
	"errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
	"strings"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

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
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	h := beacon.NewHandler(app.BeaconKeeper)

	msg := &types.MsgRegisterBeacon{
		Moniker: "moniker",
		Name: "name",
		Owner: testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgRegisterBeacon(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(2)

	h := beacon.NewHandler(app.BeaconKeeper)

	existsMoniker := test_helpers.GenerateRandomString(24)
	beacon := types.Beacon{
		Moniker:         existsMoniker,
		Name:            "this exists",
		Owner:           testAddrs[1].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, beacon)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRegisterBeacon
	}{
		{
			name:          "empty owner address",
			msg:           &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name: "name",
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name:          "invalid owner address",
			msg:           &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name: "name",
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name:          "name too big",
			msg:           &types.MsgRegisterBeacon{
				Moniker: "moniker",
				Name: test_helpers.GenerateRandomString(129),
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit"),
		},
		{
			name:          "moniker too big",
			msg:           &types.MsgRegisterBeacon{
				Moniker: test_helpers.GenerateRandomString(65),
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit"),
		},
		{
			name:          "zero length moniker",
			msg:           &types.MsgRegisterBeacon{
				Moniker: "",
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrMissingData, "unable to register beacon - must have a moniker"),
		},
		{
			name:          "beacon exists with moniker",
			msg:           &types.MsgRegisterBeacon{
				Moniker: existsMoniker,
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrBeaconAlreadyRegistered,
				"beacon already registered with moniker '%s' - id: %d, owner: %s",
				existsMoniker, 1, testAddrs[1].String(),
				),
		},
		{
			name:          "successful",
			msg:           &types.MsgRegisterBeacon{
				Moniker: test_helpers.GenerateRandomString(24),
				Name: test_helpers.GenerateRandomString(24),
				Owner: testAddrs[0].String(),
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
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	h := beacon.NewHandler(app.BeaconKeeper)

	beacon := types.Beacon{
		Moniker:         test_helpers.GenerateRandomString(24),
		Name:            "new beacon",
		Owner:           testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, beacon)
	require.Nil(t, err)

	msg := &types.MsgRecordBeaconTimestamp{
		BeaconId: 1,
		Hash: test_helpers.GenerateRandomString(64),
		Owner: testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgRecordBeaconTimestamp(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(2)

	h := beacon.NewHandler(app.BeaconKeeper)

	beacon := types.Beacon{
		Moniker:         test_helpers.GenerateRandomString(24),
		Name:            "new beacon",
		Owner:           testAddrs[0].String(),
	}

	_, err := app.BeaconKeeper.RegisterNewBeacon(ctx, beacon)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRecordBeaconTimestamp
	}{
		{
			name: "empty owner address",
			msg: &types.MsgRecordBeaconTimestamp{
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name:          "invalid owner address",
			msg:           &types.MsgRecordBeaconTimestamp{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name:          "hash too large",
			msg:           &types.MsgRecordBeaconTimestamp{
				Owner: testAddrs[0].String(),
				Hash: test_helpers.GenerateRandomString(67),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "hash too big. 66 character limit"),
		},
		{
			name:          "beacon not registered",
			msg:           &types.MsgRecordBeaconTimestamp{
				Owner: testAddrs[0].String(),
				Hash: test_helpers.GenerateRandomString(24),
				BeaconId: 2,
			},
			expectedError: sdkerrors.Wrap(types.ErrBeaconDoesNotExist, "beacon has not been registered yet"),
		},
		{
			name:          "not beacon owner",
			msg:           &types.MsgRecordBeaconTimestamp{
				Owner: testAddrs[1].String(),
				Hash: test_helpers.GenerateRandomString(24),
				BeaconId: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrNotBeaconOwner, "you are not the owner of this beacon"),
		},
		{
			name:          "successful",
			msg:           &types.MsgRecordBeaconTimestamp{
				Owner: testAddrs[0].String(),
				Hash: test_helpers.GenerateRandomString(24),
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
