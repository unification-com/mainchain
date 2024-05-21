package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/enterprise/keeper"
	"testing"

	"github.com/unification-com/mainchain/app"
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
	testApp := simapp.Setup(s.T(), false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(testApp, ctx)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, testApp.EnterpriseKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = testApp
	s.ctx = ctx
	s.queryClient = queryClient
	s.addrs = simapp.AddTestAddrsIncremental(testApp, ctx, 10, sdk.NewInt(30000000))
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
				Denom:             "test",
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
