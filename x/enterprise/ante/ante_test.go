package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/ante"
	"github.com/unification-com/mainchain/x/enterprise/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

const TestChainID = "und-unit-test-chain"

func fundAccount(ctx sdk.Context, bk bankkeeper.Keeper, addr sdk.AccAddress, amtCoins sdk.Coins) error {
	err := bk.MintCoins(ctx, types.ModuleName, amtCoins)
	if err != nil {
		return err
	}
	err = bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amtCoins)
	if err != nil {
		return err
	}
	return nil
}

func TestCheckLockedUndDecoratorModuleAndSupplyInsufficientFunds(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// artificially add locked FUND without minting first
	toLock := sdk.NewCoin(actualFeeDenom, accAmt)
	lockeUnd := types.LockedUnd{
		Owner:  addr.String(),
		Amount: toLock,
	}

	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockeUnd)

	feeInt := int64(1)
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

}

func TestCheckLockedUndDecoratorSuccessfulUnlock(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	accAmt := sdk.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt)))

	feeInt := int64(1)
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)
}

func TestCheckLockedUndDecoratorSkipIfNothingLocked(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	feeInt := int64(actualRegFeeAmt)
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)
}
