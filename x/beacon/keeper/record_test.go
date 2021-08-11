package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
)

func TestSetGetBeaconTimestamp(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	numToRecord := uint64(1000)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			beaconTimestamp := types.BeaconTimestamp{}
			beaconTimestamp.BeaconId = bID
			beaconTimestamp.Owner = addr.String()
			beaconTimestamp.TimestampId = tsID
			beaconTimestamp.Hash = test_helpers.GenerateRandomString(32)
			beaconTimestamp.SubmitTime = uint64(time.Now().Unix())

			err := app.BeaconKeeper.SetBeaconTimestamp(ctx, beaconTimestamp)
			require.NoError(t, err)

			btsDb, found := app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsID)
			require.True(t, found)
			require.True(t, BeaconTimestampEqual(btsDb, beaconTimestamp))
		}
	}
}

func TestIsBeaconTimestampRecorded(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
		require.NoError(t, err)

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			hash := test_helpers.GenerateRandomString(32)
			subTime := uint64(time.Now().Unix())
			timestamp := types.BeaconTimestamp{}
			timestamp.BeaconId = bID
			timestamp.Owner = addr.String()
			timestamp.TimestampId = tsID
			timestamp.Hash = hash
			timestamp.SubmitTime = subTime

			err := app.BeaconKeeper.SetBeaconTimestamp(ctx, timestamp)
			require.NoError(t, err)

			isRecorded := app.BeaconKeeper.IsBeaconTimestampRecordedByID(ctx, bID, tsID)
			require.True(t, isRecorded)

			isRecorded1 := app.BeaconKeeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, hash, 0)
			require.True(t, isRecorded1)

			isRecorded2 := app.BeaconKeeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, "", subTime)
			require.True(t, isRecorded2)

			isRecorded3 := app.BeaconKeeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, hash, subTime)
			require.True(t, isRecorded3)
		}
	}
}

func TestGetWrkChainBlockHashes(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)

		expectedB := types.Beacon{}
		expectedB.Owner = addr.String()
		expectedB.Moniker = moniker
		expectedB.Name = name

		bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)

		require.NoError(t, err)

		var testTimestamps []types.BeaconTimestamp

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			subTime := uint64(time.Now().Unix())
			hash := test_helpers.GenerateRandomString(32)

			timestamp := types.BeaconTimestamp{}
			timestamp.BeaconId = bID
			timestamp.Owner = addr.String()
			timestamp.TimestampId = tsID
			timestamp.Hash = hash
			timestamp.SubmitTime = subTime

			testTimestamps = append(testTimestamps, timestamp)

			err := app.BeaconKeeper.SetBeaconTimestamp(ctx, timestamp)
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
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	unauthorisedAddrs := test_helpers.GenerateRandomTestAccounts(1)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(20)
		moniker := test_helpers.GenerateRandomString(12)

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
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	numToRecord := uint64(100)

	name := test_helpers.GenerateRandomString(128)
	moniker := test_helpers.GenerateRandomString(64)

	expectedB := types.Beacon{}
	expectedB.Owner = testAddrs[0].String()
	expectedB.Moniker = moniker
	expectedB.Name = name

	bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numToRecord; tsID++ {
		subTime := uint64(time.Now().Unix())
		hash := test_helpers.GenerateRandomString(32)

		expectedTs := types.BeaconTimestamp{}
		expectedTs.BeaconId = bID
		expectedTs.Owner = testAddrs[0].String()
		expectedTs.TimestampId = tsID
		expectedTs.Hash = hash
		expectedTs.SubmitTime = subTime

		retTsID, err := app.BeaconKeeper.RecordNewBeaconTimestamp(ctx, bID, hash, subTime, testAddrs[0].String())
		require.NoError(t, err)
		require.True(t, retTsID == expectedTs.TimestampId)

		timestampDb, found := app.BeaconKeeper.GetBeaconTimestampByID(ctx, bID, tsID)
		require.True(t, found)
		require.True(t, BeaconTimestampEqual(timestampDb, expectedTs))

		beacon, found := app.BeaconKeeper.GetBeacon(ctx, bID)
		require.True(t, found)
		require.Equal(t, retTsID, beacon.LastTimestampId, "not equal: exp = %d, act = %d", retTsID, beacon.LastTimestampId)
	}

}
