package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/app/test_helpers"
	"strconv"
	"strings"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/keeper"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const custom = "custom"

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) types.Params {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryParameters}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryParameters}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var params types.Params
	require.NoError(t, cdc.UnmarshalJSON(bz, &params))

	return params
}

func getQueriedBeacon(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, bID uint64) types.Beacon {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryBeacon, strconv.FormatUint(bID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryBeacon, strconv.FormatUint(bID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var b types.Beacon
	require.NoError(t, cdc.UnmarshalJSON(bz, &b))

	return b
}

func getQueriedBeaconTimestamp(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, bID, tsID uint64) types.BeaconTimestampLegacy {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryBeaconTimestamp, strconv.FormatUint(bID, 10), strconv.FormatUint(tsID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryBeaconTimestamp, strconv.FormatUint(bID, 10), strconv.FormatUint(tsID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var bts types.BeaconTimestampLegacy
	require.NoError(t, cdc.UnmarshalJSON(bz, &bts))

	return bts
}

func getQueriedBeaconsFiltered(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, page, limit int, moniker string, owner sdk.AccAddress) types.QueryResBeacons {

	params := types.QueryBeaconsFilteredRequest{
		Moniker:    moniker,
		Owner:      owner.String(),
		Pagination: nil,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryBeacons}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{keeper.QueryBeacons}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingWcs types.QueryResBeacons
	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingWcs))
	return matchingWcs
}

func setupTest(t *testing.T) (*app.App, sdk.Context, *codec.LegacyAmino, sdk.Querier) {
	testApp := test_helpers.Setup(t, false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := testApp.LegacyAmino()
	querier := keeper.NewLegacyQuerier(testApp.BeaconKeeper, legacyQuerierCdc)

	return testApp, ctx, legacyQuerierCdc, querier
}

func TestQueryParams(t *testing.T) {

	testApp, ctx, legacyQuerierCdc, querier := setupTest(t)

	paramsNew := types.NewParams(9999, 999, 999, "somecoin", 9999, 99999)
	testApp.BeaconKeeper.SetParams(ctx, paramsNew)

	params := getQueriedParams(t, ctx, legacyQuerierCdc, querier)
	require.True(t, ParamsEqual(paramsNew, params))
}

func TestInvalidQuerier(t *testing.T) {

	_, ctx, _, querier := setupTest(t)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, "nosuchpath"}, "/"),
		Data: []byte{},
	}

	expextedErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: nosuchpath")

	_, err := querier(ctx, []string{"nosuchpath"}, query)

	require.Equal(t, expextedErr.Error(), err.Error())
}

func TestQueryBeaconByID(t *testing.T) {

	testApp, ctx, legacyQuerierCdc, querier := setupTest(t)

	var testBeacons []types.Beacon

	for i, addr := range TestAddrs {
		bID := uint64(i + 1)
		b := types.Beacon{}
		b.Owner = addr.String()
		b.BeaconId = bID
		b.LastTimestampId = 0
		b.Moniker = GenerateRandomString(12)
		b.Name = GenerateRandomString(20)

		err := testApp.BeaconKeeper.SetBeacon(ctx, b)
		require.NoError(t, err)
		testBeacons = append(testBeacons, b)
	}

	for i, tB := range testBeacons {
		bID := i + 1
		b := getQueriedBeacon(t, ctx, legacyQuerierCdc, querier, uint64(bID))
		require.True(t, BeaconEqual(tB, b))
	}
}

func TestQueryBeaconTimestampByID(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest(t)

	var testBeaconTs []types.BeaconTimestampLegacy
	numTimestamps := uint64(100)
	bID := uint64(1)
	addr := TestAddrs[0]

	b := types.Beacon{}
	b.Owner = addr.String()
	b.BeaconId = bID
	b.LastTimestampId = 0
	b.Moniker = GenerateRandomString(12)
	b.Name = GenerateRandomString(20)

	err := testApp.BeaconKeeper.SetBeacon(ctx, b)
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numTimestamps; tsID++ {
		ts := types.BeaconTimestamp{}
		ts.TimestampId = tsID
		ts.Hash = GenerateRandomString(32)
		ts.SubmitTime = uint64(time.Now().Unix())

		err := testApp.BeaconKeeper.SetBeaconTimestamp(ctx, bID, ts)
		require.NoError(t, err)

		expectedTsRes := types.BeaconTimestampLegacy{
			BeaconID:    bID,
			TimestampID: tsID,
			SubmitTime:  ts.SubmitTime,
			Hash:        ts.Hash,
			Owner:       addr.String(),
		}

		testBeaconTs = append(testBeaconTs, expectedTsRes)
	}

	for i, tBts := range testBeaconTs {
		tsID := uint64(i + 1)
		ts := getQueriedBeaconTimestamp(t, ctx, legacyQuerierCdc, querier, bID, tsID)
		require.True(t, BeaconTimestampLegacyEqual(tBts, ts))
	}
}

func TestQueryBeaconsFiltered(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest(t)

	numToRegister := 10

	for _, addr := range TestAddrs {
		var monikers []string
		var testBeacons []types.Beacon

		for i := 0; i < numToRegister; i++ {
			moniker := GenerateRandomString(12)
			name := GenerateRandomString(20)

			b := types.Beacon{}
			b.Owner = addr.String()
			b.Moniker = moniker
			b.Name = name

			bID, _ := testApp.BeaconKeeper.RegisterNewBeacon(ctx, b)

			bDb, _ := testApp.BeaconKeeper.GetBeacon(ctx, bID)

			monikers = append(monikers, moniker)
			testBeacons = append(testBeacons, bDb)
		}

		beaconsForAddress := getQueriedBeaconsFiltered(t, ctx, legacyQuerierCdc, querier, 1, 100, "", addr)
		require.True(t, len(beaconsForAddress) == numToRegister && len(beaconsForAddress) == len(testBeacons))
		for _, tB := range testBeacons {
			for _, b := range beaconsForAddress {
				if b.BeaconId == tB.BeaconId {
					require.True(t, BeaconEqual(tB, b))
				}
			}
		}

		for _, m := range monikers {
			beaconForMoniker := getQueriedBeaconsFiltered(t, ctx, legacyQuerierCdc, querier, 1, 100, m, sdk.AccAddress{})
			require.True(t, len(beaconForMoniker) >= 1)
			require.True(t, beaconForMoniker[0].Owner == addr.String())
			for _, tB := range testBeacons {
				if tB.Moniker == beaconForMoniker[0].Moniker {
					require.True(t, BeaconEqual(beaconForMoniker[0], tB))
				}
			}
		}
	}
}
