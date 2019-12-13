package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Tests for Highest Purchase Order ID

func TestSetGetParams(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	entSrc := TestAddrs[2]
	denom := "testc"
	params := types.NewParams(entSrc, denom)

	keeper.SetParams(ctx, params)

	paramsDb := keeper.GetParams(ctx)

	require.True(t, ParamsEqual(params, paramsDb))
	require.True(t, paramsDb.EntSource.String() == entSrc.String())
	require.True(t, paramsDb.Denom == denom)
}
