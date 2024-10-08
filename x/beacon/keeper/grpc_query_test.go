package keeper_test

import (
	gocontext "context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/beacon/types"
	"time"
)

func (s *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx, queryClient := s.app, s.ctx, s.queryClient

	testParams := types.Params{
		FeeRegister:         240,
		FeeRecord:           24,
		FeePurchaseStorage:  12,
		Denom:               "tnund",
		DefaultStorageLimit: 200,
		MaxStorageLimit:     300,
	}

	app.BeaconKeeper.SetParams(ctx, testParams)
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	s.NoError(err)
	s.Equal(testParams, paramsResp.Params)
}

func (s *KeeperTestSuite) TestGRPCQueryBeacon() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

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
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), bID)
				dbBeacon, found := app.BeaconKeeper.GetBeacon(ctx, uint64(1))
				s.Require().True(found)

				expBeacon = dbBeacon
			},
			true,
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			beaconRes, err := queryClient.Beacon(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(&expBeacon, beaconRes.Beacon)
			} else {
				s.Require().Error(err)
				s.Require().Nil(beaconRes)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryBeaconsFiltered() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

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
					expectedB.Moniker = simapp.GenerateRandomString(12)
					expectedB.Name = simapp.GenerateRandomString(24)
					expectedB.RegTime = uint64(ctx.BlockTime().Unix())

					bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
					s.Require().NoError(err)
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
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			beacons, err := queryClient.BeaconsFiltered(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)

				s.Require().Len(beacons.GetBeacons(), len(expRes.GetBeacons()))
				for i := 0; i < len(beacons.GetBeacons()); i++ {
					s.Require().Equal(beacons.GetBeacons()[i].String(), expRes.GetBeacons()[i].String())
				}
			} else {
				s.Require().Error(err)
				s.Require().Nil(beacons)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryBeaconTimestamp() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	var (
		req    *types.QueryBeaconTimestampRequest
		expRes types.QueryBeaconTimestampResponse
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
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), bID)

				expectedTs := types.BeaconTimestamp{
					Hash:        simapp.GenerateRandomString(32),
					SubmitTime:  uint64(time.Now().Unix()),
					TimestampId: 1,
				}

				err = app.BeaconKeeper.SetBeaconTimestamp(ctx, bID, expectedTs)
				s.Require().NoError(err)

				expRes = types.QueryBeaconTimestampResponse{
					Timestamp: &expectedTs,
					Owner:     addrs[0].String(),
					BeaconId:  bID,
				}

			},
			true,
		},
	}

	for _, testCase := range testCases {
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			timestampRes, err := queryClient.BeaconTimestamp(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(&expRes, timestampRes)
			} else {
				s.Require().Error(err)
				s.Require().Nil(timestampRes)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryBeaconStorage() {
	app, ctx, queryClient, addrs := s.app, s.ctx, s.queryClient, s.addrs

	var (
		req    *types.QueryBeaconStorageRequest
		expRes types.QueryBeaconStorageResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryBeaconStorageRequest{}
			},
			false,
		},
		{
			"zero beacon id request",
			func() {
				req = &types.QueryBeaconStorageRequest{BeaconId: 0}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryBeaconStorageRequest{BeaconId: 1}

				expectedB := types.Beacon{}
				expectedB.Owner = addrs[0].String()
				expectedB.LastTimestampId = 0
				expectedB.Moniker = "moniker"
				expectedB.Name = "name"

				bID, err := app.BeaconKeeper.RegisterNewBeacon(ctx, expectedB)
				s.Require().NoError(err)
				s.Require().Equal(uint64(1), bID)

				_, _, err = app.BeaconKeeper.RecordNewBeaconTimestamp(ctx, bID, "somehash", uint64(time.Now().Unix()))
				s.Require().NoError(err)

				expRes = types.QueryBeaconStorageResponse{
					BeaconId:       bID,
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
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			timestampRes, err := queryClient.BeaconStorage(gocontext.Background(), req)

			if testCase.expPass {
				s.Require().NoError(err)
				s.Require().Equal(&expRes, timestampRes)
			} else {
				s.Require().Error(err)
				s.Require().Nil(timestampRes)
			}
		})
	}
}
