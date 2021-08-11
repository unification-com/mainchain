package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)


var testParams = types.NewParams(24, 2, test_helpers.TestDenomination)

func TestSetGetParams(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	paramsDb := app.WrkchainKeeper.GetParams(ctx)

	require.True(t, paramsDb.FeeRegister == testParams.FeeRegister)
	require.True(t, paramsDb.FeeRecord == testParams.FeeRecord)
	require.True(t, paramsDb.Denom == testParams.Denom)
}

func TestGetParamDenom(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamDenom(ctx)

	require.Equal(t, ret, testParams.Denom)
}

func TestGetParamRegistrationFee(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamRegistrationFee(ctx)

	require.Equal(t, ret, testParams.FeeRegister)
}

func TestGetParamRecordFee(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetParamRecordFee(ctx)

	require.Equal(t, ret, testParams.FeeRecord)
}

func TestGetZeroFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetZeroFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, 0)

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRegistrationFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoin(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRecordFeeAsCoin(ctx)

	paramCoin := sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetZeroFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetZeroFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, 0)}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRegistrationFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRegistrationFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRegister))}

	require.True(t, ret.IsEqual(paramCoin))
}

func TestGetRecordFeeAsCoins(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.WrkchainKeeper.SetParams(ctx, testParams)

	ret := app.WrkchainKeeper.GetRecordFeeAsCoins(ctx)

	paramCoin := sdk.Coins{sdk.NewInt64Coin(testParams.Denom, int64(testParams.FeeRecord))}

	require.True(t, ret.IsEqual(paramCoin))
}
