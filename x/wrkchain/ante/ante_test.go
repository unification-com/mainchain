package ante_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	"github.com/unification-com/mainchain/x/wrkchain/ante"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"testing"
)

const TestChainID = "und-unit-test-chain"

func fundAccount(ctx sdk.Context, bk bankkeeper.Keeper, addr sdk.AccAddress, amtCoins sdk.Coins) error {
	err := bk.MintCoins(ctx, enttypes.ModuleName, amtCoins)
	if err != nil {
		return err
	}
	err = bk.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, addr, amtCoins)
	if err != nil {
		return err
	}
	return nil
}

func TestCorrectWrkChainFeeDecoratorAddressNotExist(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, test_helpers.TestDenomination))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fee payer does not exist
	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, int64(actualFeeAmt)))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", addr)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectWrkChainFeeDecoratorRejectTooLittleFeeInTx(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, test_helpers.TestDenomination))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// Register
	feeInt := int64(actualRegFeeAmt - 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr := sdkerrors.Wrap(types.ErrInsufficientWrkChainFee, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt - 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr1 := sdkerrors.Wrap(types.ErrInsufficientWrkChainFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorRejectTooMuchFeeInTx(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, test_helpers.TestDenomination))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt + 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr := sdkerrors.Wrap(types.ErrTooMuchWrkChainFee, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt + 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr1 := sdkerrors.Wrap(types.ErrTooMuchWrkChainFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorRejectIncorrectDenomFeeInTx(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, test_helpers.TestDenomination))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt)
	feeDenom := "rubbish"

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr := sdkerrors.Wrap(types.ErrIncorrectFeeDenomination, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx1")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err1.Error(), "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFunds(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
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

	// Register
	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRegFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFundsWithLocked(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRecFeeAmt - 2))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	lockedUnd := enttypes.LockedUnd{
		Owner:  addr.String(),
		Amount: sdk.NewInt64Coin(actualFeeDenom, 1),
	}
	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

	withLocked := initCoins.Add(lockedUnd.Amount)

	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRegFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectWrkChainFeeDecoratorAcceptValidTx(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRegFeeAmt))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// Register
	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeSufficientLocked(t *testing.T) {
	app := test_helpers.Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	txGen := test_helpers.EncodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
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

	lockedUnd := enttypes.LockedUnd{
		Owner:  addr.String(),
		Amount: sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt)),
	}
	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)
}
