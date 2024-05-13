package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetParams(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := test_helpers.GenerateRandomTestAccounts(3)
	entSigners := testAddrs[1].String() + "," + testAddrs[2].String()
	denom := "testc"
	params := types.NewParams(denom, 1, 3600, entSigners)

	app.EnterpriseKeeper.SetParams(ctx, params)

	paramsDb := app.EnterpriseKeeper.GetParams(ctx)

	require.True(t, params == paramsDb)
	require.True(t, paramsDb.Denom == denom)
}
