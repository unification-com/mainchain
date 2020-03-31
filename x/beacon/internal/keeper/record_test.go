package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

func TestSetGetBeaconTimestamp(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, addr)
		require.NoError(t, err)

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			beaconTimestamp := types.NewBeaconTimestamp()
			beaconTimestamp.BeaconID = bID
			beaconTimestamp.Owner = addr
			beaconTimestamp.TimestampID = tsID
			beaconTimestamp.Hash = GenerateRandomString(32)
			beaconTimestamp.SubmitTime = uint64(time.Now().Unix())

			err := keeper.SetBeaconTimestamp(ctx, beaconTimestamp)
			require.NoError(t, err)

			btsDb := keeper.GetBeaconTimestampByID(ctx, bID, tsID)
			require.True(t, BeaconTimestampEqual(btsDb, beaconTimestamp))
		}
	}
}

func TestIsBeaconTimestampRecorded(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, addr)
		require.NoError(t, err)

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			hash := GenerateRandomString(32)
			subTime := uint64(time.Now().Unix())
			timestamp := types.NewBeaconTimestamp()
			timestamp.BeaconID = bID
			timestamp.Owner = addr
			timestamp.TimestampID = tsID
			timestamp.Hash = hash
			timestamp.SubmitTime = subTime

			err := keeper.SetBeaconTimestamp(ctx, timestamp)
			require.NoError(t, err)

			isRecorded := keeper.IsBeaconTimestampRecordedByID(ctx, bID, tsID)
			require.True(t, isRecorded)

			isRecorded1 := keeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, hash, 0)
			require.True(t, isRecorded1)

			isRecorded2 := keeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, "", subTime)
			require.True(t, isRecorded2)

			isRecorded3 := keeper.IsBeaconTimestampRecordedByHashTime(ctx, bID, hash, subTime)
			require.True(t, isRecorded3)
		}
	}
}

func TestGetWrkChainBlockHashes(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, addr)
		require.NoError(t, err)

		var testTimestamps []types.BeaconTimestamp

		for tsID := uint64(1); tsID <= numToRecord; tsID++ {
			timestamp := types.NewBeaconTimestamp()
			timestamp.BeaconID = bID
			timestamp.Owner = addr
			timestamp.TimestampID = tsID
			timestamp.Hash = GenerateRandomString(32)
			timestamp.SubmitTime = uint64(time.Now().Unix())

			testTimestamps = append(testTimestamps, timestamp)

			err := keeper.SetBeaconTimestamp(ctx, timestamp)
			require.NoError(t, err)
		}

		allTimestamps := keeper.GetAllBeaconTimestamps(ctx, bID)
		require.True(t, len(allTimestamps) == int(numToRecord) && len(allTimestamps) == len(testTimestamps))

		for i := 0; i < int(numToRecord); i++ {
			require.True(t, BeaconTimestampEqual(allTimestamps[i], testTimestamps[i]))
		}
	}
}

func TestIsAuthorisedToRecord(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	unauthorisedAddr := sdk.AccAddress(pubKey.Address())

	for _, addr := range TestAddrs {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)

		bID, err := keeper.RegisterBeacon(ctx, moniker, name, addr)
		require.NoError(t, err)

		isAuthorised := keeper.IsAuthorisedToRecord(ctx, bID, addr)
		require.True(t, isAuthorised)

		isAuthorised = keeper.IsAuthorisedToRecord(ctx, bID, unauthorisedAddr)
		require.False(t, isAuthorised)
	}
}

func TestRecordBeaconTimestamps(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	name := GenerateRandomString(20)
	moniker := GenerateRandomString(12)

	bID, err := keeper.RegisterBeacon(ctx, moniker, name, TestAddrs[0])
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numToRecord; tsID++ {
		expectedTs := types.NewBeaconTimestamp()
		expectedTs.BeaconID = bID
		expectedTs.Owner = TestAddrs[0]
		expectedTs.TimestampID = tsID
		expectedTs.Hash = GenerateRandomString(32)
		expectedTs.SubmitTime = uint64(time.Now().Unix())

		retTsID, err := keeper.RecordBeaconTimestamp(ctx, bID, expectedTs.Hash, expectedTs.SubmitTime, expectedTs.Owner)
		require.NoError(t, err)
		require.True(t, retTsID == expectedTs.TimestampID)

		timestampDb := keeper.GetBeaconTimestampByID(ctx, bID, tsID)
		require.True(t, BeaconTimestampEqual(timestampDb, expectedTs))

		beacon := keeper.GetBeacon(ctx, bID)
		require.Equal(t, retTsID, beacon.LastTimestampID, "not equal: exp = %d, act = %d", retTsID, beacon.LastTimestampID)
	}

}

func TestRecordBeaconTimestampsFail(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	name := GenerateRandomString(20)
	moniker := GenerateRandomString(12)

	bID, err := keeper.RegisterBeacon(ctx, moniker, name, TestAddrs[0])
	require.NoError(t, err)

	testCases := []struct {
		beaconID    uint64
		subTime     uint64
		hash        string
		owner       sdk.AccAddress
		expectedErr error
		expectedID  uint64
	}{
		{0, 0, "", sdk.AccAddress{}, types.ErrBeaconDoesNotExist, 0},
		{99, 0, "", sdk.AccAddress{}, types.ErrBeaconDoesNotExist, 0},
		{bID, 1, "hash", TestAddrs[1], sdkerrors.Wrapf(types.ErrNotBeaconOwner, "%s not authorised to record hashes for this beacon", TestAddrs[1]), 0},
		{bID, 1, "hash", sdk.AccAddress{}, sdkerrors.Wrapf(types.ErrNotBeaconOwner, "%s not authorised to record hashes for this beacon", sdk.AccAddress{}), 0},
		{bID, 1, "", TestAddrs[0], sdkerrors.Wrap(types.ErrMissingData, "must include owner, id, submit time and hash"), 0},
		{bID, 0, "timstamphash", TestAddrs[0], sdkerrors.Wrap(types.ErrMissingData, "must include owner, id, submit time and hash"), 0},
		{bID, 1, "timstamphash", TestAddrs[0], nil, 1},
		{bID, 1, "timstamphash", TestAddrs[0], sdkerrors.Wrap(types.ErrBeaconTimestampAlreadyRecorded, "timestamp hash timstamphash already recorded at time 1"), 0},
	}

	for _, tc := range testCases {
		tsID, err := keeper.RecordBeaconTimestamp(ctx, tc.beaconID, tc.hash, tc.subTime, tc.owner)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
		} else {
			require.Nil(t, err)
		}
		require.True(t, tsID == tc.expectedID)
	}
}
