package keeper_test

import (
	simapp "github.com/unification-com/mainchain/app"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *app.App
	ctx         sdk.Context
	queryClient types.QueryClient
	addrs       []sdk.AccAddress
}

func (suite *KeeperTestSuite) SetupTest() {
	testApp := simapp.Setup(suite.T(), false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(testApp, ctx)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, testApp.EnterpriseKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = testApp
	suite.ctx = ctx
	suite.queryClient = queryClient
	suite.addrs = simapp.AddTestAddrsIncremental(testApp, ctx, 10, sdk.NewInt(30000000))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
