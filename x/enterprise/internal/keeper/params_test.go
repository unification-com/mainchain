package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Tests for Highest Purchase Order ID

func TestSetGetParams(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	entSigners := TestAddrs[1].String() + "," + TestAddrs[2].String()
	denom := "testc"
	params := types.NewParams(denom, 1, 3600, entSigners)

	keeper.SetParams(ctx, params)

	paramsDb := keeper.GetParams(ctx)

	require.True(t, ParamsEqual(params, paramsDb))
	require.True(t, paramsDb.Denom == denom)
}
