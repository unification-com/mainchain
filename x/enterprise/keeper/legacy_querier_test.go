package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/app/test_helpers"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"strconv"
	"strings"
	"testing"
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

func getQueriedTotalLocked(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) sdk.Coin {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryTotalLocked}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryTotalLocked}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var totalLocked sdk.Coin
	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))

	return totalLocked
}

func getQueriedTotalUnLocked(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) sdk.Coin {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryTotalUnlocked}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryTotalUnlocked}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var totalLocked sdk.Coin
	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))

	return totalLocked
}

func getQueriedTotalSupply(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) sdk.Coins {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryTotalSupply}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryTotalSupply}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var totalLocked sdk.Coins
	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))

	return totalLocked
}

func getQueriedEnterpriseSupply(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) types.UndSupply {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryEnterpriseSupply}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryEnterpriseSupply}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var totalLocked types.UndSupply
	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))

	return totalLocked
}

func getQueriedPurchaseOrder(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, poID uint64) types.EnterpriseUndPurchaseOrder {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var po types.EnterpriseUndPurchaseOrder
	require.NoError(t, cdc.UnmarshalJSON(bz, &po))

	return po
}

func getQueriedPurchaseOrders(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, page, limit int, purchaseOrderStatus types.PurchaseOrderStatus, purchaserAddr sdk.AccAddress) types.QueryResPurchaseOrders {

	params := types.NewQueryPurchaseOrdersParams(page, limit, purchaseOrderStatus, purchaserAddr)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryPurchaseOrders}, "/"),
		Data: cdc.MustMarshalJSON(params),
	}

	bz, err := querier(ctx, []string{keeper.QueryPurchaseOrders}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var matchingOrders types.QueryResPurchaseOrders
	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingOrders))

	return matchingOrders
}

func getQueriedWhitelist(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) []string {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryWhitelist}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{keeper.QueryWhitelist}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var whitelist []string
	require.NoError(t, cdc.UnmarshalJSON(bz, &whitelist))

	return whitelist
}

func getQueriedIsAddressWhitelisted(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, addr sdk.AccAddress) bool {

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, keeper.QueryWhitelisted, addr.String()}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{keeper.QueryWhitelisted, addr.String()}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var isWhitelisted bool
	require.NoError(t, cdc.UnmarshalJSON(bz, &isWhitelisted))

	return isWhitelisted
}

func setupTest() (*app.App, sdk.Context, *codec.LegacyAmino, sdk.Querier) {
	testApp := test_helpers.Setup(false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := testApp.LegacyAmino()
	querier := keeper.NewLegacyQuerier(testApp.EnterpriseKeeper, legacyQuerierCdc)

	return testApp, ctx, legacyQuerierCdc, querier
}

func TestLegacyQueryParams(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	addresses := TestAddrs
	entSigners := addresses[1].String() + "," + addresses[2].String()
	paramsNew := types.NewParams("something", 1, 3600, entSigners)

	testApp.EnterpriseKeeper.SetParams(ctx, paramsNew)
	params := getQueriedParams(t, ctx, legacyQuerierCdc, querier)
	require.True(t, ParamsEqual(paramsNew, params))
}

func TestLegacyQueryTotalLocked(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	denom := TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	_ = testApp.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)

	totalLocked := getQueriedTotalLocked(t, ctx, legacyQuerierCdc, querier)
	require.True(t, locked.IsEqual(totalLocked))
}

func TestLegacyQueryTotalUnLocked(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	denom := TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)
	toUnlock := sdk.NewInt64Coin(denom, int64(500))

	_ = testApp.EnterpriseKeeper.MintCoinsAndLock(ctx, TestAddrs[0], locked)
	_ = testApp.EnterpriseKeeper.UnlockCoinsForFees(ctx, TestAddrs[0], sdk.NewCoins(toUnlock))

	totalUnLocked := getQueriedTotalUnLocked(t, ctx, legacyQuerierCdc, querier)

	require.Equal(t, toUnlock, totalUnLocked)
}

func TestLegacyQueryTotalSupply(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	denom := TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)
	toUnlock := sdk.NewInt64Coin(denom, int64(500))

	_ = testApp.EnterpriseKeeper.MintCoinsAndLock(ctx, TestAddrs[0], locked)
	_ = testApp.EnterpriseKeeper.UnlockCoinsForFees(ctx, TestAddrs[0], sdk.NewCoins(toUnlock))

	totalSupplyFromEnt := getQueriedTotalSupply(t, ctx, legacyQuerierCdc, querier)
	require.Equal(t, sdk.Coins{toUnlock}, totalSupplyFromEnt)
}

func TestLegacyQueryEnterpriseSupply(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	denom := TestDenomination
	amount := int64(1000)
	toMint := sdk.NewInt64Coin(denom, amount)
	toUnlock := sdk.NewInt64Coin(denom, int64(500))
	stillLockedAfterUnlock := toMint.Sub(toUnlock)

	_ = testApp.EnterpriseKeeper.MintCoinsAndLock(ctx, TestAddrs[0], toMint)
	_ = testApp.EnterpriseKeeper.UnlockCoinsForFees(ctx, TestAddrs[0], sdk.NewCoins(toUnlock))

	totalSupplyFromEnt := getQueriedEnterpriseSupply(t, ctx, legacyQuerierCdc, querier)
	require.Equal(t, toUnlock.Amount.Uint64(), totalSupplyFromEnt.Amount)
	require.Equal(t, stillLockedAfterUnlock.Amount.Uint64(), totalSupplyFromEnt.Locked)
	require.Equal(t, toMint.Amount.Uint64(), totalSupplyFromEnt.Total)
}

func TestLegacyQueryPurchaseOrder(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	testAddrs := TestAddrs
	var testPos []types.EnterpriseUndPurchaseOrder

	for i, addr := range testAddrs {
		_ = testApp.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		poId := i + 1
		po := types.EnterpriseUndPurchaseOrder{}
		po.Purchaser = addr.String()
		po.Id = uint64(poId)
		po.Amount = sdk.NewInt64Coin(TestDenomination, int64(i+1))
		po.Status = RandomStatus()

		err := testApp.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
		require.NoError(t, err)
		testPos = append(testPos, po)
	}

	for i, tPo := range testPos {
		poId := i + 1
		po := getQueriedPurchaseOrder(t, ctx, legacyQuerierCdc, querier, uint64(poId))
		require.True(t, tPo.Id == po.Id)
		require.True(t, tPo.Purchaser == po.Purchaser)
		require.True(t, tPo.Amount.String() == po.Amount.String())
		require.True(t, tPo.Status == po.Status)
	}
}

func TestLegacyQueryPurchaseOrdersFilters(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	numTests := 100
	testAddrs := GenerateRandomAccounts(numTests)

	var testPosAccepted []types.EnterpriseUndPurchaseOrder
	var testPosRejected []types.EnterpriseUndPurchaseOrder
	var testPosRaised []types.EnterpriseUndPurchaseOrder
	var testPosCompleted []types.EnterpriseUndPurchaseOrder

	for i, addr := range testAddrs {
		_ = testApp.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		status := RandomStatus()
		poId := i + 1
		po := types.EnterpriseUndPurchaseOrder{}
		po.Purchaser = addr.String()
		po.Id = uint64(poId)
		po.Amount = sdk.NewInt64Coin(TestDenomination, int64(i+1))
		po.Status = status

		err := testApp.EnterpriseKeeper.SetPurchaseOrder(ctx, po)
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

	queryRaised := getQueriedPurchaseOrders(t, ctx, legacyQuerierCdc, querier, 1, numTests, types.StatusRaised, sdk.AccAddress{})
	require.True(t, len(queryRaised) == len(testPosRaised))

	queryCompleted := getQueriedPurchaseOrders(t, ctx, legacyQuerierCdc, querier, 1, numTests, types.StatusCompleted, sdk.AccAddress{})
	require.True(t, len(queryCompleted) == len(testPosCompleted))

	queryAccepted := getQueriedPurchaseOrders(t, ctx, legacyQuerierCdc, querier, 1, numTests, types.StatusAccepted, sdk.AccAddress{})
	require.True(t, len(queryAccepted) == len(testPosAccepted))

	queryRejected := getQueriedPurchaseOrders(t, ctx, legacyQuerierCdc, querier, 1, numTests, types.StatusRejected, sdk.AccAddress{})
	require.True(t, len(queryRejected) == len(testPosRejected))

	for _, tPo := range testPosRaised {
		for _, po := range queryRaised {
			if po.Id == tPo.Id {
				require.True(t, tPo.Id == po.Id)
				require.True(t, tPo.Purchaser == po.Purchaser)
				require.True(t, tPo.Amount.String() == po.Amount.String())
				require.True(t, tPo.Status == po.Status)
			}
		}
	}

	for _, tPo := range testPosCompleted {
		for _, po := range queryCompleted {
			if po.Id == tPo.Id {
				require.True(t, tPo.Id == po.Id)
				require.True(t, tPo.Purchaser == po.Purchaser)
				require.True(t, tPo.Amount.String() == po.Amount.String())
				require.True(t, tPo.Status == po.Status)
			}
		}
	}

	for _, tPo := range testPosAccepted {
		for _, po := range queryAccepted {
			if po.Id == tPo.Id {
				require.True(t, tPo.Id == po.Id)
				require.True(t, tPo.Purchaser == po.Purchaser)
				require.True(t, tPo.Amount.String() == po.Amount.String())
				require.True(t, tPo.Status == po.Status)
			}
		}
	}

	for _, tPo := range testPosRejected {
		for _, po := range queryRejected {
			if po.Id == tPo.Id {
				require.True(t, tPo.Id == po.Id)
				require.True(t, tPo.Purchaser == po.Purchaser)
				require.True(t, tPo.Amount.String() == po.Amount.String())
				require.True(t, tPo.Status == po.Status)
			}
		}
	}
}

func TestLegacyQueryWhitelist(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()

	numTests := 100
	testAddrs := GenerateRandomAccounts(numTests)

	for _, addr := range testAddrs {
		_ = testApp.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
	}

	whitelist := getQueriedWhitelist(t, ctx, legacyQuerierCdc, querier)

	require.True(t, len(whitelist) == len(testAddrs))
}

func TestLegacyQueryAddressIsWhitelisted(t *testing.T) {
	testApp, ctx, legacyQuerierCdc, querier := setupTest()
	numTests := 100
	testAddrs := GenerateRandomAccounts(numTests)

	for i, addr := range testAddrs {
		if i < 50 {
			_ = testApp.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		}
	}

	for i, addr := range testAddrs {
		isWhiteListed := getQueriedIsAddressWhitelisted(t, ctx, legacyQuerierCdc, querier, addr)
		if i < 50 {
			require.True(t, isWhiteListed)
		} else {
			require.False(t, isWhiteListed)
		}
	}
}
