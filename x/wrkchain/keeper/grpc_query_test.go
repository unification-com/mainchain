package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"time"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	testParams := types.Params{
		FeeRegister:         240,
		FeeRecord:           24,
		FeePurchaseStorage:  12,
		Denom:               "tnund",
		DefaultStorageLimit: 200,
		MaxStorageLimit:     300,
	}

	app.WrkchainKeeper.SetParams(ctx, testParams)
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	suite.NoError(err)
	suite.Equal(testParams, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryWrkChain() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req   *types.QueryWrkChainRequest
		expWc types.WrkChain
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryWrkChainRequest{}
			},
			false,
		},
		{
			"non existing wrkchain request",
			func() {
				req = &types.QueryWrkChainRequest{WrkchainId: 3}
			},
			false,
		},
		{
			"zero wrkchain id request",
			func() {
				req = &types.QueryWrkChainRequest{WrkchainId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {

				req = &types.QueryWrkChainRequest{WrkchainId: 1}

				bID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "lhbohbob", "tm", addrs[0])
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(1), bID)
				dbWrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, uint64(1))
				suite.Require().True(found)

				expWc = dbWrkchain
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			wcRes, err := queryClient.WrkChain(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expWc, wcRes.Wrkchain)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(wcRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryWrkChainsFiltered() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testWrkchains := []types.WrkChain{}

	var (
		req    *types.QueryWrkChainsFilteredRequest
		expRes *types.QueryWrkChainsFilteredResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryWrkChainsFilteredRequest{}
			},
			true,
		},
		{
			"request wrkchains with limit 3",
			func() {
				// create 5 test wrkchains
				for i := 0; i < 5; i++ {

					moniker := simapp.GenerateRandomString(12)
					name := simapp.GenerateRandomString(24)
					gHash := simapp.GenerateRandomString(64)

					wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, gHash, "tm", addrs[0])
					suite.Require().NoError(err)
					expectedWc, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
					suite.Require().True(found)
					testWrkchains = append(testWrkchains, expectedWc)
				}

				req = &types.QueryWrkChainsFilteredRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &types.QueryWrkChainsFilteredResponse{
					Wrkchains: testWrkchains[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				req = &types.QueryWrkChainsFilteredRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &types.QueryWrkChainsFilteredResponse{
					Wrkchains: testWrkchains[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &types.QueryWrkChainsFilteredRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &types.QueryWrkChainsFilteredResponse{
					Wrkchains: testWrkchains[:2],
				}
			},
			true,
		},
		{
			"request with moniker filter",
			func() {
				req = &types.QueryWrkChainsFilteredRequest{
					Moniker: testWrkchains[0].Moniker,
				}

				expRes = &types.QueryWrkChainsFilteredResponse{
					Wrkchains: testWrkchains[:1],
				}
			},
			true,
		},
		{
			"request with owner filter",
			func() {
				req = &types.QueryWrkChainsFilteredRequest{
					Owner: testWrkchains[0].Owner,
				}

				expRes = &types.QueryWrkChainsFilteredResponse{
					Wrkchains: testWrkchains,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			wrkchains, err := queryClient.WrkChainsFiltered(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(wrkchains.GetWrkchains(), len(expRes.GetWrkchains()))
				for i := 0; i < len(wrkchains.GetWrkchains()); i++ {
					suite.Require().Equal(wrkchains.GetWrkchains()[i].String(), expRes.GetWrkchains()[i].String())
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(wrkchains)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryWrkchainBlock() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req    *types.QueryWrkChainBlockRequest
		expRes types.QueryWrkChainBlockResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryWrkChainBlockRequest{}
			},
			false,
		},
		{
			"zero wrkchain id request",
			func() {
				req = &types.QueryWrkChainBlockRequest{WrkchainId: 0}
			},
			false,
		},
		{
			"zero block height request",
			func() {
				req = &types.QueryWrkChainBlockRequest{WrkchainId: 1, Height: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryWrkChainBlockRequest{WrkchainId: 1, Height: 1}

				wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", addrs[0])
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(1), wcID)

				expectedBlock := types.WrkChainBlock{
					Blockhash: simapp.GenerateRandomString(32),
					SubTime:   uint64(time.Now().Unix()),
					Height:    1,
				}

				err = app.WrkchainKeeper.SetWrkChainBlock(ctx, wcID, expectedBlock)
				suite.Require().NoError(err)

				expRes = types.QueryWrkChainBlockResponse{
					Block:      &expectedBlock,
					WrkchainId: wcID,
					Owner:      addrs[0].String(),
				}

			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			blockRes, err := queryClient.WrkChainBlock(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expRes, blockRes)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(blockRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryWrkChainStorage() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req    *types.QueryWrkChainStorageRequest
		expRes types.QueryWrkChainStorageResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryWrkChainStorageRequest{}
			},
			false,
		},
		{
			"zero wrkchain id request",
			func() {
				req = &types.QueryWrkChainStorageRequest{WrkchainId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryWrkChainStorageRequest{WrkchainId: 1}

				wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", addrs[0])
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(1), wcId)

				_, err = app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcId, 24, "somehash", "parenthash", "hash1", "hash2", "hash3")
				suite.Require().NoError(err)

				expRes = types.QueryWrkChainStorageResponse{
					WrkchainId:     wcId,
					Owner:          addrs[0].String(),
					CurrentLimit:   types.DefaultStorageLimit,
					CurrentUsed:    1,
					Max:            types.DefaultMaxStorageLimit,
					MaxPurchasable: types.DefaultMaxStorageLimit - types.DefaultStorageLimit,
				}

			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			timestampRes, err := queryClient.WrkChainStorage(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expRes, timestampRes)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(timestampRes)
			}
		})
	}
}
