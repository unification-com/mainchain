package keeper_test

import (
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

var testParams = types.NewParams(24, 2, 2, sdk.DefaultBondDenom, 100, 200)

func TestSetGetParams(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	paramsDb := app.WrkchainKeeper.GetParams(ctx)

	require.True(t, paramsDb.FeeRegister == testParams.FeeRegister)
	require.True(t, paramsDb.FeeRecord == testParams.FeeRecord)
	require.True(t, paramsDb.FeePurchaseStorage == testParams.FeePurchaseStorage)
	require.True(t, paramsDb.Denom == testParams.Denom)
	require.True(t, paramsDb.DefaultStorageLimit == testParams.DefaultStorageLimit)
	require.True(t, paramsDb.MaxStorageLimit == testParams.MaxStorageLimit)
}

func TestGetParamDenom(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamDenom(ctx)

	require.Equal(t, ret, testParams.Denom)
}

func TestGetParamRegistrationFee(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamRegistrationFee(ctx)

	require.Equal(t, ret, testParams.FeeRegister)
}

func TestGetParamRecordFee(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamRecordFee(ctx)

	require.Equal(t, ret, testParams.FeeRecord)
}

func TestGetParamPurchaseStorageFee(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamPurchaseStorageFee(ctx)

	require.Equal(t, ret, testParams.FeePurchaseStorage)
}

func TestGetParamDefaultStorageLimit(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamDefaultStorageLimit(ctx)

	require.Equal(t, ret, testParams.DefaultStorageLimit)
}

func TestGetParamMaxStorageLimit(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamMaxStorageLimit(ctx)

	require.Equal(t, ret, testParams.MaxStorageLimit)
}

func TestGetZeroFeeAsCoin(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetZeroFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, 0)

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoin(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRegistrationFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoin(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRecordFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetPurchaseStorageFeeAsCoin(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetPurchaseStorageFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeePurchaseStorage))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetZeroFeeAsCoins(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetZeroFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, 0)}

	require.True(t, ret.Equal(paramCoin))
}

func TestGetRegistrationFeeAsCoins(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRegistrationFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))}

	require.True(t, ret.Equal(paramCoin))
}

func TestGetRecordFeeAsCoins(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRecordFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))}

	require.True(t, ret.Equal(paramCoin))
}

func TestGetPurchaseStorageFeeAsCoins(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetPurchaseStorageFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeePurchaseStorage))}

	require.True(t, ret.Equal(paramCoin))
}
