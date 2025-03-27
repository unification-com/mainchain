package ante_test

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/stretchr/testify/suite"
	simapp "github.com/unification-com/mainchain/app"
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	"math/rand"
	"testing"
	"time"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/enterprise/ante"
	"github.com/unification-com/mainchain/x/enterprise/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

const TestChainID = "und-unit-test-chain"

type AnteTestSuite struct {
	suite.Suite

	app         *simapp.App
	ctx         sdk.Context
	txGen       client.TxConfig
	anteHandler sdk.AnteHandler
	privKey     *ed25519.PrivKey
	addr        sdk.AccAddress
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

func (s *AnteTestSuite) SetupTest() {
	app := simapphelpers.Setup(s.T())
	ctx := app.BaseApp.NewContext(false)

	// set some default params
	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	wrkParams.Denom = sdk.DefaultBondDenom
	wrkParams.FeeRegister = 100
	wrkParams.FeeRecord = 1
	wrkParams.FeePurchaseStorage = 100
	_ = app.WrkchainKeeper.SetParams(ctx, wrkParams)

	beaconParams := app.BeaconKeeper.GetParams(ctx)
	beaconParams.Denom = sdk.DefaultBondDenom
	beaconParams.FeeRegister = 100
	beaconParams.FeeRecord = 1
	beaconParams.FeePurchaseStorage = 100
	_ = app.BeaconKeeper.SetParams(ctx, beaconParams)

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	anteHandler := sdk.ChainAnteDecorators(feeDecorator)

	privK, addr := simapphelpers.AddTestAccForTxSigning(app, ctx, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)))

	// Note - default total supply with this setup is 100000001000100nund

	s.app = app
	s.ctx = ctx
	s.txGen = app.GetTxConfig()
	s.anteHandler = anteHandler
	s.privKey = privK
	s.addr = addr
}

func (s *AnteTestSuite) TestAnteHandler() {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	testCases := []struct {
		name           string
		msgs           []sdk.Msg
		feeToSend      sdk.Coins
		toLock         sdk.Coin
		expectErr      bool
		expErrMsg      string
		expTotalSupply sdk.Coin
		expAccLocked   sdk.Coin
		expTotalLocked sdk.Coin
		expSpent       sdk.Coin
	}{
		{
			name: "full fee correctly minted for MsgRegisterWrkChain no locked left",
			msgs: []sdk.Msg{
				&wrkchaintypes.MsgRegisterWrkChain{
					Moniker:     "test1",
					Name:        "test1",
					GenesisHash: "test",
					BaseType:    "geth",
					Owner:       s.addr.String(),
				},
			},
			feeToSend:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			toLock:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 100),
			expectErr:      false,
			expErrMsg:      "",
			expTotalSupply: sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000001000200),
			expAccLocked:   sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expTotalLocked: sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expSpent:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 100),
		},
		{
			name: "full fee correctly minted for MsgRegisterBeacon no locked left",
			msgs: []sdk.Msg{
				&beacontypes.MsgRegisterBeacon{
					Moniker: "test1",
					Name:    "test1",
					Owner:   s.addr.String(),
				},
			},
			feeToSend:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			toLock:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 100),
			expectErr:      false,
			expErrMsg:      "",
			expTotalSupply: sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000001000300),
			expAccLocked:   sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expTotalLocked: sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expSpent:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 200),
		},
		{
			name: "full fee correctly minted for MsgRegisterBeacon 50 locked left",
			msgs: []sdk.Msg{
				&beacontypes.MsgRegisterBeacon{
					Moniker: "test1",
					Name:    "test1",
					Owner:   s.addr.String(),
				},
			},
			feeToSend:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			toLock:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 150),
			expectErr:      false,
			expErrMsg:      "",
			expTotalSupply: sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000001000400),
			expAccLocked:   sdk.NewInt64Coin(sdk.DefaultBondDenom, 50),
			expTotalLocked: sdk.NewInt64Coin(sdk.DefaultBondDenom, 50),
			expSpent:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 300),
		},
		{
			name: "remaining locked correctly minted for MsgRegisterBeacon no locked left",
			msgs: []sdk.Msg{
				&beacontypes.MsgRegisterBeacon{
					Moniker: "test1",
					Name:    "test1",
					Owner:   s.addr.String(),
				},
			},
			feeToSend:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			toLock:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), // still has 50 from previous test
			expectErr:      false,
			expErrMsg:      "",
			expTotalSupply: sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000001000450),
			expAccLocked:   sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expTotalLocked: sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expSpent:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 350),
		},
		{
			name: "nothing locked, nothing minted",
			msgs: []sdk.Msg{
				&beacontypes.MsgRegisterBeacon{
					Moniker: "test1",
					Name:    "test1",
					Owner:   s.addr.String(),
				},
			},
			feeToSend:      sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)),
			toLock:         sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expectErr:      false,
			expErrMsg:      "",
			expTotalSupply: sdk.NewInt64Coin(sdk.DefaultBondDenom, 100000001000450), // no change to total supply
			expAccLocked:   sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expTotalLocked: sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
			expSpent:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 350),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_ = s.app.EnterpriseKeeper.CreateAndLockEFUND(s.ctx, s.addr, tc.toLock)

			tx, _ := simtestutil.GenSignedMockTx(r, s.txGen, tc.msgs, tc.feeToSend, uint64(0), simapphelpers.SimAppChainID, []uint64{0}, []uint64{0}, s.privKey)

			_, err := s.anteHandler(s.ctx, tx, false)

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().NoError(err)
			}

			totalSupplyAfter := s.app.BankKeeper.GetSupply(s.ctx, sdk.DefaultBondDenom)
			s.Require().Equal(tc.expTotalSupply, totalSupplyAfter)

			accLockedAfter := s.app.EnterpriseKeeper.GetLockedUndForAccount(s.ctx, s.addr)
			s.Require().Equal(tc.expAccLocked, accLockedAfter.Amount)

			totalLockedAfter := s.app.EnterpriseKeeper.GetTotalLockedUnd(s.ctx)
			s.Require().Equal(tc.expTotalLocked, totalLockedAfter)

			spentEfundForAcc := s.app.EnterpriseKeeper.GetSpentEFUNDForAccount(s.ctx, s.addr)
			s.Require().Equal(tc.expSpent, spentEfundForAcc.Amount)

			totalSpentEfund := s.app.EnterpriseKeeper.GetTotalSpentEFUND(s.ctx)
			s.Require().Equal(tc.expSpent, totalSpentEfund)
		})
	}

}

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

func TestCheckLockedUndDecoratorModuleAndSupplyMinting(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(true)
	txGen := app.GetTxConfig()

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualFeeDenom := wrkParams.Denom
	actualFeeAmount := wrkParams.FeeRegister

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := mathmod.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// add locked eFUND
	feeCoin := sdk.NewCoin(actualFeeDenom, mathmod.NewIntFromUint64(actualFeeAmount))
	lockeUnd := types.LockedUnd{
		Owner:  addr.String(),
		Amount: feeCoin,
	}

	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)

	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockeUnd)

	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(feeCoin)

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	require.Equal(t, totalSupplyBefore.Add(feeCoin).String(), totalSupplyAfter.String())

}

func TestOnlyMintsAmountLocked(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(true)
	txGen := app.GetTxConfig()

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualFeeDenom := wrkParams.Denom
	actualFeeAmount := wrkParams.FeeRegister

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := mathmod.NewInt(int64(actualFeeAmount / 2))
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// add locked eFUND, half the value of the actual reg fee
	toLock := sdk.NewCoin(actualFeeDenom, mathmod.NewIntFromUint64(actualFeeAmount/2))
	lockeUnd := types.LockedUnd{
		Owner:  addr.String(),
		Amount: toLock,
	}

	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)

	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockeUnd)

	feeCoin := sdk.NewCoin(actualFeeDenom, mathmod.NewIntFromUint64(actualFeeAmount))
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(feeCoin)

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	require.Equal(t, totalSupplyBefore.Add(toLock).String(), totalSupplyAfter.String())

}

func TestCheckLockedUndDecoratorSuccessfulUnlock(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(true)
	txGen := app.GetTxConfig()

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	accAmt := mathmod.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	_ = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, sdk.NewInt64Coin(actualFeeDenom, int64(actualRegFeeAmt)))

	feeInt := int64(1)
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)
}

func TestCheckLockedUndDecoratorSkipIfNothingLocked(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(true)
	txGen := app.GetTxConfig()

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualRegFeeAmt := wrkParams.FeeRegister
	actualFeeDenom := wrkParams.Denom

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := mathmod.NewInt(int64(1))
	initCoins := sdk.NewCoins(sdk.NewCoin(actualFeeDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	feeInt := int64(actualRegFeeAmt)
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(sdk.NewInt64Coin(actualFeeDenom, feeInt))

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)
}

func TestNoMintingIfInsufficientBalanceAndLocked(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(true)
	txGen := app.GetTxConfig()

	feeDecorator := ante.NewCheckLockedUndDecorator(app.EnterpriseKeeper)
	antehandler := sdk.ChainAnteDecorators(feeDecorator)

	wrkParams := app.WrkchainKeeper.GetParams(ctx)
	actualFeeDenom := wrkParams.Denom
	actualFeeAmount := wrkParams.FeeRegister

	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	// fund the account
	accAmt := mathmod.NewInt(int64(actualFeeAmount / 2))
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, accAmt))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(ctx, app.BankKeeper, addr, initCoins)
	require.NoError(t, err)

	// add locked eFUND, half the value of the actual reg fee
	toLock := sdk.NewCoin(actualFeeDenom, mathmod.NewIntFromUint64(1))
	lockeUnd := types.LockedUnd{
		Owner:  addr.String(),
		Amount: toLock,
	}

	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)

	_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockeUnd)

	feeCoin := sdk.NewCoin(actualFeeDenom, mathmod.NewIntFromUint64(actualFeeAmount))
	msg := wrkchaintypes.NewMsgRegisterWrkChain("test", "hash", "Test", "geth", addr)
	fee := sdk.NewCoins(feeCoin)

	tx, _ := simtestutil.GenSignedMockTx(r, txGen, []sdk.Msg{msg}, fee, uint64(0), TestChainID, []uint64{0}, []uint64{0}, privK)

	_, err = antehandler(ctx, tx, false)
	require.NoError(t, err)

	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	require.Equal(t, totalSupplyBefore.String(), totalSupplyAfter.String())

	// check locked is still 1
	lockedAfter := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
	require.Equal(t, lockeUnd.Amount, lockedAfter.Amount)
}
