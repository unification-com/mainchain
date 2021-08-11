package wrkchain_test

import (
	"errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"strings"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain"
	"github.com/unification-com/mainchain/x/wrkchain/keeper"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := wrkchain.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognised wrkchain message type"))
}

func TestValidMsgRegisterWrkChain(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	msg := &types.MsgRegisterWrkChain{
		Moniker: "moniker",
		Name: "name",
		GenesisHash: "lhvviuvi",
		BaseType: "tendermint",
		Owner: testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgRegisterWrkChain(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(2)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	existsMoniker := test_helpers.GenerateRandomString(24)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, existsMoniker, "this exists", "boiob", "tendermint", testAddrs[1])
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRegisterWrkChain
	}{
		{
			name:          "empty owner address",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name: "name",
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name:          "invalid owner address",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name: "name",
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name:          "name too big",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name: test_helpers.GenerateRandomString(129),
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit"),
		},
		{
			name:          "moniker too big",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: test_helpers.GenerateRandomString(65),
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit"),
		},
		{
			name:          "zero length moniker",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: "",
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrap(types.ErrMissingData, "unable to register wrkchain - must have a moniker"),
		},
		{
			name:          "wrkchain exists with moniker",
			msg:           &types.MsgRegisterWrkChain{
				Moniker: existsMoniker,
				Name: "name",
				Owner: testAddrs[0].String(),
			},
			expectedError: sdkerrors.Wrapf(
				types.ErrWrkChainAlreadyRegistered,
				"wrkchain already registered with moniker '%s' - id: %d, owner: %s",
				existsMoniker, 1, testAddrs[1].String(),
			),
		},
		{
			name:          "successful",
			msg:           &types.MsgRegisterWrkChain{
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

func TestValidMsgRecordWrkChainBlock(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(
		ctx,
		test_helpers.GenerateRandomString(24),
		test_helpers.GenerateRandomString(24),
		"boiob",
		"tendermint",
		testAddrs[0],
		)
	require.Nil(t, err)

	msg := &types.MsgRecordWrkChainBlock{
		WrkchainId: 1,
		BlockHash: test_helpers.GenerateRandomString(64),
		Owner: testAddrs[0].String(),
		Height: 1,
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgRecordWrkChainBlock(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(2)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(
		ctx,
		test_helpers.GenerateRandomString(24),
		test_helpers.GenerateRandomString(24),
		"boiob",
		"tendermint",
		testAddrs[0],
	)
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRecordWrkChainBlock
	}{
		{
			name: "empty owner address",
			msg: &types.MsgRecordWrkChainBlock{
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name:          "invalid owner address",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name:          "zero height",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(66),
				Height: 0,
			},
			expectedError: sdkerrors.Wrap(types.ErrInvalidData, "height must be > 0"),
		},
		{
			name:          "blockhash too large",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(67),
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "block hash too big. 66 character limit"),
		},
		{
			name:          "parenthash too large",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(64),
				ParentHash: test_helpers.GenerateRandomString(67),
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "parent hash too big. 66 character limit"),
		},
		{
			name:          "hash1 too large",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(64),
				ParentHash: test_helpers.GenerateRandomString(66),
				Hash1: test_helpers.GenerateRandomString(67),
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "hash1 too big. 66 character limit"),
		},
		{
			name:          "hash2 too large",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(64),
				ParentHash: test_helpers.GenerateRandomString(66),
				Hash1: test_helpers.GenerateRandomString(66),
				Hash2: test_helpers.GenerateRandomString(67),
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "hash2 too big. 66 character limit"),
		},
		{
			name:          "hash3 too large",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(64),
				ParentHash: test_helpers.GenerateRandomString(66),
				Hash1: test_helpers.GenerateRandomString(66),
				Hash2: test_helpers.GenerateRandomString(66),
				Hash3: test_helpers.GenerateRandomString(67),
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrContentTooLarge, "hash3 too big. 66 character limit"),
		},
		{
			name:          "wrkchain not registered",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(24),
				WrkchainId: 2,
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrWrkChainDoesNotExist, "wrkchain has not been registered yet"),
		},
		{
			name:          "not wrkchain owner",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[1].String(),
				BlockHash: test_helpers.GenerateRandomString(24),
				WrkchainId: 1,
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrNotWrkChainOwner, "you are not the owner of this wrkchain"),
		},
		{
			name:          "successful",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(24),
				WrkchainId: 1,
				Height: 1,
			},
			expectedError: nil,
		},
		{
			name:          "height already recorded",
			msg:           &types.MsgRecordWrkChainBlock{
				Owner: testAddrs[0].String(),
				BlockHash: test_helpers.GenerateRandomString(24),
				WrkchainId: 1,
				Height: 1,
			},
			expectedError: sdkerrors.Wrap(types.ErrWrkChainBlockAlreadyRecorded, "wrkchain block hashes have already been recorded for this height"),
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
