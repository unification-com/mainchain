package ante_test

import (
	"fmt"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"
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
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, 2, test_helpers.TestDenomination, 200, 300))

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
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, 2, test_helpers.TestDenomination, 200, 300))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
	feeInt1 := int64(actualRecFeeAmt - 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt1, feeDenom)
	expectedErr1 := sdkerrors.Wrap(types.ErrInsufficientWrkChainFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	expectedFees := actualPurchaseAmt * numToPurchase // fee is per slot
	feeInt2 := int64(actualPurchaseAmt*numToPurchase) - 1
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", expectedFees, actualFeeDenom, feeInt2, feeDenom)
	expectedErr2 := sdkerrors.Wrap(types.ErrInsufficientWrkChainFee, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	expectedFees3 := (actualPurchaseAmt * numToPurchase) + actualRegFeeAmt + actualRecFeeAmt
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 3, expected fees: %d%s, sent fees: %d%s", expectedFees3, actualFeeDenom, multiFees, feeDenom)
	expectedErr3 := sdkerrors.Wrap(types.ErrInsufficientWrkChainFee, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectWrkChainFeeDecoratorRejectTooMuchFeeInTx(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, 2, test_helpers.TestDenomination, 200, 300))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
	feeInt1 := int64(actualRecFeeAmt + 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt1, feeDenom)
	expectedErr1 := sdkerrors.Wrap(types.ErrTooMuchWrkChainFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	expectedFees := actualPurchaseAmt * numToPurchase // fee is per slot
	feeInt2 := int64(actualPurchaseAmt*numToPurchase) + 1
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", expectedFees, actualFeeDenom, feeInt2, feeDenom)
	expectedErr2 := sdkerrors.Wrap(types.ErrTooMuchWrkChainFee, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	expectedFees3 := (actualPurchaseAmt * numToPurchase) + actualRegFeeAmt + actualRecFeeAmt
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 3, expected fees: %d%s, sent fees: %d%s", expectedFees3, actualFeeDenom, multiFees, feeDenom)
	expectedErr3 := sdkerrors.Wrap(types.ErrTooMuchWrkChainFee, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectWrkChainFeeDecoratorRejectIncorrectDenomFeeInTx(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.WrkchainKeeper.SetParams(ctx, types.NewParams(24, 2, 2, test_helpers.TestDenomination, 200, 300))

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr2 := sdkerrors.Wrap(types.ErrIncorrectFeeDenomination, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr3 := sdkerrors.Wrap(types.ErrIncorrectFeeDenomination, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFunds(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, feeInt2, actualFeeDenom)

	_, err = antehandler(ctx, tx2, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	expectedErr3 := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, multiFees, actualFeeDenom)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFundsWithLocked(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, feeInt2, actualFeeDenom)

	_, err = antehandler(ctx, tx2, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	expectedErr3 := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, multiFees, actualFeeDenom)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectWrkChainFeeDecoratorAcceptValidTx(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRegFeeAmt * 3))
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
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)

	// PurchaseStorageAction
	err = app.WrkchainKeeper.SetWrkChain(ctx, types.WrkChain{
		WrkchainId:   1,
		Moniker:      "test",
		Name:         "test",
		Type:         "test",
		Genesis:      "genesishash",
		LowestHeight: 0,
		Lastblock:    0,
		NumBlocks:    0,
		RegTime:      0,
		Owner:        addr.String(),
	})
	require.NoError(t, err)

	err = app.WrkchainKeeper.SetWrkChainStorageLimit(ctx, 1, 10)
	require.NoError(t, err)

	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx2, false)
	require.NoError(t, err)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx3, false)
	require.NoError(t, err)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeSufficientLocked(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualPurchaseAmt := wrkParams.FeePurchaseStorage
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
		Amount: sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt*3)),
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
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1, "test", "test", "", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)

	// PurchaseStorageAction
	_ = app.WrkchainKeeper.SetWrkChain(ctx, types.WrkChain{
		WrkchainId:   1,
		Moniker:      "test",
		Name:         "test",
		Type:         "test",
		Genesis:      "genesishash",
		LowestHeight: 0,
		Lastblock:    0,
		NumBlocks:    0,
		RegTime:      0,
		Owner:        addr.String(),
	})

	err = app.WrkchainKeeper.SetWrkChainStorageLimit(ctx, 1, 10)
	require.NoError(t, err)

	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx2, false)
	require.NoError(t, err)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx3, false)
	require.NoError(t, err)
}

func TestExceedsMaxStorageDecoratorInvalidTx(t *testing.T) {
	app := test_helpers.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := test_helpers.GetAppEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.BankKeeper, app.AccountKeeper, app.WrkchainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wcParams := app.WrkchainKeeper.GetParams(ctx)
	actualPurchaseAmt := wcParams.FeePurchaseStorage
	actualFeeDenom := wcParams.Denom
	startInStateLimit := uint64(100)
	numToPurchase := test_helpers.TestMaxStorage - startInStateLimit + 1
	wcId := uint64(1)

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualPurchaseAmt * numToPurchase))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// PurchaseStorageAction
	_ = app.WrkchainKeeper.SetWrkChain(ctx, types.WrkChain{
		WrkchainId:   wcId,
		Moniker:      "test",
		Name:         "test",
		Type:         "test",
		Genesis:      "genesishash",
		LowestHeight: 0,
		Lastblock:    0,
		NumBlocks:    0,
		RegTime:      0,
		Owner:        addr.String(),
	})

	err = app.WrkchainKeeper.SetWrkChainStorageLimit(ctx, 1, startInStateLimit)
	require.NoError(t, err)

	feeInt := int64(actualPurchaseAmt * numToPurchase)
	msg := types.NewMsgPurchaseWrkChainStateStorage(1, numToPurchase, addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := test_helpers.GenTx(txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)

	expectedErr := sdkerrors.Wrapf(types.ErrExceedsMaxStorage,
		"num slots exceeds max for wrkchain %d. Max can purchase: %d. Want in Msgs: %d", wcId, test_helpers.TestMaxStorage-startInStateLimit, numToPurchase)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
}
