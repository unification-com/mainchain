package ante_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain-cosmos/x/enterprise"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain-cosmos/simapp"
	"github.com/unification-com/mainchain-cosmos/simapp/helpers"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/ante"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

const TestChainID = "und-unit-test-chain"

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.UndSimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{ChainID: TestChainID})
	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.WrkChainKeeper.SetParams(ctx, wrkchain.DefaultParams())
	app.EnterpriseKeeper.SetParams(ctx, enterprise.DefaultParams())

	return app, ctx
}

func TestCorrectWrkChainFeeDecoratorAddressNotExist(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fee payer does not exist
	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, int64(actualFeeAmt)))
	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", addr)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectWrkChainFeeDecoratorRejectTooLittleFeeInTx(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// Register
	feeInt := int64(actualRegFeeAmt - 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash","Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr := types.ErrInsufficientWrkChainFee(types.DefaultCodespace, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr, err, "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt - 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("insufficient fee to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr1 := types.ErrInsufficientWrkChainFee(types.DefaultCodespace, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1, err1, "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorRejectTooMuchFeeInTx(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt + 1)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRegFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr := types.ErrTooMuchWrkChainFee(types.DefaultCodespace, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx")
	require.Equal(t, expectedErr, err, "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt + 1)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err1 := antehandler(ctx, tx1, false)

	errMsg1 := fmt.Sprintf("too much fee sent to pay for WrkChain tx. numMsgs in tx: 1, expected fees: %d%s, sent fees: %d%s", actualRecFeeAmt, actualFeeDenom, feeInt, feeDenom)
	expectedErr1 := types.ErrTooMuchWrkChainFee(types.DefaultCodespace, errMsg1)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr1, err1, "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorRejectIncorrectDenomFeeInTx(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	feeInt := int64(actualRegFeeAmt)
	feeDenom := "rubbish"

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)

	errMsg := fmt.Sprintf("incorrect fee denomination. expected %s", actualFeeDenom)
	expectedErr := types.ErrIncorrectFeeDenomination(types.DefaultCodespace, errMsg)

	require.NotNil(t, err, "Did not error on invalid tx1")
	require.Equal(t, expectedErr, err, "unexpected type of error: %s", err)

	// Record
	feeInt = int64(actualRecFeeAmt)
	msg1 := types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err1 := antehandler(ctx, tx1, false)

	require.NotNil(t, err1, "Did not error on invalid tx")
	require.Equal(t, expectedErr, err1, "unexpected type of error: %s", err1)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFunds(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRecFeeAmt - 1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, initCoins)

	// Register
	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRegFeeAmt, actualFeeDenom)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	msg1 :=  types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, initCoins, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeInsufficientFundsWithLocked(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRecFeeAmt - 2))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, initCoins)

	lockedUnd := enterprise.LockedUnd{
		Owner: addr,
		Amount: sdk.NewInt64Coin(actualFeeDenom, 1),
	}
	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

	withLocked := initCoins.Add(sdk.NewCoins(lockedUnd.Amount))

	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	expectedErr := sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRegFeeAmt, actualFeeDenom)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}

	// Record
	msg1 :=  types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	expectedErr = sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
		"insufficient und to pay for fees. unlocked und: %s, including locked und: %s, fee: %d%s", initCoins, withLocked, actualRecFeeAmt, actualFeeDenom)

	_, err = antehandler(ctx, tx1, false)
	require.NotNil(t, err, "Did not error on invalid tx")

	if err != nil {
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err)
	}
}


func TestCorrectWrkChainFeeDecoratorAcceptValidTx(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRegFeeAmt))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, initCoins)

	// Register
	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 :=  types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)
}

func TestCorrectWrkChainFeeDecoratorCorrectFeeSufficientLocked(t *testing.T) {
	app, ctx := createTestApp(true)

	feeDecorator := ante.NewCorrectWrkChainFeeDecorator(app.AccountKeeper, app.WrkChainKeeper, app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualRecFeeAmt := wrkParams.FeeRecord
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(0))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, initCoins)

	lockedUnd := enterprise.LockedUnd{
		Owner: addr,
		Amount: sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt)),
	}
	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

	feeInt := int64(actualRegFeeAmt)
	feeDenom := actualFeeDenom

	msg := types.NewMsgRegisterWrkChain("test", "hash", "Test", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)
	require.NoError(t, err)

	// Record
	msg1 :=  types.NewMsgRecordWrkChainBlock(1, 1,"test", "test","", "", "", addr)
	fee1 := sdk.NewCoins(sdk.NewInt64Coin(feeDenom, int64(actualRecFeeAmt)))

	tx1 := helpers.GenTx(
		[]sdk.Msg{msg1},
		fee1,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err = antehandler(ctx, tx1, false)
	require.NoError(t, err)
}
