package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

// Tests for Highest BEACON ID

func TestSetGetHighestBeaconID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i <= 1000; i++ {
		keeper.SetHighestBeaconID(ctx, i)
		bID, err := keeper.GetHighestBeaconID(ctx)
		require.NoError(t, err)
		require.True(t, bID == i)
	}
}

func TestSetGetHighestBeaconIDNotSet(t *testing.T) {

	keyBeacon := sdk.NewKVStoreKey(types.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)


	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyBeacon, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	require.Nil(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "und-unit-test-chain"}, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)

	cdc := makeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	keeper := NewKeeper(
		keyBeacon, pk.Subspace(types.DefaultParamspace), cdc,
	)

	expectedErr := sdkerrors.Wrapf(types.ErrInvalidGenesis, "initial beacon ID hasn't been set")
	bId, err := keeper.GetHighestBeaconID(ctx)

	require.Equal(t, expectedErr.Error(), err.Error())
	require.Equal(t, uint64(0x0), bId)
}

// Tests for Get/Set BEACONs

func TestSetGetBeacon(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 1000)

	bID := uint64(1)
	for _, addr := range TestAddrs {

		b := types.NewBeacon()
		b.Owner = addr
		b.BeaconID = bID
		b.LastTimestampID = 1
		b.Moniker = GenerateRandomString(12)
		b.Name = GenerateRandomString(20)

		err := keeper.SetBeacon(ctx, b)
		require.NoError(t, err)

		isRegistered := keeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)

		bDb := keeper.GetBeacon(ctx, bID)
		require.True(t, BeaconEqual(bDb, b))

		bDbOwner := keeper.GetBeaconOwner(ctx, bID)
		require.True(t, bDbOwner.String() == addr.String())

		bID = bID + 1
	}
}

func TestSetLastTimestampID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	lastTimestampID := uint64(0)
	bID := uint64(1)

	b := types.NewBeacon()
	b.Owner = TestAddrs[0]
	b.BeaconID = bID
	b.LastTimestampID = lastTimestampID
	b.Moniker = GenerateRandomString(12)
	b.Name = GenerateRandomString(20)

	err := keeper.SetBeacon(ctx, b)
	require.NoError(t, err)

	for i := uint64(1); i <= 1000; i++ {
		err := keeper.SetLastTimestampID(ctx, bID, i)
		require.NoError(t, err)

		bDb := keeper.GetBeacon(ctx, bID)
		require.True(t, bDb.LastTimestampID == i)
		lastTimestampID = i
	}

	// check can't set last block to < current last block
	oldTsID := lastTimestampID - 1
	err = keeper.SetLastTimestampID(ctx, bID, oldTsID)
	require.NoError(t, err)
	wcDb := keeper.GetBeacon(ctx, bID)
	require.True(t, wcDb.LastTimestampID == lastTimestampID)

}

func TestEmptySetBeaconValuesReturnError(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	b0 := types.NewBeacon()

	b1 := b0
	b1.Owner = TestAddrs[0]

	b2 := b1
	b2.BeaconID = uint64(1)

	b3 := b2
	b3.Moniker = "moniker"

	testCases := []struct {
		b           types.Beacon
		expectedErr error
	}{
		{b0, sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have owner")},
		{b1, sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - id must be positive non-zero")},
		{b2, sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have a moniker")},
		{b3, nil},
	}

	for _, tc := range testCases {
		err := keeper.SetBeacon(ctx, tc.b)
		if tc.expectedErr == nil {
			require.NoError(t, err)
		} else {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
		}
	}
}

// Tests for Registering a new BEACON

func TestRegisterBeacon(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 1000)

	i, _ := keeper.GetHighestBeaconID(ctx)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(128)
		moniker := GenerateRandomString(64)

		expectedB := types.NewBeacon()
		expectedB.Owner = addr
		expectedB.BeaconID = i
		expectedB.LastTimestampID = 0
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, addr)
		require.NoError(t, err)
		require.True(t, bID == expectedB.BeaconID)

		isRegistered := keeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)

		bDb := keeper.GetBeacon(ctx, bID)

		require.True(t, BeaconEqual(bDb, expectedB))

		bDbOwner := keeper.GetBeaconOwner(ctx, bID)
		require.True(t, bDbOwner.String() == addr.String())

		i = i + 1
	}
}

func TestFailRegisterNewBeacon(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	testCases := []struct {
		moniker     string
		name        string
		owner       sdk.AccAddress
		expectedErr error
		expectedBID uint64
	}{
		{"testmoniker", "", TestAddrs[0], nil, 1},
		{"", "", TestAddrs[0], sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have a moniker"), 0},
		{"", "", sdk.AccAddress{}, sdkerrors.Wrap(types.ErrMissingData, "unable to set beacon - must have owner"), 0},
	}

	for _, tc := range testCases {
		wcID, err := keeper.RegisterBeacon(ctx, tc.moniker, tc.name, tc.owner)
		if tc.expectedErr == nil {
			require.NoError(t, err)
		} else {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
		}
		require.True(t, wcID == tc.expectedBID)
	}
}

func TestHighestBeaconIdAfterRegister(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i < 1000; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		owner := TestAddrs[1]

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, owner)
		require.NoError(t, err)

		nextID, _ := keeper.GetHighestBeaconID(ctx)
		expectedNextID := bID + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestBeaconIsRegisteredAfterRegister(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i < 1000; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		owner := TestAddrs[1]

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, owner)
		require.NoError(t, err)

		isRegistered := keeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)
	}
}

func TestGetBeaconFilter(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToReg := 100
	lastMoniker := ""

	for i := 0; i < numToReg; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		owner := TestAddrs[1]

		_, _ = keeper.RegisterBeacon(ctx, moniker, name, owner)
		lastMoniker = moniker
	}

	params := types.NewQueryBeaconParams(1, 1000, "", TestAddrs[1])
	results := keeper.GetBeaconsFiltered(ctx, params)
	require.True(t, len(results) == numToReg)

	params = types.NewQueryBeaconParams(1, 1000, lastMoniker, sdk.AccAddress{})
	results = keeper.GetBeaconsFiltered(ctx, params)
	require.True(t, len(results) == 1)
}
