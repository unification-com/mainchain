package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetParams(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(3)
	entSigners := testAddrs[1].String() + "," + testAddrs[2].String()
	denom := "testc"
	params := types.NewParams(denom, 1, 3600, entSigners)

	app.EnterpriseKeeper.SetParams(ctx, params)

	paramsDb := app.EnterpriseKeeper.GetParams(ctx)

	require.True(t, params == paramsDb)
	require.True(t, paramsDb.Denom == denom)
}
