package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// Tests for Highest Purchase Order ID

func TestSetGetParams(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	newFeeReg := uint64(1000)
	newFeeRec := uint64(100)
	denom := "somecoin"
	params := types.NewParams(newFeeReg, newFeeRec, denom)

	keeper.SetParams(ctx, params)

	paramsDb := keeper.GetParams(ctx)

	require.True(t, ParamsEqual(params, paramsDb))
	require.True(t, paramsDb.FeeRegister == newFeeReg)
	require.True(t, paramsDb.FeeRecord == newFeeRec)
	require.True(t, paramsDb.Denom == denom)
}
