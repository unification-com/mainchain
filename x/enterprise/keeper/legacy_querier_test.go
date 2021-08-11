package keeper_test

//	"github.com/unification-com/mainchain/x/enterprise/types"
//)
//
//const custom = "custom"
//
//func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.Params {
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryParameters}, "/"),
//		Data: []byte{},
//	}
//
//	bz, err := querier(ctx, []string{QueryParameters}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var params types.Params
//	require.NoError(t, cdc.UnmarshalJSON(bz, &params))
//
//	return params
//}
//
//func getQueriedTotalLocked(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) sdk.Coin {
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryTotalLocked}, "/"),
//		Data: []byte{},
//	}
//
//	bz, err := querier(ctx, []string{QueryTotalLocked}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var totalLocked sdk.Coin
//	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))
//
//	return totalLocked
//}
//
//func getQueriedTotalUnLocked(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) sdk.Coin {
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryTotalUnlocked}, "/"),
//		Data: []byte{},
//	}
//
//	bz, err := querier(ctx, []string{QueryTotalUnlocked}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var totalLocked sdk.Coin
//	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))
//
//	return totalLocked
//}
//
//func getQueriedTotalSupply(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) sdk.Coin {
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryTotalSupply}, "/"),
//		Data: []byte{},
//	}
//
//	bz, err := querier(ctx, []string{QueryTotalSupply}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var totalLocked sdk.Coin
//	require.NoError(t, cdc.UnmarshalJSON(bz, &totalLocked))
//
//	return totalLocked
//}
//
//func getQueriedPurchaseOrder(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, poID uint64) types.EnterpriseUndPurchaseOrder {
//
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, "/"),
//		Data: nil,
//	}
//
//	bz, err := querier(ctx, []string{QueryGetPurchaseOrder, strconv.FormatUint(poID, 10)}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var po types.EnterpriseUndPurchaseOrder
//	require.NoError(t, cdc.UnmarshalJSON(bz, &po))
//
//	return po
//}
//
//func getQueriedPurchaseOrders(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, page, limit int, purchaseOrderStatus types.PurchaseOrderStatus, purchaserAddr sdk.AccAddress) types.QueryResPurchaseOrders {
//
//	params := types.NewQueryPurchaseOrdersParams(page, limit, purchaseOrderStatus, purchaserAddr)
//
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryPurchaseOrders}, "/"),
//		Data: cdc.MustMarshalJSON(params),
//	}
//
//	bz, err := querier(ctx, []string{QueryPurchaseOrders}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var matchingOrders types.QueryResPurchaseOrders
//	require.NoError(t, cdc.UnmarshalJSON(bz, &matchingOrders))
//
//	return matchingOrders
//}
//
//func getQueriedWhitelist(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.WhitelistAddresses {
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryWhitelist}, "/"),
//		Data: []byte{},
//	}
//
//	bz, err := querier(ctx, []string{QueryWhitelist}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var whitelist types.WhitelistAddresses
//	require.NoError(t, cdc.UnmarshalJSON(bz, &whitelist))
//
//	return whitelist
//}
//
//func getQueriedIsAddressWhitelisted(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, addr sdk.AccAddress) bool {
//
//	query := abci.RequestQuery{
//		Path: strings.Join([]string{custom, types.QuerierRoute, QueryWhitelisted, addr.String()}, "/"),
//		Data: nil,
//	}
//
//	bz, err := querier(ctx, []string{QueryWhitelisted, addr.String()}, query)
//	require.NoError(t, err)
//	require.NotNil(t, bz)
//
//	var isWhitelisted bool
//	require.NoError(t, cdc.UnmarshalJSON(bz, &isWhitelisted))
//
//	return isWhitelisted
//}
//
//func TestQueryParams(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//	addresses := GenerateRandomAddresses(3)
//	entSigners := addresses[1].String() + "," + addresses[2].String()
//	paramsNew := types.NewParams("something", 1, 3600, entSigners)
//
//	keeper.SetParams(ctx, paramsNew)
//	params := getQueriedParams(t, ctx, keeper.cdc, querier)
//	require.True(t, ParamsEqual(paramsNew, params))
//}
//
//func TestQueryTotalLocked(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//
//	denom := TestDenomination
//	amount := int64(1000)
//	locked := sdk.NewInt64Coin(denom, amount)
//
//	_ = keeper.SetTotalLockedUnd(ctx, locked)
//
//	totalLocked := getQueriedTotalLocked(t, ctx, keeper.cdc, querier)
//	require.True(t, locked.IsEqual(totalLocked))
//}
//
//func TestQueryTotalUnLocked(t *testing.T) {
//	ctx, _, keeper, _, supplyKeeper := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//
//	denom := TestDenomination
//	amount := int64(1000)
//	locked := sdk.NewInt64Coin(denom, amount)
//
//	_ = keeper.SetTotalLockedUnd(ctx, locked)
//	totalSupply := supplyKeeper.GetSupply(ctx).GetTotal()
//
//	totalUnLocked := getQueriedTotalUnLocked(t, ctx, keeper.cdc, querier)
//
//	diff := totalSupply.Sub(sdk.Coins{totalUnLocked})
//
//	require.Equal(t, sdk.Coins{locked}, diff)
//}
//
//func TestQueryTotalSupply(t *testing.T) {
//	ctx, _, keeper, _, supplyKeeper := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//
//	totalSupply := supplyKeeper.GetSupply(ctx).GetTotal()
//
//	totalSupplyFromEnt := getQueriedTotalSupply(t, ctx, keeper.cdc, querier)
//	require.Equal(t, totalSupply, sdk.Coins{totalSupplyFromEnt})
//}
//
//func TestQueryPurchaseOrder(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//	testAddrs := GenerateRandomAddresses(100)
//	var testPos []types.EnterpriseUndPurchaseOrder
//
//	for i, addr := range testAddrs {
//		_ = keeper.AddAddressToWhitelist(ctx, addr)
//		poId := i + 1
//		po := types.NewEnterpriseUndPurchaseOrder()
//		po.Purchaser = addr
//		po.PurchaseOrderID = uint64(poId)
//		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i+1))
//		po.Status = RandomStatus()
//
//		err := keeper.SetPurchaseOrder(ctx, po)
//		require.NoError(t, err)
//		testPos = append(testPos, po)
//	}
//
//	for i, tPo := range testPos {
//		poId := i + 1
//		po := getQueriedPurchaseOrder(t, ctx, keeper.cdc, querier, uint64(poId))
//		require.True(t, PurchaseOrderEqual(tPo, po))
//	}
//}
//
//func TestQueryPurchaseOrdersFilters(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//	numTests := 100
//	testAddrs := GenerateRandomAddresses(numTests)
//
//	var testPosAccepted []types.EnterpriseUndPurchaseOrder
//	var testPosRejected []types.EnterpriseUndPurchaseOrder
//	var testPosRaised []types.EnterpriseUndPurchaseOrder
//	var testPosCompleted []types.EnterpriseUndPurchaseOrder
//
//	for i, addr := range testAddrs {
//		_ = keeper.AddAddressToWhitelist(ctx, addr)
//		status := RandomStatus()
//		poId := i + 1
//		po := types.NewEnterpriseUndPurchaseOrder()
//		po.Purchaser = addr
//		po.PurchaseOrderID = uint64(poId)
//		po.Amount = sdk.NewInt64Coin(types.DefaultDenomination, int64(i+1))
//		po.Status = status
//
//		err := keeper.SetPurchaseOrder(ctx, po)
//		require.NoError(t, err)
//
//		switch status {
//		case types.StatusRaised:
//			testPosRaised = append(testPosRaised, po)
//		case types.StatusCompleted:
//			testPosCompleted = append(testPosCompleted, po)
//		case types.StatusAccepted:
//			testPosAccepted = append(testPosAccepted, po)
//		case types.StatusRejected:
//			testPosRejected = append(testPosRejected, po)
//		}
//	}
//
//	queryRaised := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusRaised, sdk.AccAddress{})
//	require.True(t, len(queryRaised) == len(testPosRaised))
//
//	queryCompleted := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusCompleted, sdk.AccAddress{})
//	require.True(t, len(queryCompleted) == len(testPosCompleted))
//
//	queryAccepted := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusAccepted, sdk.AccAddress{})
//	require.True(t, len(queryAccepted) == len(testPosAccepted))
//
//	queryRejected := getQueriedPurchaseOrders(t, ctx, keeper.cdc, querier, 1, numTests, types.StatusRejected, sdk.AccAddress{})
//	require.True(t, len(queryRejected) == len(testPosRejected))
//
//	for _, tPo := range testPosRaised {
//		for _, po := range queryRaised {
//			if po.PurchaseOrderID == tPo.PurchaseOrderID {
//				require.True(t, PurchaseOrderEqual(tPo, po))
//			}
//		}
//	}
//
//	for _, tPo := range testPosCompleted {
//		for _, po := range queryCompleted {
//			if po.PurchaseOrderID == tPo.PurchaseOrderID {
//				require.True(t, PurchaseOrderEqual(tPo, po))
//			}
//		}
//	}
//
//	for _, tPo := range testPosAccepted {
//		for _, po := range queryAccepted {
//			if po.PurchaseOrderID == tPo.PurchaseOrderID {
//				require.True(t, PurchaseOrderEqual(tPo, po))
//			}
//		}
//	}
//
//	for _, tPo := range testPosRejected {
//		for _, po := range queryRejected {
//			if po.PurchaseOrderID == tPo.PurchaseOrderID {
//				require.True(t, PurchaseOrderEqual(tPo, po))
//			}
//		}
//	}
//}
//
//func TestQueryWhitelist(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//	numTests := 100
//	testAddrs := GenerateRandomAddresses(numTests)
//
//	for _, addr := range testAddrs {
//		_ = keeper.AddAddressToWhitelist(ctx, addr)
//	}
//
//	whitelist := getQueriedWhitelist(t, ctx, keeper.cdc, querier)
//
//	require.True(t, len(whitelist) == len(testAddrs))
//}
//
//func TestQueryAddressIsWhitelisted(t *testing.T) {
//	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
//	querier := NewLegacyQuerier(keeper)
//	numTests := 100
//	testAddrs := GenerateRandomAddresses(numTests)
//
//	for i, addr := range testAddrs {
//		if i < 50 {
//			_ = keeper.AddAddressToWhitelist(ctx, addr)
//		}
//	}
//
//	for i, addr := range testAddrs {
//		isWhiteListed := getQueriedIsAddressWhitelisted(t, ctx, keeper.cdc, querier, addr)
//		if i < 50 {
//			require.True(t, isWhiteListed)
//		} else {
//			require.False(t, isWhiteListed)
//		}
//	}
//}
