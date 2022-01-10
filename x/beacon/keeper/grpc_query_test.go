package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
	"time"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	testParams := types.Params{
		FeeRegister: 240,
		FeeRecord:   24,
		Denom:       "tnund",
	}

	app.BeaconKeeper.SetParams(ctx, testParams)
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	suite.NoError(err)
	suite.Equal(testParams, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryBeacon() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req       *types.QueryBeaconRequest
		expBeacon types.Beacon
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryBeaconRequest{}
			},
			false,
		},
		{
			"non existing beacon request",
			func() {
				req = &types.QueryBeaconRequest{BeaconId: 3}
			},
			false,
		},
		{
			"zero beacon id request",
			func() {
				req = &types.QueryBeaconRequest{BeaconId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {

				req = &types.QueryBeaconRequest{BeaconId: 1}

				expectedB := types.Beacon{}
				expectedB.Owner = addrs[0].String()
				expectedB.LastTimestampId = 0
				expectedB.Moniker = "moniker"
				expectedB.Name = "name"

				bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(1), bID)
				dbBeacon, found := app.BeaconKeeper.GetBeacon(ctx, uint64(1))
				suite.Require().True(found)

				expBeacon = dbBeacon
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			beaconRes, err := queryClient.Beacon(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expBeacon, beaconRes.Beacon)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(beaconRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryBeaconsFiltered() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testBeacons := []types.Beacon{}

	var (
		req    *types.QueryBeaconsFilteredRequest
		expRes *types.QueryBeaconsFilteredResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryBeaconsFilteredRequest{}
			},
			true,
		},
		{
			"request beacons with limit 3",
			func() {
				// create 5 test beacons
				for i := 0; i < 5; i++ {
					expectedB := types.Beacon{}
					expectedB.Owner = addrs[0].String()
					expectedB.LastTimestampId = 0
					expectedB.Moniker = test_helpers.GenerateRandomString(12)
					expectedB.Name = test_helpers.GenerateRandomString(24)

					bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
					suite.Require().NoError(err)
					expectedB.BeaconId = bID
					testBeacons = append(testBeacons, expectedB)
				}

				req = &types.QueryBeaconsFilteredRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &types.QueryBeaconsFilteredResponse{
					Beacons: testBeacons[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				req = &types.QueryBeaconsFilteredRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &types.QueryBeaconsFilteredResponse{
					Beacons: testBeacons[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &types.QueryBeaconsFilteredRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &types.QueryBeaconsFilteredResponse{
					Beacons: testBeacons[:2],
				}
			},
			true,
		},
		{
			"request with moniker filter",
			func() {
				req = &types.QueryBeaconsFilteredRequest{
					Moniker: testBeacons[0].Moniker,
				}

				expRes = &types.QueryBeaconsFilteredResponse{
					Beacons: testBeacons[:1],
				}
			},
			true,
		},
		{
			"request with owner filter",
			func() {
				req = &types.QueryBeaconsFilteredRequest{
					Owner: testBeacons[0].Owner,
				}

				expRes = &types.QueryBeaconsFilteredResponse{
					Beacons: testBeacons,
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			beacons, err := queryClient.BeaconsFiltered(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(beacons.GetBeacons(), len(expRes.GetBeacons()))
				for i := 0; i < len(beacons.GetBeacons()); i++ {
					suite.Require().Equal(beacons.GetBeacons()[i].String(), expRes.GetBeacons()[i].String())
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(beacons)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryBeaconTimestamp() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req    *types.QueryBeaconTimestampRequest
		expRes types.BeaconTimestamp
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryBeaconTimestampRequest{}
			},
			false,
		},
		{
			"zero beacon id request",
			func() {
				req = &types.QueryBeaconTimestampRequest{BeaconId: 0}
			},
			false,
		},
		{
			"zero timestamp id request",
			func() {
				req = &types.QueryBeaconTimestampRequest{BeaconId: 1, TimestampId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryBeaconTimestampRequest{BeaconId: 1, TimestampId: 1}

				expectedB := types.Beacon{}
				expectedB.Owner = addrs[0].String()
				expectedB.LastTimestampId = 0
				expectedB.Moniker = "moniker"
				expectedB.Name = "name"

				bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
				suite.Require().NoError(err)
				suite.Require().Equal(uint64(1), bID)

				expectedTs := types.BeaconTimestamp{
					BeaconId:    bID,
					Hash:        test_helpers.GenerateRandomString(32),
					SubmitTime:  uint64(time.Now().Unix()),
					Owner:       addrs[0].String(),
					TimestampId: 1,
				}

				err = app.BeaconKeeper.SetBeaconTimestamp(ctx, expectedTs)
				suite.Require().NoError(err)

				expRes = expectedTs

			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			timestampRes, err := queryClient.BeaconTimestamp(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expRes, timestampRes.Timestamp)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(timestampRes)
			}
		})
	}
}
