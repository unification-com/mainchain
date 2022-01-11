package keeper_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
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

func getQueriedWrkChain(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, wcID uint64) types.WrkChain {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryWrkChain, strconv.FormatUint(wcID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryWrkChain, strconv.FormatUint(wcID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var wc types.WrkChain
	require.NoError(t, cdc.UnmarshalJSON(bz, &wc))

	return wc
}

func getQueriedWrkChainBlock(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, wcID, height uint64) types.WrkChainBlock {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryWrkChainBlock, strconv.FormatUint(wcID, 10), strconv.FormatUint(height, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryWrkChainBlock, strconv.FormatUint(wcID, 10), strconv.FormatUint(height, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var wcb types.WrkChainBlock
	require.NoError(t, cdc.UnmarshalJSON(bz, &wcb))

	return wcb
}

func getQueriedWrkChainsFiltered(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, page, limit int, moniker string, owner sdk.AccAddress) types.QueryResWrkChains {

	params := types.QueryWrkChainsFilteredRequest{
		Moniker:    moniker,
		Owner:      owner.String(),
		Pagination: nil,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryWrkChainsFiltered}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{keeper.QueryWrkChainsFiltered}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingWcs types.QueryResWrkChains
	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingWcs))

	return matchingWcs
}

func setupTest() (*app.App, sdk.Context, *codec.LegacyAmino, sdk.Querier) {
	testApp := test_helpers.Setup(false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := testApp.LegacyAmino()
	querier := keeper.NewLegacyQuerier(testApp.WrkchainKeeper, legacyQuerierCdc)

	return testApp, ctx, legacyQuerierCdc, querier
}

func TestQueryParams(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	paramsNew := types.NewParams(9999, 999, "somecoin")

	testApp.WrkchainKeeper.SetParams(ctx, paramsNew)
	params := getQueriedParams(t, ctx, legacyQuerierCdc, querier)
	require.True(t, ParamsEqual(paramsNew, params))
}

func TestInvalidQuerier(t *testing.T) {
	_, ctx, _, querier := setupTest()

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, "nosuchpath"}, "/"),
		Data: []byte{},
	}

	expextedErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: nosuchpath")

	_, err := querier(ctx, []string{"nosuchpath"}, query)

	require.Equal(t, expextedErr.Error(), err.Error())
}

func TestQueryWrkChainByID(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	var testWcs []types.WrkChain

	for i, addr := range TestAddrs {
		wcID := uint64(i + 1)
		wc := types.WrkChain{}
		wc.Owner = addr.String()
		wc.WrkchainId = wcID
		wc.Lastblock = 0
		wc.RegTime = uint64(time.Now().Unix())
		wc.Moniker = GenerateRandomString(12)
		wc.Name = GenerateRandomString(20)
		wc.Genesis = GenerateRandomString(32)

		err := testApp.WrkchainKeeper.SetWrkChain(ctx, wc)
		require.NoError(t, err)
		testWcs = append(testWcs, wc)
	}

	for i, tWc := range testWcs {
		wcID := i + 1
		wc := getQueriedWrkChain(t, ctx, legacyQuerierCdc, querier, uint64(wcID))
		require.True(t, WRKChainEqual(tWc, wc))
	}
}

func TestQueryWrkChainBlockByHeight(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	var testWcBlocks []types.WrkChainBlock
	numBlocks := uint64(100)
	wcID := uint64(1)
	addr := TestAddrs[0]

	wc := types.WrkChain{}
	wc.Owner = addr.String()
	wc.WrkchainId = wcID
	wc.Lastblock = 0
	wc.RegTime = uint64(time.Now().Unix())
	wc.Moniker = GenerateRandomString(12)
	wc.Name = GenerateRandomString(20)
	wc.Genesis = GenerateRandomString(32)

	err := testApp.WrkchainKeeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for h := uint64(1); h <= numBlocks; h++ {
		block := types.WrkChainBlock{}
		block.WrkchainId = wcID
		block.Owner = addr.String()
		block.Height = h
		block.Blockhash = GenerateRandomString(32)
		block.Parenthash = GenerateRandomString(32)
		block.Hash1 = GenerateRandomString(32)
		block.Hash2 = GenerateRandomString(32)
		block.Hash3 = GenerateRandomString(32)
		block.SubTime = uint64(time.Now().Unix())

		err = testApp.WrkchainKeeper.SetWrkChainBlock(ctx, block)
		require.NoError(t, err)
		testWcBlocks = append(testWcBlocks, block)
	}

	for i, tWcb := range testWcBlocks {
		height := uint64(i + 1)
		block := getQueriedWrkChainBlock(t, ctx, legacyQuerierCdc, querier, wcID, height)
		require.True(t, WRKChainBlockEqual(tWcb, block))
	}
}

func TestQueryWrkChainsFiltered(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	numToRegister := 10

	for _, addr := range TestAddrs {
		var monikers []string
		var testWcs []types.WrkChain

		for i := 0; i < numToRegister; i++ {
			moniker := GenerateRandomString(12)
			name := GenerateRandomString(20)
			genesisHash := GenerateRandomString(32)
			wcID, _ := testApp.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)

			wc, ok := testApp.WrkchainKeeper.GetWrkChain(ctx, wcID)
			require.True(t, ok)

			monikers = append(monikers, moniker)
			testWcs = append(testWcs, wc)
		}

		wcsForAddress := getQueriedWrkChainsFiltered(t, ctx, legacyQuerierCdc, querier, 1, 100, "", addr)
		require.True(t, len(wcsForAddress) == numToRegister && len(wcsForAddress) == len(testWcs))
		for _, tWc := range testWcs {
			for _, wc := range wcsForAddress {
				if wc.WrkchainId == tWc.WrkchainId {
					require.True(t, WRKChainEqual(tWc, wc))
				}
			}
		}

		for _, m := range monikers {
			wcForMoniker := getQueriedWrkChainsFiltered(t, ctx, legacyQuerierCdc, querier, 1, 100, m, sdk.AccAddress{})
			require.True(t, len(wcForMoniker) == 1)
			require.True(t, wcForMoniker[0].Owner == addr.String())
			for _, tWc := range testWcs {
				if tWc.Moniker == wcForMoniker[0].Moniker {
					require.True(t, WRKChainEqual(wcForMoniker[0], tWc))
				}
			}
		}
	}
}
