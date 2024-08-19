package keeper_test

import (
	simapp "github.com/unification-com/mainchain/app"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/beacon/types"
)

func TestSetGetBeaconTimestamp(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(10)

	numToRecord := uint64(1000)

	for _, addr := range testAddrs {
		name := simapp.GenerateRandomString(20)
		moniker := simapp.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			beaconTimestamp := types.BeaconTimestamp{}
			beaconTimestamp.TimestampId = tsID
			beaconTimestamp.Hash = simapp.GenerateRandomString(32)
			beaconTimestamp.SubmitTime = uint64(time.Now().Unix())

			err := app.BeaconKeeper.SetBeaconTimestamp(ctx, bID, beaconTimestamp)
			require.NoError(t, err)

			btsDb, found := app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsID)
			require.True(t, found)
			require.True(t, BeaconTimestampEqual(btsDb, beaconTimestamp))
		}
	}
}

func TestGetBeaconTimestamp(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	for _, addr := range testAddrs {
		name := simapp.GenerateRandomString(20)
		moniker := simapp.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)

		require.NoError(t, err)

		var testTimestamps []types.BeaconTimestamp

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			subTime := uint64(time.Now().Unix())
			hash := simapp.GenerateRandomString(32)

			timestamp := types.BeaconTimestamp{}
			timestamp.TimestampId = tsID
			timestamp.Hash = hash
			timestamp.SubmitTime = subTime

			testTimestamps = append(testTimestamps, timestamp)

			err := app.BeaconKeeper.SetBeaconTimestamp(ctx, bID, timestamp)
			require.NoError(t, err)
		}

		allTimestamps := app.BeaconKeeper.GetAllBeaconTimestamps(ctx, bID)
		require.True(t, len(allTimestamps) == int(numToRecord) && len(allTimestamps) == len(testTimestamps))

		for i := 0; i < int(numToRecord); i++ {
			require.True(t, BeaconTimestampEqual(allTimestamps[i], testTimestamps[i]))
		}
	}
}

func TestIsAuthorisedToRecord(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(10)

	unauthorisedAddrs := simapp.GenerateRandomTestAccounts(1)

	for _, addr := range testAddrs {
		name := simapp.GenerateRandomString(20)
		moniker := simapp.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		isAuthorised := app.BeaconKeeper.IsAuthorisedToRecord(ctx, bID, addr)
		require.True(t, isAuthorised)

		isAuthorised = app.BeaconKeeper.IsAuthorisedToRecord(ctx, bID, unauthorisedAddrs[0])
		require.False(t, isAuthorised)
	}
}

func TestRecordBeaconTimestamps(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	numToRecord := uint64(1000)
	recordLimit := uint64(200)

	name := simapp.GenerateRandomString(128)
	moniker := simapp.GenerateRandomString(64)

	expectedB := types.Beacon{}
	expectedB.Owner = testAddrs[0].String()
	expectedB.Moniker = moniker
	expectedB.Name = name

	bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
	require.NoError(t, err)

	// set the record limit
	err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, bID, recordLimit)
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numToRecord; tsID++ {
		subTime := uint64(time.Now().Unix())
		hash := simapp.GenerateRandomString(32)

		expectedTs := types.BeaconTimestamp{}
		expectedTs.TimestampId = tsID
		expectedTs.Hash = hash
		expectedTs.SubmitTime = subTime

		retTsID, deletedTsId, err := app.BeaconKeeper.RecordNewBeaconTimestamp(ctx, bID, hash, subTime)
		require.NoError(t, err)
		require.True(t, retTsID == expectedTs.TimestampId)

		timestampDb, found := app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsID)
		require.True(t, found)
		require.True(t, BeaconTimestampEqual(timestampDb, expectedTs))

		beacon, found := app.BeaconKeeper.GetBeacon(ctx, bID)
		require.True(t, found)
		require.Equal(t, retTsID, beacon.LastTimestampId, "not equal: exp = %d, act = %d", retTsID, beacon.LastTimestampId)

		if deletedTsId > 0 {
			_, found := app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, deletedTsId)
			require.False(t, found)
		}
	}

	beacon, found := app.BeaconKeeper.GetBeacon(ctx, bID)
	require.True(t, found)
	require.True(t, beacon.NumInState == recordLimit)
	require.True(t, beacon.FirstIdInState == numToRecord-recordLimit+1)

	// should still be in state
	tsCount := uint64(0)
	for tsId := beacon.FirstIdInState; tsId <= beacon.LastTimestampId; tsId++ {
		_, found = app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsId)
		require.True(t, found)
		tsCount++
	}
	require.True(t, tsCount == recordLimit)

	// should no longer be in state
	for tsId := uint64(1); tsId <= recordLimit; tsId++ {
		_, found = app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsId)
		require.False(t, found)
	}
}

func TestIncreaseInStateStorage(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	recordLimitIncrease := uint64(200)

	name := simapp.GenerateRandomString(128)
	moniker := simapp.GenerateRandomString(64)

	expectedB := types.Beacon{}
	expectedB.Owner = testAddrs[0].String()
	expectedB.Moniker = moniker
	expectedB.Name = name

	bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
	require.NoError(t, err)

	beacon, found := app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
	require.True(t, found)
	require.True(t, beacon.InStateLimit == types.DefaultStorageLimit)

	err = app.BeaconKeeper.IncreaseInStateStorage(ctx, bID, recordLimitIncrease)
	require.NoError(t, err)

	beacon, found = app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
	require.True(t, found)
	require.True(t, beacon.InStateLimit == types.DefaultStorageLimit+recordLimitIncrease)
}

func TestIncreaseInStateStorageWithTimestampRecording(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := simapp.GenerateRandomTestAccounts(1)

	numToRecord := uint64(500)
	recordLimit := uint64(100)
	increaseAmount := uint64(50)

	name := simapp.GenerateRandomString(128)
	moniker := simapp.GenerateRandomString(64)

	expectedB := types.Beacon{}
	expectedB.Owner = testAddrs[0].String()
	expectedB.Moniker = moniker
	expectedB.Name = name

	bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
	require.NoError(t, err)

	// set the record limit
	err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, bID, recordLimit)
	require.NoError(t, err)

	// record initial timestamps
	for i := uint64(1); i <= numToRecord; i++ {
		hash := simapp.GenerateRandomString(32)
		subTime := uint64(time.Now().Unix())
		_, _, err := app.BeaconKeeper.RecordNewBeaconTimestamp(ctx, bID, hash, subTime)
		require.NoError(t, err)
	}

	// sanity check
	beacon, found := app.BeaconKeeper.GetBeacon(ctx, bID)
	require.True(t, found)
	require.True(t, beacon.NumInState == recordLimit)
	require.True(t, beacon.FirstIdInState == numToRecord-recordLimit+1)

	// increase storage capacity
	err = app.BeaconKeeper.IncreaseInStateStorage(ctx, bID, increaseAmount)
	require.NoError(t, err)
	beaconStorage, found := app.BeaconKeeper.GetBeaconStorageLimit(ctx, bID)
	require.True(t, found)
	require.True(t, beaconStorage.InStateLimit == recordLimit+increaseAmount)

	// record new timestamps
	for i := uint64(1); i <= numToRecord; i++ {
		hash := simapp.GenerateRandomString(32)
		subTime := uint64(time.Now().Unix())
		_, _, err := app.BeaconKeeper.RecordNewBeaconTimestamp(ctx, bID, hash, subTime)
		require.NoError(t, err)
	}

	// check final result
	beacon, found = app.BeaconKeeper.GetBeacon(ctx, bID)
	require.True(t, found)
	require.True(t, beacon.NumInState == recordLimit+increaseAmount)
	require.True(t, beacon.FirstIdInState == (numToRecord*2)-recordLimit-increaseAmount+1)
}
