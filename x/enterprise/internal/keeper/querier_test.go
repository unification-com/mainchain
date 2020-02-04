package keeper

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
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

func getQueriedPurchaseOrder(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, poID uint64) types.EnterpriseUndPurchaseOrder {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var po types.EnterpriseUndPurchaseOrder
	require.NoError(t, cdc.UnmarshalJSON(bz, &po))

	return po
}

func getQueriedPurchaseOrders(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, page, limit int, purchaseOrderStatus types.PurchaseOrderStatus, purchaserAddr sdk.AccAddress) types.QueryResPurchaseOrders {

	params := types.NewQueryPurchaseOrdersParams(page, limit, purchaseOrderStatus, purchaserAddr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, QueryPurchaseOrders}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{QueryPurchaseOrders}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingOrders types.QueryResPurchaseOrders
	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingOrders))

	return matchingOrders
}

func TestQueryParams(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	querier := NewQuerier(keeper)
	addresses := GenerateRandomAddresses(3)
	entSigners := addresses[1].String() + "," + addresses[2].String()
	paramsNew := types.NewParams("something", 1, 3600, entSigners)

	keeper.SetParams(ctx, paramsNew)
	params := getQueriedParams(t, ctx, keeper.cdc, querier)
	require.True(t, ParamsEqual(paramsNew, params))
}

func TestQueryPurchaseOrder(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	querier := NewQuerier(keeper)
	testAddrs := GenerateRandomAddresses(100)
	var testPos []types.EnterpriseUndPurchaseOrder

	for i, addr := range testAddrs {
		poId := i + 1
		po := types.NewEnterpriseUndPurchaseOrder()
		po.Purchaser = addr
		po.PurchaseOrderID = uint64(poId)
		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i+1))
		po.Status = RandomStatus()

		err := keeper.SetPurchaseOrder(ctx, po)
		require.NoError(t, err)
		testPos = append(testPos, po)
	}

	for i, tPo := range testPos {
		poId := i + 1
		po := getQueriedPurchaseOrder(t, ctx, keeper.cdc, querier, uint64(poId))
		require.True(t, PurchaseOrderEqual(tPo, po))
	}
}

func TestQueryPurchaseOrdersFilters(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	querier := NewQuerier(keeper)
	numTests := 100
	testAddrs := GenerateRandomAddresses(numTests)

	var testPosAccepted []types.EnterpriseUndPurchaseOrder
	var testPosRejected []types.EnterpriseUndPurchaseOrder
	var testPosRaised []types.EnterpriseUndPurchaseOrder
	var testPosCompleted []types.EnterpriseUndPurchaseOrder

	for i, addr := range testAddrs {
		status := RandomStatus()
		poId := i + 1
		po := types.NewEnterpriseUndPurchaseOrder()
		po.Purchaser = addr
		po.PurchaseOrderID = uint64(poId)
		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i+1))
		po.Status = status

		err := keeper.SetPurchaseOrder(ctx, po)
		require.NoError(t, err)

		switch status {
		case types.StatusRaised:
			testPosRaised = append(testPosRaised, po)
		case types.StatusCompleted:
			testPosCompleted = append(testPosCompleted, po)
		case types.StatusAccepted:
			testPosAccepted = append(testPosAccepted, po)
		case types.StatusRejected:
			testPosRejected = append(testPosRejected, po)
		}
	}

	queryRaised := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusRaised, sdk.AccAddress{})
	require.True(t, len(queryRaised) == len(testPosRaised))

	queryCompleted := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusCompleted, sdk.AccAddress{})
	require.True(t, len(queryCompleted) == len(testPosCompleted))

	queryAccepted := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusAccepted, sdk.AccAddress{})
	require.True(t, len(queryAccepted) == len(testPosAccepted))

	queryRejected := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusRejected, sdk.AccAddress{})
	require.True(t, len(queryRejected) == len(testPosRejected))

	for _, tPo := range testPosRaised {
		for _, po := range queryRaised {
			if po.PurchaseOrderID == tPo.PurchaseOrderID {
				require.True(t, PurchaseOrderEqual(tPo, po))
			}
		}
	}

	for _, tPo := range testPosCompleted {
		for _, po := range queryCompleted {
			if po.PurchaseOrderID == tPo.PurchaseOrderID {
				require.True(t, PurchaseOrderEqual(tPo, po))
			}
		}
	}

	for _, tPo := range testPosAccepted {
		for _, po := range queryAccepted {
			if po.PurchaseOrderID == tPo.PurchaseOrderID {
				require.True(t, PurchaseOrderEqual(tPo, po))
			}
		}
	}

	for _, tPo := range testPosRejected {
		for _, po := range queryRejected {
			if po.PurchaseOrderID == tPo.PurchaseOrderID {
				require.True(t, PurchaseOrderEqual(tPo, po))
			}
		}
	}
}
