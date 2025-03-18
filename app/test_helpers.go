package app

import (
	"context"
	"encoding/json"
	"fmt"
	undtypes "github.com/unification-com/mainchain/types"
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/log"
	mathmod "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v8/testing/mock"
	"github.com/stretchr/testify/require"

	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	streamtypes "github.com/unification-com/mainchain/x/stream/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	TestDenomination   = "stake"
	SimAppChainID      = "simulation-app"
	charset            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	TestDefaultStorage = 200
	TestMaxStorage     = 300
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// App testing.
var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// SetupOptions defines arguments that are passed into `Simapp` constructor.
type SetupOptions struct {
	Logger  log.Logger
	DB      *dbm.MemDB
	AppOpts servertypes.AppOptions
}

func setup(withGenesis bool, invCheckPeriod uint) (*App, GenesisState) {
	db := dbm.NewMemDB()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod

	app := NewApp(log.NewNopLogger(), db, nil, true, appOptions)

	fmt.Println(app.DefaultGenesis())
	// note - SDK sim returns app, app.DefaultGenesis().
	// since FUND modules use nund to test, need to convert genesis
	if withGenesis {
		genesisState := app.DefaultGenesis()
		genesisState = convertGenesisStateToNund(app, genesisState)
		return app, genesisState
	}

	return app, GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(t *testing.T, isCheckTx bool) *App {
	t.Helper()
	config := sdk.GetConfig()
	if config.GetBech32AccountAddrPrefix() != undtypes.Bech32PrefixAccAddr {
		SetConfig()
	}

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := cmtypes.NewValidator(pubKey, 1)
	valSet := cmtypes.NewValidatorSet([]*cmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(TestDenomination, mathmod.NewInt(100000000000000))),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	return app
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, valSet *cmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *App {
	t.Helper()

	app, genesisState := setup(true, 5)
	genesisState, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(&abci.RequestInitChain{
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: simtestutil.DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	},
	)

	require.NoError(t, err)

	return app
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *App, ctx context.Context, accNum int, accAmt mathmod.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, simtestutil.CreateIncrementalAccounts)
}

func addTestAddrs(app *App, ctx context.Context, accNum int, accAmt mathmod.Int, strategy simtestutil.GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	if err != nil {
		panic(err)
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}

func initAccountWithCoins(app *App, ctx context.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, enttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

func convertGenesisStateToNund(app *App, genesisState map[string]json.RawMessage) map[string]json.RawMessage {

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.BondDenom = TestDenomination
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, nil, nil)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	govGenesis := govtypesv1.DefaultGenesisState()
	govGenesis.Params.MinDeposit = sdk.Coins{sdk.NewCoin(TestDenomination, mathmod.NewIntFromUint64(10000000))}
	genesisState[govtypes.ModuleName] = app.AppCodec().MustMarshalJSON(govGenesis)

	crisisGenesis := crisistypes.NewGenesisState(sdk.NewCoin(TestDenomination, mathmod.NewIntFromUint64(1000)))
	genesisState[crisistypes.ModuleName] = app.AppCodec().MustMarshalJSON(crisisGenesis)

	return genesisState
}

func SetKeeperTestParamsAndDefaultValues(app *App, ctx context.Context) {
	// Todo - migrate modules from sdk.Context to context.Context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	app.BeaconKeeper.SetParams(sdkCtx, beacontypes.NewParams(24, 2, 2, TestDenomination, TestDefaultStorage, TestMaxStorage))
	app.WrkchainKeeper.SetParams(sdkCtx, wrkchaintypes.NewParams(24, 2, 2, TestDenomination, TestDefaultStorage, TestMaxStorage))
	app.EnterpriseKeeper.SetParams(sdkCtx, enttypes.Params{
		EntSigners:        sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
		Denom:             TestDenomination,
		MinAccepts:        1,
		DecisionTimeLimit: 1000,
	})
	_ = app.EnterpriseKeeper.SetTotalLockedUnd(sdkCtx, sdk.NewInt64Coin(TestDenomination, 0))
	_ = app.EnterpriseKeeper.SetTotalSpentEFUND(sdkCtx, sdk.NewInt64Coin(TestDenomination, 0))
	app.EnterpriseKeeper.SetHighestPurchaseOrderID(sdkCtx, 1)

	// set to 0%. Individual unit tests will set specific %
	app.StreamKeeper.SetParams(sdkCtx, streamtypes.Params{ValidatorFee: "0.00"})
}

func GenerateRandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomString generates a random string given a length, based on a set character set
func GenerateRandomString(length int) string {
	return GenerateRandomStringWithCharset(length, charset)
}

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func GenerateRandomTestAccounts(accNum int) []sdk.AccAddress {
	return createRandomAccounts(accNum)
}

func RandInBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *App, ctx context.Context, accNum int, accAmt mathmod.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

func AddTestAddrsWithExtraNonBondCoin(app *App, ctx context.Context, accNum int, accAmt mathmod.Int, extraCoin sdk.Coin) []sdk.AccAddress {
	testAddrs := createRandomAccounts(accNum)
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	if err != nil {
		panic(err)
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, accAmt), extraCoin)

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs
}
