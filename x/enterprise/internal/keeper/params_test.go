package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

// Tests for Highest Purchase Order ID

func TestSetGetParams(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	var entSigners []sdk.AccAddress
	entSigners = append(entSigners, TestAddrs[1])
	entSigners = append(entSigners, TestAddrs[2])
	denom := "testc"
	params := types.NewParams(entSigners, denom, 1, 3600)

	keeper.SetParams(ctx, params)

	paramsDb := keeper.GetParams(ctx)

	require.True(t, ParamsEqual(params, paramsDb))
	require.True(t, paramsDb.Denom == denom)
}
