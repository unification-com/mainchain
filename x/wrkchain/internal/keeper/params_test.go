package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

var testParams = types.NewParams(1000, 100, "nund")

func TestSetGetParams(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	keeper.SetParams(ctx, testParams)

	paramsDb := keeper.GetParams(ctx)

	require.True(t, ParamsEqual(testParams, paramsDb))
	require.True(t, paramsDb.FeeRegister == testParams.FeeRegister)
	require.True(t, paramsDb.FeeRecord == testParams.FeeRecord)
	require.True(t, paramsDb.Denom == testParams.Denom)
}

func TestGetParamDenom(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetParamDenom(ctx)

	require.Equal(t, ret, testParams.Denom)
}

func TestGetParamRegistrationFee(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetParamRegistrationFee(ctx)

	require.Equal(t, ret, testParams.FeeRegister)
}

func TestGetParamRecordFee(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetParamRecordFee(ctx)

	require.Equal(t, ret, testParams.FeeRecord)
}

func TestGetZeroFeeAsCoin(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetZeroFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, 0)

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoin(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetRegistrationFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoin(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetRecordFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetZeroFeeAsCoins(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetZeroFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, 0)}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoins(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetRegistrationFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoins(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	keeper.SetParams(ctx, testParams)

	ret := keeper.GetRecordFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))}

	require.True(t, ret.IsEqual(paramCoin))
}
