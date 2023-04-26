package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
)

var testParams = types.NewParams(24, 2, 2, test_helpers.TestDenomination, 100, 200)

func TestSetGetParams(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	paramsDb := app.BeaconKeeper.GetParams(ctx)

	require.True(t, paramsDb.FeeRegister == testParams.FeeRegister)
	require.True(t, paramsDb.FeeRecord == testParams.FeeRecord)
	require.True(t, paramsDb.FeePurchaseStorage == testParams.FeePurchaseStorage)
	require.True(t, paramsDb.Denom == testParams.Denom)
	require.True(t, paramsDb.DefaultStorageLimit == testParams.DefaultStorageLimit)
	require.True(t, paramsDb.MaxStorageLimit == testParams.MaxStorageLimit)
}

func TestGetParamDenom(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamDenom(ctx)

	require.Equal(t, ret, testParams.Denom)
}

func TestGetParamRegistrationFee(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamRegistrationFee(ctx)

	require.Equal(t, ret, testParams.FeeRegister)
}

func TestGetParamRecordFee(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamRecordFee(ctx)

	require.Equal(t, ret, testParams.FeeRecord)
}

func TestGetParamPurchaseStorageFee(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamPurchaseStorageFee(ctx)

	require.Equal(t, ret, testParams.FeePurchaseStorage)
}

func TestGetParamDefaultStorageLimit(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamDefaultStorageLimit(ctx)

	require.Equal(t, ret, testParams.DefaultStorageLimit)
}

func TestGetParamMaxStorageLimit(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetParamMaxStorageLimit(ctx)

	require.Equal(t, ret, testParams.MaxStorageLimit)
}

func TestGetZeroFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetZeroFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, 0)

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetRegistrationFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetRecordFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetPurchaseStorageFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetPurchaseStorageFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeePurchaseStorage))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetZeroFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetZeroFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, 0)}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetRegistrationFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetRecordFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetPurchaseStorageFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.BeaconKeeper.SetParams(ctx, testParams)

	ret := app.BeaconKeeper.GetPurchaseStorageFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeePurchaseStorage))}

	require.True(t, ret.IsEqual(paramCoin))
}
