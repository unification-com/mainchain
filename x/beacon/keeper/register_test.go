package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
)

// Tests for Highest BEACON ID

func TestSetGetHighestBeaconID(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	for i := uint64(1); i <= 1000; i++ {
		app.BeaconKeeper.SetHighestBeaconID(ctx, i)
		bID, err := app.BeaconKeeper.GetHighestBeaconID(ctx)
		require.NoError(t, err)
		require.True(t, bID == i)
	}
}

func TestSetGetHighestBeaconIDNotSet(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	_, err := app.BeaconKeeper.GetHighestBeaconID(ctx)
	require.Error(t, err)
}

// Tests for Get/Set BEACONs

func TestSetGetBeacon(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	bID := uint64(1)
	for _, addr := range testAddrs {

		moniker := test_helpers.GenerateRandomString(12)
		name := test_helpers.GenerateRandomString(20)

		b := types.Beacon{}
		b.Owner = addr.String()
		b.BeaconId = bID
		b.LastTimestampId = 1
		b.Moniker = moniker
		b.Name = name

		err := app.BeaconKeeper.SetBeacon(ctx, b)
		require.NoError(t, err)

		isRegistered := app.BeaconKeeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)

		bDb, found := app.BeaconKeeper.GetBeacon(ctx, bID)
		require.True(t, found)

		require.True(t, bDb.Owner == addr.String())
		require.True(t, bDb.BeaconId == bID)
		require.True(t, bDb.LastTimestampId == 1)
		require.True(t, bDb.Moniker == moniker)
		require.True(t, bDb.Name == name)

		bID = bID + 1
	}
}

func TestSetLastTimestampID(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	lastTimestampID := uint64(0)
	bID := uint64(1)

	b := types.Beacon{}
	b.Owner = testAddrs[0].String()
	b.BeaconId = bID
	b.LastTimestampId = lastTimestampID
	b.Moniker = test_helpers.GenerateRandomString(12)
	b.Name = test_helpers.GenerateRandomString(20)

	err := app.BeaconKeeper.SetBeacon(ctx, b)
	require.NoError(t, err)

	for i := uint64(1); i <= 1000; i++ {
		err := app.BeaconKeeper.SetLastTimestampID(ctx, bID, i)
		require.NoError(t, err)

		bDb, found := app.BeaconKeeper.GetBeacon(ctx, bID)
		require.True(t, found)
		require.True(t, bDb.LastTimestampId == i)
		lastTimestampID = i
	}

	// check can't set last block to < current last block
	oldTsID := lastTimestampID - 1
	err = app.BeaconKeeper.SetLastTimestampID(ctx, bID, oldTsID)
	require.NoError(t, err)
	wcDb, found := app.BeaconKeeper.GetBeacon(ctx, bID)
	require.True(t, found)
	require.True(t, wcDb.LastTimestampId == lastTimestampID)

}

// Tests for Registering a new BEACON

func TestRegisterBeacon(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	i, _ := app.BeaconKeeper.GetHighestBeaconID(ctx)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(128)
		moniker := test_helpers.GenerateRandomString(64)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.BeaconId = i
		expectedB.LastTimestampId = 0
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)
		require.True(t, bID == expectedB.BeaconId)

		isRegistered := app.BeaconKeeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)

		bDb, found := app.BeaconKeeper.GetBeacon(ctx, bID)
		require.True(t, found)

		require.True(t, BeaconEqual(bDb, expectedB))

		bDbOwner := app.BeaconKeeper.GetBeaconOwner(ctx, bID)
		require.True(t, bDbOwner.String() == addr.String())

		i = i + 1
	}
}

func TestFailSetLastTimestampId(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	err := app.BeaconKeeper.SetLastTimestampID(ctx, 1, 1)
	require.Error(t, err)
}

func TestHighestBeaconIdAfterRegister(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	for i := uint64(1); i < 1000; i++ {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)
		owner := testAddrs[0].String()
		expectedB := types.Beacon{}
		expectedB.Owner = owner
		expectedB.BeaconId = i
		expectedB.LastTimestampId = 0
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		nextID, _ := app.BeaconKeeper.GetHighestBeaconID(ctx)
		expectedNextID := bID + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestBeaconIsRegisteredAfterRegister(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	for i := uint64(1); i < 1000; i++ {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)
		owner := testAddrs[0].String()

		expectedB := types.Beacon{}
		expectedB.Owner = owner
		expectedB.BeaconId = i
		expectedB.LastTimestampId = 0
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		isRegistered := app.BeaconKeeper.IsBeaconRegistered(ctx, bID)
		require.True(t, isRegistered)
	}
}

func TestGetBeaconFilter(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	numToReg := 100
	lastMoniker := ""

	for i := 0; i < numToReg; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		owner := TestAddrs[1]
		expectedB := types.Beacon{}
		expectedB.Owner = owner.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		_, _ = app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		lastMoniker = moniker
	}

	params := types.QueryBeaconsFilteredRequest{
		Owner: TestAddrs[1].String(),
	}

	results := app.BeaconKeeper.GetBeaconsFiltered(ctx, params)
	require.True(t, len(results) == numToReg)

	params = types.QueryBeaconsFilteredRequest{
		Moniker: lastMoniker,
	}

	results = app.BeaconKeeper.GetBeaconsFiltered(ctx, params)
	require.True(t, len(results) == 1)
}
