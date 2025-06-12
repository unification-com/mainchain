package helpers

import (
	"context"
	"math/rand"
	"time"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/app"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *app.App, ctx context.Context, accNum int, accAmt mathmod.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, simtestutil.CreateIncrementalAccounts)
}

func AddTestAccForTxSigning(app *app.App, ctx context.Context, initCoins sdk.Coins) (*ed25519.PrivKey, sdk.AccAddress) {
	privK := ed25519.GenPrivKey()
	pubK := privK.PubKey()
	addr := sdk.AccAddress(pubK.Address())

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)

	initAccountWithCoins(app, ctx, addr, initCoins)

	return privK, addr
}

func addTestAddrs(app *app.App, ctx context.Context, accNum int, accAmt mathmod.Int, strategy simtestutil.GenerateAccountStrategy) []sdk.AccAddress {
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

func initAccountWithCoins(app *app.App, ctx context.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, enttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
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
		privK := mock.NewPV()
		pubK := privK.PrivKey.PubKey()
		testAddrs[i] = sdk.AccAddress(pubK.Address())
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
func AddTestAddrs(app *app.App, ctx context.Context, accNum int, accAmt mathmod.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

func AddTestAddrsWithExtraNonBondCoin(app *app.App, ctx context.Context, accNum int, accAmt mathmod.Int, extraCoin sdk.Coin) []sdk.AccAddress {
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
