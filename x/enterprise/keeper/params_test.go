package keeper_test

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/unification-com/mainchain/app"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetParams(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(3)
	entSigners := testAddrs[1].String() + "," + testAddrs[2].String()
	denom := "testc"
	params := types.NewParams(denom, 1, 3600, entSigners)

	app.EnterpriseKeeper.SetParams(ctx, params)

	paramsDb := app.EnterpriseKeeper.GetParams(ctx)

	require.True(t, params == paramsDb)
	require.True(t, paramsDb.Denom == denom)
}
