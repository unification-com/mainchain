package keeper_test

import (
	gocontext "context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"time"
)

func (s *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	testParams := types.Params{
		EntSigners:        addrs[0].String(),
		Denom:             sdk.DefaultBondDenom,
		MinAccepts:        1,
		DecisionTimeLimit: 600,
	}

	app.EnterpriseKeeper.SetParams(ctx, testParams)
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	s.NoError(err)
	s.Equal(testParams, paramsResp.Params)
}

func (s *KeeperTestSuite) TestGRPCQueryEnterpriseUndPurchaseOrder() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	var (
		req   *types.QueryEnterpriseUndPurchaseOrderRequest
		expPo types.EnterpriseUndPurchaseOrder
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"po id must not be 0",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{
					PurchaseOrderId: 0,
				}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"po does not exist",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{
					PurchaseOrderId: 99,
				}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"valid request",
			func() {

				poId := uint64(1)
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{PurchaseOrderId: poId}

				expectedPo := types.EnterpriseUndPurchaseOrder{
					Id:             poId,
					Purchaser:      addrs[0].String(),
					Amount:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 100),
					Status:         types.StatusRaised,
					RaiseTime:      uint64(time.Now().Unix()),
					CompletionTime: 0,
					Decisions:      nil,
				}

				err := app.EnterpriseKeeper.SetPurchaseOrder(ctx, expectedPo)
				s.Require().Nil(err)
				dbPo, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
				s.Require().True(found)

				expPo = dbPo
			},
			true,
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			poRes, err := queryClient.EnterpriseUndPurchaseOrder(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(expPo, poRes.PurchaseOrder)
			} else {
				s.Require().Error(err)
				s.Require().Nil(poRes)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryEnterpriseUndPurchaseOrders() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	testPos := []types.EnterpriseUndPurchaseOrder{}

	var (
		req    *types.QueryEnterpriseUndPurchaseOrdersRequest
		expRes *types.QueryEnterpriseUndPurchaseOrdersResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{}
			},
			true,
		},
		{
			"request pos with limit 3",
			func() {
				// create test pos
				for i := 0; i < len(addrs); i++ {
					newPo := types.EnterpriseUndPurchaseOrder{
						Purchaser: addrs[i].String(),
						Amount:    sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(i)+1),
					}

					poId, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, newPo)
					s.Require().NoError(err)
					expectedPo, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
					s.Require().True(found)
					testPos = append(testPos, expectedPo)
				}

				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 3",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[3:6],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:2],
				}
			},
			true,
		},
		{
			"request with purchaser filter",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Purchaser:  addrs[0].String(),
					Pagination: &query.PageRequest{Limit: 2},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:1],
				}
			},
			true,
		},
		{
			"request with status filter",
			func() {

				expectedPo, _ := app.EnterpriseKeeper.GetPurchaseOrder(ctx, testPos[0].Id)
				expectedPo.Status = types.StatusCompleted
				_ = app.EnterpriseKeeper.SetPurchaseOrder(ctx, expectedPo)
				testPos[0] = expectedPo

				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Status:     types.StatusCompleted,
					Pagination: &query.PageRequest{Limit: 2},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:1],
				}
			},
			true,
		},
		{
			"request with purchaser and status filters",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Purchaser:  addrs[0].String(),
					Status:     types.StatusCompleted,
					Pagination: &query.PageRequest{Limit: 2},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:1],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			pos, err := queryClient.EnterpriseUndPurchaseOrders(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)

				s.Require().Len(pos.GetPurchaseOrders(), len(expRes.GetPurchaseOrders()))
				for i := 0; i < len(pos.GetPurchaseOrders()); i++ {
					s.Require().Equal(pos.GetPurchaseOrders()[i].String(), expRes.GetPurchaseOrders()[i].String())
				}
			} else {
				s.Require().Error(err)
				s.Require().Nil(pos)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryLockedUndByAddress() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	var (
		req    *types.QueryLockedUndByAddressRequest
		expRes types.QueryLockedUndByAddressResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryLockedUndByAddressRequest{}
			},
			false,
		},
		{
			"invalid address",
			func() {
				req = &types.QueryLockedUndByAddressRequest{
					Owner: "rubbish",
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryLockedUndByAddressRequest{
					Owner: addrs[0].String(),
				}

				l := types.LockedUnd{
					Owner:  addrs[0].String(),
					Amount: sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000),
				}

				err := app.EnterpriseKeeper.SetLockedUndForAccount(ctx, l)
				s.Require().NoError(err)
				expRes = types.QueryLockedUndByAddressResponse{
					Amount: l.Amount,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			lRes, err := queryClient.LockedUndByAddress(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(&expRes, lRes)
			} else {
				s.Require().Error(err)
				s.Require().Nil(lRes)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryTotalLocked() {
	app, ctx, queryClient := s.app, s.ctx, s.queryClient

	req := &types.QueryTotalLockedRequest{}
	locked := sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	expectedRes := &types.QueryTotalLockedResponse{
		Amount: locked,
	}

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	s.Require().NoError(err)

	lRes, err := queryClient.TotalLocked(gocontext.Background(), req)

	s.Require().NoError(err)
	s.Require().Equal(expectedRes, lRes)
}

func (s *KeeperTestSuite) TestGRPCQueryWhitelist() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	var whitelistedAddrs []string

	for _, addr := range addrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		s.Require().NoError(err)
		whitelistedAddrs = append(whitelistedAddrs, addr.String())
	}

	req := &types.QueryWhitelistRequest{}
	res, err := queryClient.Whitelist(gocontext.Background(), req)
	s.Require().NoError(err)
	s.Require().Equal(whitelistedAddrs, res.Addresses)
}

func (s *KeeperTestSuite) TestGRPCQueryWhitelisted() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs
	notListed := simapphelpers.GenerateRandomTestAccounts(10)

	for _, addr := range addrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		s.Require().NoError(err)
	}

	for _, addr := range addrs {
		req := &types.QueryWhitelistedRequest{Address: addr.String()}
		res, err := queryClient.Whitelisted(gocontext.Background(), req)
		s.Require().NoError(err)
		s.Require().True(res.Whitelisted)
		s.Require().Equal(addr.String(), res.Address)
	}

	for _, addr := range notListed {
		req := &types.QueryWhitelistedRequest{Address: addr.String()}
		res, err := queryClient.Whitelisted(gocontext.Background(), req)
		s.Require().NoError(err)
		s.Require().False(res.Whitelisted)
		s.Require().Equal(addr.String(), res.Address)
	}
}

func (s *KeeperTestSuite) TestTotalSpentEFUND() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	poAmount := uint64(12345)
	totalUnlocked := uint64(0)
	toUnlock := uint64(123)
	poAmountCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(poAmount))
	toUnlockCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(toUnlock))

	for i := 0; i < len(addrs); i++ {
		newPo := types.EnterpriseUndPurchaseOrder{
			Purchaser: addrs[i].String(),
			Amount:    poAmountCoin,
		}

		poId, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, newPo)
		s.Require().NoError(err)
		expectedPo, _ := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
		expectedPo.Status = types.StatusCompleted
		_ = app.EnterpriseKeeper.SetPurchaseOrder(ctx, expectedPo)

		err = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addrs[i], poAmountCoin)
		s.Require().NoError(err)

		err = app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addrs[i], sdk.Coins{toUnlockCoin})
		s.Require().NoError(err)

		totalUnlocked += toUnlock
	}

	expectedResp := &types.QueryTotalSpentEFUNDResponse{Amount: sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(totalUnlocked))}

	req := &types.QueryTotalSpentEFUNDRequest{}
	res, err := queryClient.TotalSpentEFUND(gocontext.Background(), req)
	s.Require().NoError(err)
	s.Require().Equal(expectedResp, res)
}

func (s *KeeperTestSuite) TestSpentEFUNDByAddress() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	for i := 0; i < len(addrs); i++ {
		poAmount := uint64(i+1) * 10
		toUnlock := uint64(i + 1)
		poAmountCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(poAmount))
		toUnlockCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(toUnlock))
		newPo := types.EnterpriseUndPurchaseOrder{
			Purchaser: addrs[i].String(),
			Amount:    poAmountCoin,
		}

		poId, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, newPo)
		s.Require().NoError(err)
		expectedPo, _ := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
		expectedPo.Status = types.StatusCompleted
		_ = app.EnterpriseKeeper.SetPurchaseOrder(ctx, expectedPo)

		err = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addrs[i], poAmountCoin)
		s.Require().NoError(err)

		err = app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addrs[i], sdk.Coins{toUnlockCoin})
		s.Require().NoError(err)

		expectedResp := &types.QuerySpentEFUNDByAddressResponse{Amount: sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(toUnlock))}

		req := &types.QuerySpentEFUNDByAddressRequest{
			Address: addrs[i].String(),
		}
		res, err := queryClient.SpentEFUNDByAddress(gocontext.Background(), req)
		s.Require().NoError(err)
		s.Require().Equal(expectedResp.Amount.String(), res.Amount.String())
	}
}
