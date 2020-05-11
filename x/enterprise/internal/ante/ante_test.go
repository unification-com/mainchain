package ante_test

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain/simapp/helpers"
	undtypes "github.com/unification-com/mainchain/types"
	"github.com/unification-com/mainchain/x/beacon"
	"github.com/unification-com/mainchain/x/enterprise/internal/ante"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
	"github.com/unification-com/mainchain/x/wrkchain"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/enterprise"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/simapp"
)

const TestChainID = "und-unit-test-chain"

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.UndSimApp, sdk.Context) {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	config.SetCoinType(undtypes.CoinType)
	config.SetFullFundraiserPath(undtypes.HdWalletPath)

	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	signerAddr := sdk.AccAddress(pubKey2.Address())

	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{ChainID: TestChainID})
	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.WrkChainKeeper.SetParams(ctx, wrkchain.DefaultParams())
	app.BeaconKeeper.SetParams(ctx, beacon.DefaultParams())

	entParams := enterprise.DefaultParams()
	entParams.EntSigners = signerAddr.String()
	app.EnterpriseKeeper.SetParams(ctx, entParams)

	return app, ctx
}

func TestCheckLockedUndDecoratorModuleAndSupplyInsufficientFunds(t *testing.T) {
	app, ctx := createTestApp(false)

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt)))

	// artificially add locked FUND without minting first
	toLock := sdk.NewCoin(actualFeeDenom, accAmt)
	lockeUnd := types.LockedUnd{
		Owner:  addr,
		Amount: toLock,
	}

	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockeUnd)

	feeInt := int64(1)
	msg := wrkchain.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fee,
		TestChainID,
		[]uint64{0},
		[]uint64{0},
		privK,
	)

	_, err := antehandler(ctx, tx, false)
	require.NotNil(t, err, "Did not error on invalid tx")

}

func TestCheckLockedUndDecoratorSuccessfulUnlock(t *testing.T) {
	app, ctx := createTestApp(false)

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt)))

	_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt)))

	feeInt := int64(1)
	msg := wrkchain.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

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
}

func TestCheckLockedUndDecoratorSkipIfNothingLocked(t *testing.T) {
	app, ctx := createTestApp(false)

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkChainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := sdk.NewInt(int64(actualRegFeeAmt))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	totalSupply := initCoins
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	_, _ = app.BankKeeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt)))

	feeInt := int64(actualRegFeeAmt)
	msg := wrkchain.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

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
}
