package simapp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/unification-com/mainchain/simapp/helpers"
)

type TestAccount struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  crypto.PubKey
	Address sdk.AccAddress
}

// Setup initializes a new UndSimApp. A Nop logger is set in UndSimApp.
func Setup(isCheckTx bool) *UndSimApp {
	db := dbm.NewMemDB()
	app := NewUndSimApp(log.NewNopLogger(), db, nil, true, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

// SetupWithGenesisAccounts initializes a new UndSimApp with the passed in
// genesis accounts.
func SetupWithGenesisAccounts(genAccs []authexported.GenesisAccount) *UndSimApp {
	db := dbm.NewMemDB()
	app := NewUndSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	// initialize the chain with the passed in genesis accounts
	genesisState := NewDefaultGenesisState()

	authGenesis := auth.NewGenesisState(auth.DefaultParams(), genAccs)
	genesisStateBz := app.cdc.MustMarshalJSON(authGenesis)
	genesisState[auth.ModuleName] = genesisStateBz

	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})

	return app
}

// SetupUnitTestApp creates a simApp and context for testing with some sensible chain defaults
func SetupUnitTestApp(isCheckTx bool, genAccs int, amt int64, testDenom string) (*UndSimApp, sdk.Context, TestAccount, []TestAccount) {

	accAmt := sdk.NewInt(amt)

	app := Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{ChainID: "und-unit-test-chain"})

	skParams := app.StakingKeeper.GetParams(ctx)
	skParams.BondDenom = testDenom
	app.StakingKeeper.SetParams(ctx, skParams)

	testAccs := make([]TestAccount, genAccs)
	for i := 0; i < genAccs; i++ {
		privK := ed25519.GenPrivKey()
		pubK := privK.PubKey()
		addr := sdk.AccAddress(pubK.Address())
		testAccs[0] = TestAccount{PrivKey: privK, PubKey: pubK, Address: addr}
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(len(testAccs)))))
	prevSupply := app.SupplyKeeper.GetSupply(ctx)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))

	app.BeaconKeeper.SetHighestBeaconID(ctx, 1)
	beaconParams := app.BeaconKeeper.GetParams(ctx)
	beaconParams.Denom = app.StakingKeeper.BondDenom(ctx)
	app.BeaconKeeper.SetParams(ctx, beaconParams)

	app.WrkChainKeeper.SetHighestWrkChainID(ctx, 1)
	wrkchainParams := app.WrkChainKeeper.GetParams(ctx)
	wrkchainParams.Denom = app.StakingKeeper.BondDenom(ctx)
	app.WrkChainKeeper.SetParams(ctx, wrkchainParams)

	entPrivK := ed25519.GenPrivKey()
	entPubKey := entPrivK.PubKey()
	entAddr := sdk.AccAddress(entPubKey.Address())
	entAcc := TestAccount{PrivKey: entPrivK, PubKey: entPubKey, Address: entAddr}

	app.EnterpriseKeeper.SetHighestPurchaseOrderID(ctx, 1)
	entParams := app.EnterpriseKeeper.GetParams(ctx)
	entParams.Denom = app.StakingKeeper.BondDenom(ctx)
	entParams.EntSigners = entAddr.String()
	app.EnterpriseKeeper.SetParams(ctx, entParams)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, ta := range testAccs {
		_, err := app.BankKeeper.AddCoins(ctx, ta.Address, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return app, ctx, entAcc, testAccs
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt
func AddTestAddrs(app *UndSimApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(len(testAddrs)))))
	prevSupply := app.SupplyKeeper.GetSupply(ctx)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		_, err := app.BankKeeper.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}
	return testAddrs
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *UndSimApp, addr sdk.AccAddress, exp sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	res := app.AccountKeeper.GetAccount(ctxCheck, addr)

	require.Equal(t, exp, res.GetCoins())
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *bam.BaseApp, header abci.Header, msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {

	tx := helpers.GenTx(
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
		"",
		accNums,
		seq,
		priv...,
	)

	txBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	_, res, err := app.Simulate(txBytes, tx)

	if expSimPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	gInfo, res, err := app.Deliver(tx)

	if expPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return gInfo, res, err
}

// GenSequenceOfTxs generates a set of signed transactions of messages, such
// that they differ only by having the sequence numbers incremented between
// every transaction.
func GenSequenceOfTxs(msgs []sdk.Msg, accNums []uint64, initSeqNums []uint64, numToGenerate int, priv ...crypto.PrivKey) []auth.StdTx {
	txs := make([]auth.StdTx, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		txs[i] = helpers.GenTx(
			msgs,
			sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
			"",
			accNums,
			initSeqNums,
			priv...,
		)
		incrementAllSequenceNumbers(initSeqNums)
	}

	return txs
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}
