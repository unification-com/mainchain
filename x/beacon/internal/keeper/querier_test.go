package keeper

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

const custom = "custom"

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.Params {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryParameters}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{QueryParameters}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var params types.Params
	require.NoError(t, cdc.UnmarshalJSON(bz, &params))

	return params
}

func getQueriedBeacon(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, bID uint64) types.Beacon {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryBeacon, strconv.FormatUint(bID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{QueryBeacon, strconv.FormatUint(bID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var b types.Beacon
	require.NoError(t, cdc.UnmarshalJSON(bz, &b))

	return b
}

func getQueriedBeaconTimestamp(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, bID, tsID uint64) types.BeaconTimestamp {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryBeaconTimestamp, strconv.FormatUint(bID, 10), strconv.FormatUint(tsID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{QueryBeaconTimestamp, strconv.FormatUint(bID, 10), strconv.FormatUint(tsID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var bts types.BeaconTimestamp
	require.NoError(t, cdc.UnmarshalJSON(bz, &bts))

	return bts
}

func getQueriedBeaconTimestamps(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, params types.QueryBeaconTimestampParams) types.QueryResBeaconTimestampHashes {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryBeaconTimestamps}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{QueryBeaconTimestamps}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var bts types.QueryResBeaconTimestampHashes
	require.NoError(t, cdc.UnmarshalJSON(bz, &bts))

	return bts
}

func getQueriedBeaconsFiltered(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, page, limit int, moniker string, owner sdk.AccAddress) types.QueryResBeacons {

	params := types.NewQueryBeaconParams(page, limit, moniker, owner)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryBeacons}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{QueryBeacons}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingWcs types.QueryResBeacons
	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingWcs))

	return matchingWcs
}

func TestQueryParams(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	querier := NewQuerier(keeper)
	paramsNew := types.NewParams(9999, 999, "somecoin")

	keeper.SetParams(ctx, paramsNew)
	params := getQueriedParams(t, ctx, keeper.cdc, querier)
	require.True(t, ParamsEqual(paramsNew, params))
}

func TestQueryBeaconByID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	var testBeacons []types.Beacon

	for i, addr := range TestAddrs {
		bID := uint64(i + 1)
		b := types.NewBeacon()
		b.Owner = addr
		b.BeaconID = bID
		b.LastTimestampID = 0
		b.Moniker = GenerateRandomString(12)
		b.Name = GenerateRandomString(20)

		err := keeper.SetBeacon(ctx, b)
		require.NoError(t, err)
		testBeacons = append(testBeacons, b)
	}

	for i, tB := range testBeacons {
		bID := i + 1
		b := getQueriedBeacon(t, ctx, keeper.cdc, querier, uint64(bID))
		require.True(t, BeaconEqual(tB, b))
	}
}

func TestQueryBeaconTimestampByID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	var testBeaconTs []types.BeaconTimestamp
	numTimestamps := uint64(100)
	bID := uint64(1)
	addr := TestAddrs[0]

	b := types.NewBeacon()
	b.Owner = addr
	b.BeaconID = bID
	b.LastTimestampID = 0
	b.Moniker = GenerateRandomString(12)
	b.Name = GenerateRandomString(20)

	err := keeper.SetBeacon(ctx, b)
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numTimestamps; tsID++ {
		ts := types.NewBeaconTimestamp()
		ts.BeaconID = bID
		ts.Owner = addr
		ts.TimestampID = tsID
		ts.Hash = GenerateRandomString(32)
		ts.SubmitTime = uint64(time.Now().Unix())

		err := keeper.SetBeaconTimestamp(ctx, ts)
		require.NoError(t, err)
		testBeaconTs = append(testBeaconTs, ts)
	}

	for i, tBts := range testBeaconTs {
		tsID := uint64(i + 1)
		block := getQueriedBeaconTimestamp(t, ctx, keeper.cdc, querier, bID, tsID)
		require.True(t, BeaconTimestampEqual(tBts, block))
	}
}

func TestQueryBeaconTimestamps(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	var testBeaconTs []types.BeaconTimestamp
	numTimestamps := uint64(100)
	bID := uint64(1)
	addr := TestAddrs[0]

	b := types.NewBeacon()
	b.Owner = addr
	b.BeaconID = bID
	b.LastTimestampID = 0
	b.Moniker = GenerateRandomString(12)
	b.Name = GenerateRandomString(20)

	err := keeper.SetBeacon(ctx, b)
	require.NoError(t, err)

	for tsID := uint64(1); tsID <= numTimestamps; tsID++ {
		ts := types.NewBeaconTimestamp()
		ts.BeaconID = bID
		ts.Owner = addr
		ts.TimestampID = tsID
		ts.Hash = GenerateRandomString(32)
		ts.SubmitTime = uint64(time.Now().Unix())

		err := keeper.SetBeaconTimestamp(ctx, ts)
		require.NoError(t, err)
		testBeaconTs = append(testBeaconTs, ts)
	}

	params := types.NewQueryBeaconTimestampParams(1, 100, bID, "", 0)
	allBlocks := getQueriedBeaconTimestamps(t, ctx, keeper.cdc, querier, params)

	require.True(t, len(allBlocks) == int(numTimestamps) && len(allBlocks) == len(testBeaconTs))

	for _, tBts := range testBeaconTs {
		for _, b := range allBlocks {
			if b.TimestampID == tBts.TimestampID {
				require.True(t, BeaconTimestampEqual(tBts, b))
			}
		}
	}

}

func TestQueryBeaconsFiltered(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	numToRegister := 10

	for _, addr := range TestAddrs {
		var monikers []string
		var testBeacons []types.Beacon

		for i := 0; i < numToRegister; i++ {
			moniker := GenerateRandomString(12)
			name := GenerateRandomString(20)
			bID, _ := keeper.RegisterBeacon(ctx, moniker, name, addr)

			b := keeper.GetBeacon(ctx, bID)

			monikers = append(monikers, moniker)
			testBeacons = append(testBeacons, b)
		}

		beaconsForAddress := getQueriedBeaconsFiltered(t, ctx, keeper.cdc, querier, 1, 100, "", addr)
		require.True(t, len(beaconsForAddress) == numToRegister && len(beaconsForAddress) == len(testBeacons))
		for _, tB := range testBeacons {
			for _, b := range beaconsForAddress {
				if b.BeaconID == tB.BeaconID {
					require.True(t, BeaconEqual(tB, b))
				}
			}
		}

		for _, m := range monikers {
			beaconForMoniker := getQueriedBeaconsFiltered(t, ctx, keeper.cdc, querier, 1, 100, m, sdk.AccAddress{})
			require.True(t, len(beaconForMoniker) == 1)
			require.True(t, beaconForMoniker[0].Owner.String() == addr.String())
			for _, tB := range testBeacons {
				if tB.Moniker == beaconForMoniker[0].Moniker {
					require.True(t, BeaconEqual(beaconForMoniker[0], tB))
				}
			}
		}
	}
}
