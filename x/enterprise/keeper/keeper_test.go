package keeper_test

import (
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"testing"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"github.com/unification-com/mainchain/x/enterprise/types"
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
	testApp := simapphelpers.Setup(s.T())
	ctx := testApp.BaseApp.NewContext(false)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, testApp.EnterpriseKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = testApp
	s.ctx = ctx
	s.queryClient = queryClient
	s.addrs = simapphelpers.AddTestAddrsIncremental(testApp, ctx, 10, mathmod.NewInt(30000000))
	s.msgServer = keeper.NewMsgServerImpl(s.app.EnterpriseKeeper)
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
				EntSigners:        "",
				Denom:             "",
				MinAccepts:        0,
				DecisionTimeLimit: 0,
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			input: types.Params{
				EntSigners:        "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6,und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed",
				Denom:             sdk.DefaultBondDenom,
				MinAccepts:        2,
				DecisionTimeLimit: 2,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			expected := s.app.EnterpriseKeeper.GetParams(s.ctx)
			err := s.app.EnterpriseKeeper.SetParams(s.ctx, tc.input)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				expected = tc.input
				s.Require().NoError(err)
			}

			p := s.app.EnterpriseKeeper.GetParams(s.ctx)
			s.Require().Equal(expected, p)
		})
	}
}
