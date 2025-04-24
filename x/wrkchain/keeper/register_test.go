package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// Tests for Highest WRKChain ID

func TestSetGetHighestWRKChainID(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	for i := uint64(1); i <= 1000; i++ {
		app.WrkchainKeeper.SetHighestWrkChainID(ctx, i)
		wcID, err := app.WrkchainKeeper.GetHighestWrkChainID(ctx)
		require.NoError(t, err)
		require.True(t, wcID == i)
	}
}

//func TestSetGetHighestBeaconIDNotSet(t *testing.T) {
//	app := simapp.Setup(t, true)
//	ctx := app.BaseApp.NewContext(true)
//
//	_, err := app.WrkchainKeeper.GetHighestWrkChainID(ctx)
//	require.Error(t, err)
//}

// Tests for Get/Set WRKChains

func TestSetGetWrkChain(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)

	wcID := uint64(1)
	for _, addr := range testAddrs {

		wc := types.WrkChain{
			WrkchainId: wcID,
			Moniker:    simapphelpers.GenerateRandomString(12),
			Name:       simapphelpers.GenerateRandomString(20),
			Genesis:    simapphelpers.GenerateRandomString(32),
			BaseType:   "tendemrint",
			RegTime:    uint64(time.Now().Unix()),
			Owner:      addr.String(),
		}

		err := app.WrkchainKeeper.SetWrkChain(ctx, wc)
		require.NoError(t, err)

		// set the record limit
		err = app.WrkchainKeeper.SetWrkChainStorageLimit(ctx, wcID, types.DefaultStorageLimit)
		require.NoError(t, err)

		isRegistered := app.WrkchainKeeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)

		wcDb, found := app.WrkchainKeeper.GetWrkChain(ctx, wcID)
		require.True(t, found)
		require.True(t, WRKChainEqual(wcDb, wc))

		wcDbOwner := app.WrkchainKeeper.GetWrkChainOwner(ctx, wcID)
		require.True(t, wcDbOwner.String() == addr.String())

		wcSt, found := app.WrkchainKeeper.GetWrkChainStorageLimit(ctx, wcID)
		require.True(t, found)
		require.True(t, wcSt.InStateLimit == types.DefaultStorageLimit)

		wcID = wcID + 1
	}
}

// Tests for Registering a new WRKChain

func TestRegisterWrkChain(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)

	i, _ := app.WrkchainKeeper.GetHighestWrkChainID(ctx)

	for _, addr := range testAddrs {
		name := simapphelpers.GenerateRandomString(128)
		moniker := simapphelpers.GenerateRandomString(64)
		genesisHash := simapphelpers.GenerateRandomString(66)

		expectedWc := types.WrkChain{}
		expectedWc.Owner = addr.String()
		expectedWc.WrkchainId = i
		expectedWc.Lastblock = 0
		expectedWc.RegTime = uint64(time.Now().Unix())
		expectedWc.Moniker = moniker
		expectedWc.Name = name
		expectedWc.Genesis = genesisHash
		expectedWc.BaseType = "geth"

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)
		require.True(t, wcID == expectedWc.WrkchainId)

		isRegistered := app.WrkchainKeeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)

		wcDb, found := app.WrkchainKeeper.GetWrkChain(ctx, wcID)
		require.True(t, found)

		// hackery for reg time, otherwise following test fails
		expectedWc.RegTime = wcDb.RegTime
		require.True(t, WRKChainEqual(wcDb, expectedWc))

		wcDbOwner := app.WrkchainKeeper.GetWrkChainOwner(ctx, wcID)
		require.True(t, wcDbOwner.String() == addr.String())

		wcSt, found := app.WrkchainKeeper.GetWrkChainStorageLimit(ctx, wcID)
		require.True(t, found)
		require.True(t, wcSt.InStateLimit == simapphelpers.SimTestDefaultStorageLimit)

		i = i + 1
	}
}

func TestHighestWrkChainIdAfterRegister(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)

	for i := uint64(1); i < 1000; i++ {
		name := simapphelpers.GenerateRandomString(20)
		moniker := simapphelpers.GenerateRandomString(12)
		genesisHash := simapphelpers.GenerateRandomString(32)
		owner := testAddrs[1]

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		require.NoError(t, err)

		nextID, _ := app.WrkchainKeeper.GetHighestWrkChainID(ctx)
		expectedNextID := wcID + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestWrkChainIsRegisteredAfterRegister(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)

	for i := uint64(1); i < 1000; i++ {
		name := simapphelpers.GenerateRandomString(20)
		moniker := simapphelpers.GenerateRandomString(12)
		genesisHash := simapphelpers.GenerateRandomString(32)
		owner := testAddrs[1]

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		require.NoError(t, err)

		isRegistered := app.WrkchainKeeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)
	}
}

func TestGetWrkChainFilter(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)
	numToReg := 100
	lastMoniker := ""

	for i := 0; i < numToReg; i++ {
		name := simapphelpers.GenerateRandomString(20)
		moniker := simapphelpers.GenerateRandomString(12)
		genesisHash := simapphelpers.GenerateRandomString(32)
		owner := testAddrs[1]

		_, _ = app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		lastMoniker = moniker
	}

	params := types.QueryWrkChainsFilteredRequest{
		Owner: testAddrs[1].String(),
	}

	results := app.WrkchainKeeper.GetWrkChainsFiltered(ctx, params)
	require.True(t, len(results) == numToReg)

	params = types.QueryWrkChainsFilteredRequest{
		Moniker: lastMoniker,
	}

	results = app.WrkchainKeeper.GetWrkChainsFiltered(ctx, params)
	require.True(t, len(results) == 1)
}
