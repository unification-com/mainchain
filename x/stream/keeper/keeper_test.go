package keeper_test

import (
	"testing"

	mathmod "cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/stream/keeper"
	"github.com/unification-com/mainchain/x/stream/types"
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
	s.addrs = simapp.AddTestAddrsIncremental(app, ctx, 100, mathmod.NewInt(1000000000000000000))
	s.msgServer = keeper.NewMsgServerImpl(s.app.StreamKeeper)
}

func (s *KeeperTestSuite) TestGetAuthority() {
	authority := s.app.StreamKeeper.GetAuthority()
	s.Equal("und10d07y265gmmuvt4z0w9aw880jnsr700ja85vs4", authority)
}
