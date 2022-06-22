package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/app/test_helpers"
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
	testApp := test_helpers.Setup(false)
	ctx := testApp.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(testApp, ctx)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, testApp.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, testApp.EnterpriseKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = testApp
	suite.ctx = ctx
	suite.queryClient = queryClient
	suite.addrs = test_helpers.AddTestAddrsIncremental(testApp, ctx, 10, sdk.NewInt(30000000))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
