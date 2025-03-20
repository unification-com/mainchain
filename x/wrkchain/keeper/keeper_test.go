package keeper_test

import (
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"testing"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/wrkchain/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	queryClient types.QueryClient
	addrs       []sdk.AccAddress
	msgServer   types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := simapphelpers.Setup(s.T())
	ctx := app.BaseApp.NewContext(false)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.WrkchainKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = app
	s.ctx = ctx
	s.queryClient = queryClient
	s.addrs = simapphelpers.AddTestAddrsIncremental(app, ctx, 10, mathmod.NewInt(30000000))
	s.msgServer = keeper.NewMsgServerImpl(s.app.WrkchainKeeper)
}

func (s *KeeperTestSuite) TestParams() {
	testCases := []struct {
		name      string
		input     types.Params
		expectErr bool
	}{
		{
			name: "set invalid params",
			input: types.Params{
				FeeRegister:         0,
				FeeRecord:           0,
				FeePurchaseStorage:  0,
				Denom:               "",
				DefaultStorageLimit: 0,
				MaxStorageLimit:     0,
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			input: types.Params{
				FeeRegister:         24,
				FeeRecord:           2,
				FeePurchaseStorage:  24,
				Denom:               "test",
				DefaultStorageLimit: 99,
				MaxStorageLimit:     999,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			expected := s.app.WrkchainKeeper.GetParams(s.ctx)
			err := s.app.WrkchainKeeper.SetParams(s.ctx, tc.input)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				expected = tc.input
				s.Require().NoError(err)
			}

			p := s.app.WrkchainKeeper.GetParams(s.ctx)
			s.Require().Equal(expected, p)
		})
	}
}
