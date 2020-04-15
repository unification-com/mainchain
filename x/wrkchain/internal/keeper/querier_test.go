package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
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

func getQueriedWrkChain(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, wcID uint64) types.WrkChain {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryWrkChain, strconv.FormatUint(wcID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{QueryWrkChain, strconv.FormatUint(wcID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var wc types.WrkChain
	require.NoError(t, cdc.UnmarshalJSON(bz, &wc))

	return wc
}

func getQueriedWrkChainBlock(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, wcID, height uint64) types.WrkChainBlock {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryWrkChainBlock, strconv.FormatUint(wcID, 10), strconv.FormatUint(height, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{QueryWrkChainBlock, strconv.FormatUint(wcID, 10), strconv.FormatUint(height, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var wcb types.WrkChainBlock
	require.NoError(t, cdc.UnmarshalJSON(bz, &wcb))

	return wcb
}

func getQueriedWrkChainsFiltered(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, page, limit int, moniker string, owner sdk.AccAddress) types.QueryResWrkChains {

	params := types.NewQueryWrkChainParams(page, limit, moniker, owner)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryWrkChainsFiltered}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{QueryWrkChainsFiltered}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingWcs types.QueryResWrkChains
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

func TestInvalidQuerier(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, "nosuchpath"}, "/"),
		Data: []byte{},
	}

	expextedErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: nosuchpath")

	_, err := querier(ctx, []string{"nosuchpath"}, query)

	require.Equal(t, expextedErr.Error(), err.Error())
}

func TestQueryWrkChainByID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	var testWcs []types.WrkChain

	for i, addr := range TestAddrs {
		wcID := uint64(i + 1)
		wc := types.NewWrkChain()
		wc.Owner = addr
		wc.WrkChainID = wcID
		wc.LastBlock = 0
		wc.RegisterTime = time.Now().Unix()
		wc.Moniker = GenerateRandomString(12)
		wc.Name = GenerateRandomString(20)
		wc.GenesisHash = GenerateRandomString(32)

		err := keeper.SetWrkChain(ctx, wc)
		require.NoError(t, err)
		testWcs = append(testWcs, wc)
	}

	for i, tWc := range testWcs {
		wcID := i + 1
		wc := getQueriedWrkChain(t, ctx, keeper.cdc, querier, uint64(wcID))
		require.True(t, WRKChainEqual(tWc, wc))
	}
}

func TestQueryWrkChainBlockByHeight(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	var testWcBlocks []types.WrkChainBlock
	numBlocks := uint64(100)
	wcID := uint64(1)
	addr := TestAddrs[0]

	wc := types.NewWrkChain()
	wc.Owner = addr
	wc.WrkChainID = wcID
	wc.LastBlock = 0
	wc.RegisterTime = time.Now().Unix()
	wc.Moniker = GenerateRandomString(12)
	wc.Name = GenerateRandomString(20)
	wc.GenesisHash = GenerateRandomString(32)

	err := keeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for h := uint64(1); h <= numBlocks; h++ {
		block := types.NewWrkChainBlock()
		block.WrkChainID = wcID
		block.Owner = addr
		block.Height = h
		block.BlockHash = GenerateRandomString(32)
		block.ParentHash = GenerateRandomString(32)
		block.Hash1 = GenerateRandomString(32)
		block.Hash2 = GenerateRandomString(32)
		block.Hash3 = GenerateRandomString(32)
		block.SubmitTime = time.Now().Unix()

		err := keeper.SetWrkChainBlock(ctx, block)
		require.NoError(t, err)
		testWcBlocks = append(testWcBlocks, block)
	}

	for i, tWcb := range testWcBlocks {
		height := uint64(i + 1)
		block := getQueriedWrkChainBlock(t, ctx, keeper.cdc, querier, wcID, height)
		require.True(t, WRKChainBlockEqual(tWcb, block))
	}
}

func TestQueryWrkChainsFiltered(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)
	querier := NewQuerier(keeper)
	numToRegister := 10

	for _, addr := range TestAddrs {
		var monikers []string
		var testWcs []types.WrkChain

		for i := 0; i < numToRegister; i++ {
			moniker := GenerateRandomString(12)
			name := GenerateRandomString(20)
			genesisHash := GenerateRandomString(32)
			wcID, _ := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)

			wc := keeper.GetWrkChain(ctx, wcID)

			monikers = append(monikers, moniker)
			testWcs = append(testWcs, wc)
		}

		wcsForAddress := getQueriedWrkChainsFiltered(t, ctx, keeper.cdc, querier, 1, 100, "", addr)
		require.True(t, len(wcsForAddress) == numToRegister && len(wcsForAddress) == len(testWcs))
		for _, tWc := range testWcs {
			for _, wc := range wcsForAddress {
				if wc.WrkChainID == tWc.WrkChainID {
					require.True(t, WRKChainEqual(tWc, wc))
				}
			}
		}

		for _, m := range monikers {
			wcForMoniker := getQueriedWrkChainsFiltered(t, ctx, keeper.cdc, querier, 1, 100, m, sdk.AccAddress{})
			require.True(t, len(wcForMoniker) == 1)
			require.True(t, wcForMoniker[0].Owner.String() == addr.String())
			for _, tWc := range testWcs {
				if tWc.Moniker == wcForMoniker[0].Moniker {
					require.True(t, WRKChainEqual(wcForMoniker[0], tWc))
				}
			}
		}
	}
}
