package wrkchain_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	errorsmod "cosmossdk.io/errors"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/wrkchain/types"

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
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	msg := &types.MsgRegisterWrkChain{
		Moniker:     "moniker",
		Name:        "name",
		GenesisHash: "lhvviuvi",
		BaseType:    "tendermint",
		Owner:       testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)
}

func TestInvalidMsgRegisterWrkChain(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	existsMoniker := simapp.GenerateRandomString(24)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, existsMoniker, "this exists", "boiob", "tendermint", testAddrs[1])
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgRegisterWrkChain
	}{
		{
			name: "empty owner address",
			msg: &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name:    "name",
			},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name:    "name",
				Owner:   "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "name too big",
			msg: &types.MsgRegisterWrkChain{
				Moniker: "moniker",
				Name:    simapp.GenerateRandomString(129),
				Owner:   testAddrs[0].String(),
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "name too big. 128 character limit"),
		},
		{
			name: "moniker too big",
			msg: &types.MsgRegisterWrkChain{
				Moniker: simapp.GenerateRandomString(65),
				Name:    "name",
				Owner:   testAddrs[0].String(),
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "moniker too big. 64 character limit"),
		},
		{
			name: "zero length moniker",
			msg: &types.MsgRegisterWrkChain{
				Moniker: "",
				Name:    "name",
				Owner:   testAddrs[0].String(),
			},
			expectedError: errorsmod.Wrap(types.ErrMissingData, "unable to register wrkchain - must have a moniker"),
		},
		{
			name: "successful",
			msg: &types.MsgRegisterWrkChain{
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

func TestValidMsgRecordWrkChainBlock(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(
		ctx,
		simapp.GenerateRandomString(24),
		simapp.GenerateRandomString(24),
		"boiob",
		"tendermint",
		testAddrs[0],
	)
	require.Nil(t, err)

	msg := &types.MsgRecordWrkChainBlock{
		WrkchainId: 1,
		BlockHash:  simapp.GenerateRandomString(64),
		Owner:      testAddrs[0].String(),
		Height:     1,
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgRecordWrkChainBlock(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(
		ctx,
		simapp.GenerateRandomString(24),
		simapp.GenerateRandomString(24),
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
			name:          "empty owner address",
			msg:           &types.MsgRecordWrkChainBlock{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgRecordWrkChainBlock{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "zero height",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:     testAddrs[0].String(),
				BlockHash: simapp.GenerateRandomString(66),
				Height:    0,
			},
			expectedError: errorsmod.Wrap(types.ErrInvalidData, "height must be > 0"),
		},
		{
			name: "blockhash too large",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:     testAddrs[0].String(),
				BlockHash: simapp.GenerateRandomString(67),
				Height:    1,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "block hash too big. 66 character limit"),
		},
		{
			name: "parenthash too large",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(64),
				ParentHash: simapp.GenerateRandomString(67),
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "parent hash too big. 66 character limit"),
		},
		{
			name: "hash1 too large",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(64),
				ParentHash: simapp.GenerateRandomString(66),
				Hash1:      simapp.GenerateRandomString(67),
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "hash1 too big. 66 character limit"),
		},
		{
			name: "hash2 too large",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(64),
				ParentHash: simapp.GenerateRandomString(66),
				Hash1:      simapp.GenerateRandomString(66),
				Hash2:      simapp.GenerateRandomString(67),
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "hash2 too big. 66 character limit"),
		},
		{
			name: "hash3 too large",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(64),
				ParentHash: simapp.GenerateRandomString(66),
				Hash1:      simapp.GenerateRandomString(66),
				Hash2:      simapp.GenerateRandomString(66),
				Hash3:      simapp.GenerateRandomString(67),
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "hash3 too big. 66 character limit"),
		},
		{
			name: "wrkchain not registered",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(24),
				WrkchainId: 2,
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrWrkChainDoesNotExist, "wrkchain has not been registered yet"),
		},
		{
			name: "not wrkchain owner",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[1].String(),
				BlockHash:  simapp.GenerateRandomString(24),
				WrkchainId: 1,
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrNotWrkChainOwner, "you are not the owner of this wrkchain"),
		},
		{
			name: "successful",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(24),
				WrkchainId: 1,
				Height:     1,
			},
			expectedError: nil,
		},
		{
			name: "height too low",
			msg: &types.MsgRecordWrkChainBlock{
				Owner:      testAddrs[0].String(),
				BlockHash:  simapp.GenerateRandomString(24),
				WrkchainId: 1,
				Height:     1,
			},
			expectedError: errorsmod.Wrap(types.ErrNewHeightMustBeHigher, "wrkchain block hashes height must be > last height recorded"),
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

func TestValidMsgPurchaseWrkChainStateStorage(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	wc := types.WrkChain{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new wrkchain",
		Owner:   testAddrs[0].String(),
		Genesis: "genesis",
	}

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, wc.Moniker, wc.Name, wc.Genesis, "test", testAddrs[0])
	require.Nil(t, err)

	msg := &types.MsgPurchaseWrkChainStateStorage{
		WrkchainId: 1,
		Number:     100,
		Owner:      testAddrs[0].String(),
	}

	res, err := h(ctx, msg)

	require.NotNil(t, res)
	require.Nil(t, err)

}

func TestInvalidMsgPurchaseWrkChainStateStorage(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(2)

	h := wrkchain.NewHandler(app.WrkchainKeeper)

	wc := types.WrkChain{
		Moniker: simapp.GenerateRandomString(24),
		Name:    "new wrkchain",
		Owner:   testAddrs[0].String(),
		Genesis: "genesis",
	}

	_, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, wc.Moniker, wc.Name, wc.Genesis, "test", testAddrs[0])
	require.Nil(t, err)

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgPurchaseWrkChainStateStorage
	}{
		{
			name:          "empty owner address",
			msg:           &types.MsgPurchaseWrkChainStateStorage{},
			expectedError: errors.New("empty address string is not allowed"),
		},
		{
			name: "invalid owner address",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner: "rubbish",
			},
			expectedError: errors.New("decoding bech32 failed: invalid bech32 string length 7"),
		},
		{
			name: "cannot purchase zero",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner:  testAddrs[0].String(),
				Number: 0,
			},
			expectedError: errorsmod.Wrap(types.ErrContentTooLarge, "cannot purchase zero"),
		},
		{
			name: "wrkchain not registered",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      testAddrs[0].String(),
				Number:     10,
				WrkchainId: 2,
			},
			expectedError: errorsmod.Wrap(types.ErrWrkChainDoesNotExist, "wrkchain has not been registered yet"),
		},
		{
			name: "not wrkchain owner",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      testAddrs[1].String(),
				Number:     10,
				WrkchainId: 1,
			},
			expectedError: errorsmod.Wrap(types.ErrNotWrkChainOwner, "you are not the owner of this wrkchain"),
		},
		{
			name: "exceeds max storage",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      testAddrs[0].String(),
				Number:     simapp.TestMaxStorage,
				WrkchainId: 1,
			},
			expectedError: errorsmod.Wrap(types.ErrExceedsMaxStorage, fmt.Sprintf("%d will exceed max storage of %d", simapp.TestDefaultStorage+simapp.TestMaxStorage, simapp.TestMaxStorage)),
		},
		{
			name: "successful",
			msg: &types.MsgPurchaseWrkChainStateStorage{
				Owner:      testAddrs[0].String(),
				Number:     10,
				WrkchainId: 1,
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
