package ante_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	errorsmod "cosmossdk.io/errors"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/beacon/ante"
	"github.com/unification-com/mainchain/x/beacon/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
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

func TestCorrectBeaconFeeDecoratorAddressNotExist(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.BeaconKeeper.SetParams(ctx, types.NewParams(24, 2, 2, simapp.TestDenomination, 200, 300))
	bParams := app.BeaconKeeper.GetParams(ctx)
	actualFeeAmt := bParams.FeeRegister
	actualFeeDenom := bParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fee payer does not exist
	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, int64(actualFeeAmt)))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", addr)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectBeaconFeeDecoratorRejectTooLittleFeeInTx(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.BeaconKeeper.SetParams(ctx, types.NewParams(24, 2, 2, "testnund", 200, 300))

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// Register
	feeInt := int64(actualRegFeeAmt - 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: 1, expected fees: %v%v, sent fees: %v%v", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)

	expectedErr := errorsmod.Wrap(types.ErrInsufficientBeaconFee, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt1 := int64(actualRecFeeAmt - 1)
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt1, feeDenom)
	expectedErr1 := errorsmod.Wrap(types.ErrInsufficientBeaconFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	expectedFees := actualPurchaseAmt * numToPurchase // fee is per slot
	feeInt2 := int64(actualPurchaseAmt*numToPurchase) - 1
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", expectedFees, actualFeeDenom, feeInt2, feeDenom)
	expectedErr2 := errorsmod.Wrap(types.ErrInsufficientBeaconFee, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	expectedFees3 := (actualPurchaseAmt * numToPurchase) + actualRegFeeAmt + actualRecFeeAmt
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("insufficient fee to pay for beacon tx. numMsgs in tx: 3, expected fees: %d%s, sent fees: %d%s", expectedFees3, actualFeeDenom, multiFees, feeDenom)
	expectedErr3 := errorsmod.Wrap(types.ErrInsufficientBeaconFee, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)

}

func TestCorrectBeaconFeeDecoratorRejectTooMuchFeeInTx(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)

	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.BeaconKeeper.SetParams(ctx, types.NewParams(24, 2, 2, "testnund", 200, 300))

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt + 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr := errorsmod.Wrap(types.ErrTooMuchBeaconFee, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt1 := int64(actualRecFeeAmt + 1)
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt1, feeDenom)
	expectedErr1 := errorsmod.Wrap(types.ErrTooMuchBeaconFee, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	expectedFees := actualPurchaseAmt * numToPurchase // fee is per slot
	feeInt2 := int64(actualPurchaseAmt*numToPurchase) + 1
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", expectedFees, actualFeeDenom, feeInt2, feeDenom)
	expectedErr2 := errorsmod.Wrap(types.ErrTooMuchBeaconFee, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	expectedFees3 := (actualPurchaseAmt * numToPurchase) + actualRegFeeAmt + actualRecFeeAmt
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("too much fee sent to pay for beacon tx. numMsgs in tx: 3, expected fees: %d%s, sent fees: %d%s", expectedFees3, actualFeeDenom, multiFees, feeDenom)
	expectedErr3 := errorsmod.Wrap(types.ErrTooMuchBeaconFee, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)

}

func TestCorrectBeaconFeeDecoratorRejectIncorrectDenomFeeInTx(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)

	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	app.BeaconKeeper.SetParams(ctx, types.NewParams(24, 2, 2, "testnund", 200, 300))

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt)
	feeDenom := "rubbish"

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr := errorsmod.Wrap(types.ErrIncorrectFeeDenomination, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx1")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)

	// Record
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr1 := errorsmod.Wrap(types.ErrIncorrectFeeDenomination, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1.Error(), err1.Error(), "unexpected type of error: %s", err1)

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err2 := antehandler(ctx, tx2, false)

	errMsg2 := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr2 := errorsmod.Wrap(types.ErrIncorrectFeeDenomination, errMsg2)

	require.NotNil(t, err2, "Did not error on invalid tx")
	require.Equal(t, expectedErr2.Error(), err2.Error(), "unexpected type of error: %s", err2)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	errMsg3 := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr3 := errorsmod.Wrap(types.ErrIncorrectFeeDenomination, errMsg3)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectBeaconFeeDecoratorCorrectFeeInsufficientFunds(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)

	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

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

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRegFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, feeInt2, actualFeeDenom)

	_, err = antehandler(ctx, tx2, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	expectedErr3 := errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, multiFees, actualFeeDenom)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectBeaconFeeDecoratorCorrectFeeInsufficientFundsWithLocked(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

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

	// Register
	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr := errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRegFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	feeInt1 := int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// PurchaseStorageAction
	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	expectedErr = errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, feeInt2, actualFeeDenom)

	_, err = antehandler(ctx, tx2, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err3 := antehandler(ctx, tx3, false)

	expectedErr3 := errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, multiFees, actualFeeDenom)

	require.NotNil(t, err3, "Did not error on invalid tx")
	require.Equal(t, expectedErr3.Error(), err3.Error(), "unexpected type of error: %s", err3)
}

func TestCorrectBeaconFeeDecoratorAcceptValidTx(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

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

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	feeInt1 := int64(actualRecFeeAmt)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)

	// PurchaseStorageAction
	err = app.BeaconKeeper.SetBeacon(ctx, types.Beacon{
		BeaconId:        1,
		Moniker:         "test",
		Name:            "test",
		LastTimestampId: 0,
		FirstIdInState:  0,
		NumInState:      0,
		RegTime:         0,
		Owner:           addr.String(),
	})
	require.NoError(t, err)

	err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, 1, 10)
	require.NoError(t, err)

	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx2, false)
	require.NoError(t, err)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx3, false)
	require.NoError(t, err)
}

func TestCorrectBeaconFeeDecoratorCorrectFeeSufficientLocked(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualRegFeeAmt := bParams.FeeRegister
	actualRecFeeAmt := bParams.FeeRecord
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom

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

	msg := types.NewMsgRegisterBeacon("test", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 := types.NewMsgRecordBeaconTimestamp(1, "test", 1, addr)
	feeInt1 := int64(actualRecFeeAmt)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt1))

	tx1, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg1}, fee1, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)

	// PurchaseStorageAction
	_ = app.BeaconKeeper.SetBeacon(ctx, types.Beacon{
		BeaconId:        1,
		Moniker:         "test",
		Name:            "test",
		LastTimestampId: 0,
		FirstIdInState:  0,
		NumInState:      0,
		RegTime:         0,
		Owner:           addr.String(),
	})

	err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, 1, 10)
	require.NoError(t, err)

	numToPurchase := uint64(10)
	feeInt2 := int64(actualPurchaseAmt * numToPurchase)
	msg2 := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee2 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt2))

	tx2, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg2}, fee2, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx2, false)
	require.NoError(t, err)

	// Multi Msg
	multiFees := feeInt + feeInt1 + feeInt2
	fee3 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, multiFees))
	tx3, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg, msg1, msg2}, fee3, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx3, false)
	require.NoError(t, err)
}

func TestExceedsMaxStorageDecoratorInvalidTx(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapp.Setup(t, true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	encodingConfig := simapp.MakeEncodingConfig()
	txGen := encodingConfig.TxConfig

	feeDecorator := ante.NewCorrectBeaconFeeDecorator(app.BankKeeper, app.AccountKeeper, app.BeaconKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	bParams := app.BeaconKeeper.GetParams(ctx)
	actualPurchaseAmt := bParams.FeePurchaseStorage
	actualFeeDenom := bParams.Denom
	startInStateLimit := uint64(100)
	numToPurchase := simapp.TestMaxStorage - startInStateLimit + 1
	bId := uint64(1)

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
	_ = app.BeaconKeeper.SetBeacon(ctx, types.Beacon{
		BeaconId:        bId,
		Moniker:         "test",
		Name:            "test",
		LastTimestampId: 0,
		FirstIdInState:  0,
		NumInState:      0,
		RegTime:         0,
		Owner:           addr.String(),
	})

	err = app.BeaconKeeper.SetBeaconStorageLimit(ctx, bId, startInStateLimit)
	require.NoError(t, err)

	feeInt := int64(actualPurchaseAmt * numToPurchase)
	msg := types.NewMsgPurchaseBeaconStateStorage(1, numToPurchase, addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)

	expectedErr := errorsmod.Wrapf(types.ErrExceedsMaxStorage,
		"num slots exceeds max for beacon %d. Max can purchase: %d. Want in Msgs: %d", bId, simapp.TestMaxStorage-startInStateLimit, numToPurchase)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
}
