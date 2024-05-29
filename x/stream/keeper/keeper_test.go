package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/stream/keeper"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *simapp.App
	ctx         sdk.Context
	queryClient types.QueryClient
	addrs       []sdk.AccAddress
	msgServer   types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(s.T(), true)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.StreamKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = app
	s.ctx = ctx
	s.queryClient = queryClient
	s.addrs = simapp.AddTestAddrsIncremental(app, ctx, 10, sdk.NewInt(30000000))
	s.msgServer = keeper.NewMsgServerImpl(s.app.StreamKeeper)
}
func (s *KeeperTestSuite) TestParams() {
	invalidFee, _ := sdk.NewDecFromStr("1.01")
	validFee, _ := sdk.NewDecFromStr("0.24")

	testCases := []struct {
		name      string
		input     types.Params
		expectErr bool
	}{
		{
			name: "set invalid params",
			input: types.Params{
				ValidatorFee: invalidFee,
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			input: types.Params{
				ValidatorFee: validFee,
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			expected := s.app.StreamKeeper.GetParams(s.ctx)
			err := s.app.StreamKeeper.SetParams(s.ctx, tc.input)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				expected = tc.input
				s.Require().NoError(err)
			}

			p := s.app.StreamKeeper.GetParams(s.ctx)
			s.Require().Equal(expected, p)
		})
	}
}
