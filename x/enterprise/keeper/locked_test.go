package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetTotalLockedUnd(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	denom := test_helpers.TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	require.NoError(t, err)

	lockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)

	require.True(t, lockedDb.IsEqual(locked))
	require.True(t, lockedDb.Denom == denom)
	require.True(t, lockedDb.Amount.Int64() == amount)
}

func TestGetTotalUnlocked(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	test_helpers.AddTestAddrs(app, ctx, 1, sdk.NewInt(20000))

	denom := test_helpers.TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	require.NoError(t, err)

	totUnlocked := app.EnterpriseKeeper.GetTotalUnLockedUnd(ctx)
	totalSupply := app.BankKeeper.GetSupply(ctx).GetTotal()

	diff := totalSupply.Sub(sdk.Coins{totUnlocked})

	require.Equal(t, sdk.Coins{locked}, diff)
}

func TestGetTotalUndSupply(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)
	test_helpers.AddTestAddrs(app, ctx, 1, sdk.NewInt(20000))

	totalSupply := app.BankKeeper.GetSupply(ctx).GetTotal()
	totalSupplyFromEnt := app.EnterpriseKeeper.GetTotalUndSupply(ctx)
	require.Equal(t, totalSupply, sdk.Coins{totalSupplyFromEnt})
}

func TestSetGetLockedUndForAccount(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	testAddresses := test_helpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		denom := test_helpers.TestDenomination

		locked := types.LockedUnd{
			Owner:  addr.String(),
			Amount: sdk.NewInt64Coin(denom, amount),
		}

		err := app.EnterpriseKeeper.SetLockedUndForAccount(ctx, locked)
		require.NoError(t, err)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)

		require.True(t, locked.Owner == lockedDb.Owner)
		require.True(t, lockedDb.Amount.IsEqual(locked.Amount))

		lockedDbAmount := app.EnterpriseKeeper.GetLockedUndAmountForAccount(ctx, addr)
		require.True(t, lockedDbAmount.IsEqual(locked.Amount))
	}
}

func (suite *KeeperTestSuite) TestIsLocked() {
	app, ctx, addrs := suite.app, suite.ctx, suite.addrs

	denom := test_helpers.TestDenomination

	var (
		l    types.LockedUnd
		addr sdk.AccAddress
	)

	testCases := []struct {
		msg         string
		malleate    func()
		expIsLocked bool
	}{
		{
			"zero value",
			func() {
				addr = addrs[0]
				l = types.LockedUnd{
					Owner:  addr.String(),
					Amount: sdk.NewInt64Coin(denom, 0),
				}
			},
			false,
		},
		{
			"valid value",
			func() {
				addr = addrs[2]
				l = types.LockedUnd{
					Owner:  addr.String(),
					Amount: sdk.NewInt64Coin(denom, 100),
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			err := app.EnterpriseKeeper.SetLockedUndForAccount(ctx, l)
			suite.Require().NoError(err)

			isLocked := app.EnterpriseKeeper.IsLocked(ctx, addr)

			suite.Require().Equal(testCase.expIsLocked, isLocked)
		})
	}
}

func TestMintCoinsAndLock(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalAmount := int64(0)

	testAddresses := test_helpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		totalAmount = totalAmount + amount

		toMint := sdk.NewInt64Coin(test_helpers.TestDenomination, amount)

		err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)
		require.NoError(t, err)

		isLocked := app.EnterpriseKeeper.IsLocked(ctx, addr)
		require.True(t, isLocked)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(toMint))
	}

	totalLocked := sdk.NewInt64Coin(test_helpers.TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := app.EnterpriseKeeper.GetTotalSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Uint64())

	entAccount := app.EnterpriseKeeper.GetEnterpriseAccount(ctx)
	entAccountCoins := app.BankKeeper.GetAllBalances(ctx, entAccount.GetAddress())
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromAccK := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := app.BankKeeper.GetAllBalances(ctx, entAccFromAccK.GetAddress())
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestUnlockCoinsForFees(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalAmount := int64(0)

	testAddresses := test_helpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(test_helpers.RandInBetween(1000, 100000))
		amountToUnlock := int64(test_helpers.RandInBetween(1, 999))
		totalAmount = totalAmount + amountToMint

		toMint := sdk.NewInt64Coin(test_helpers.TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(test_helpers.TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalAmount = totalAmount - amountToUnlock

		expectedLocked := toMint.Sub(toUnlock)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(expectedLocked))
	}

	totalLocked := sdk.NewInt64Coin(test_helpers.TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := app.EnterpriseKeeper.GetTotalSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Uint64())

	entAccount := app.EnterpriseKeeper.GetEnterpriseAccount(ctx)
	entAccountCoins := app.BankKeeper.GetAllBalances(ctx, entAccount.GetAddress())
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromAccK := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := app.BankKeeper.GetAllBalances(ctx, entAccFromAccK.GetAddress())
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestGetTotalSupplyWithLockedNundRemoved(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	test_helpers.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalSupply := sdk.NewCoins(sdk.NewInt64Coin(test_helpers.TestDenomination, 0))
	totalMinted := sdk.NewCoins(sdk.NewInt64Coin(test_helpers.TestDenomination, 0))

	testAddresses := test_helpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(test_helpers.RandInBetween(1000, 100000))
		amountToUnlock := int64(test_helpers.RandInBetween(1, 999))

		toMint := sdk.NewInt64Coin(test_helpers.TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(test_helpers.TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalSupply = totalSupply.Add(toUnlock)
		totalMinted = totalMinted.Add(toMint)

		totalSupplyDb := app.EnterpriseKeeper.GetTotalSupplyWithLockedNundRemoved(ctx)
		require.True(t, totalSupplyDb.IsEqual(totalSupply))

		totalMintedDb := app.BankKeeper.GetSupply(ctx).GetTotal()
		require.True(t, totalMintedDb.IsEqual(totalMinted))
	}
}
