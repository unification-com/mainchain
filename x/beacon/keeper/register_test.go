package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/unification-com/mainchain/app"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/beacon/types"
)

// Tests for Highest BEACON ID

func TestSetGetHighestBeaconID(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	for i := uint64(1); i <= 1000; i++ {
		app.BeaconKeeper.SetHighestBeaconID(ctx, i)
		bID, err := app.BeaconKeeper.GetHighestBeaconID(ctx)
		require.NoError(t, err)
		require.True(t, bID == i)
	}
}

// Tests for Get/Set BEACONs

func TestSetGetBeacon(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(10)

	bID := uint64(1)
	for _, addr := range testAddrs {

		moniker := simapp.GenerateRandomString(12)
		name := simapp.GenerateRandomString(20)

		b := types.Beacon{}
		b.Owner = addr.String()
		b.BeaconId = bID
		b.LastTimestampId = 1
		b.Moniker = moniker
		b.Name = name

		err := app.BeaconKeeper.SetBeacon(ctx, b)
		require.NoError(t, err)

		// set the record limit
		err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, bID, types.DefaultStorageLimit)
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

		bSt, found := app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
		require.True(t, found)
		require.True(t, bSt.InStateLimit == types.DefaultStorageLimit)

		bID = bID + 1
	}
}

// Tests for Registering a new BEACON

func TestRegisterBeacon(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(10)

	i, _ := app.BeaconKeeper.GetHighestBeaconID(ctx)

	for _, addr := range testAddrs {
		name := simapp.GenerateRandomString(128)
		moniker := simapp.GenerateRandomString(64)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.BeaconId = i
		expectedB.LastTimestampId = 0
		expectedB.Moniker = moniker
		expectedB.Name = name
		expectedB.RegTime = uint64(ctx.BlockTime().Unix())

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

		bSt, found := app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
		require.True(t, found)
		require.True(t, bSt.InStateLimit == types.DefaultStorageLimit)

		i = i + 1
	}
}

func TestHighestBeaconIdAfterRegister(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	for i := uint64(1); i < 1000; i++ {
		name := simapp.GenerateRandomString(20)
		moniker := simapp.GenerateRandomString(12)
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
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	for i := uint64(1); i < 1000; i++ {
		name := simapp.GenerateRandomString(20)
		moniker := simapp.GenerateRandomString(12)
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

		bSt, found := app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
		require.True(t, found)
		require.True(t, bSt.InStateLimit == types.DefaultStorageLimit)
	}
}

func TestGetBeaconFilter(t *testing.T) {
	app := simapp.Setup(t, false)
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
