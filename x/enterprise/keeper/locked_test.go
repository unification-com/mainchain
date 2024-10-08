package keeper_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	simapp "github.com/unification-com/mainchain/app"
	"math/rand"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetTotalLockedUnd(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	denom := simapp.TestDenomination
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
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(20000))

	denom := simapp.TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	require.NoError(t, err)

	totUnlocked := app.EnterpriseKeeper.GetTotalUnLockedUnd(ctx)
	totalSupply := app.BankKeeper.GetSupply(ctx, denom)

	diff := totalSupply.Sub(totUnlocked)

	require.Equal(t, locked, diff)
}

func TestGetTotalUndSupply(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(20000))

	totalSupply := app.BankKeeper.GetSupply(ctx, simapp.TestDenomination)
	totalSupplyFromEnt := app.EnterpriseKeeper.GetTotalUndSupply(ctx)
	require.Equal(t, totalSupply, totalSupplyFromEnt)
}

func TestSetGetLockedUndForAccount(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	testAddresses := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		denom := simapp.TestDenomination

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

func (s *KeeperTestSuite) TestIsLocked() {
	app, ctx, addrs := s.app, s.ctx, s.addrs

	denom := simapp.TestDenomination

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
		s.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			err := app.EnterpriseKeeper.SetLockedUndForAccount(ctx, l)
			s.Require().NoError(err)

			isLocked := app.EnterpriseKeeper.IsLocked(ctx, addr)

			s.Require().Equal(testCase.expIsLocked, isLocked)
		})
	}
}

func TestMintCoinsAndLock(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalAmount := int64(0)

	testAddresses := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		totalAmount = totalAmount + amount

		toMint := sdk.NewInt64Coin(simapp.TestDenomination, amount)

		err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)
		require.NoError(t, err)

		isLocked := app.EnterpriseKeeper.IsLocked(ctx, addr)
		require.True(t, isLocked)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(toMint))
	}

	totalLocked := sdk.NewInt64Coin(simapp.TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := app.EnterpriseKeeper.GetEnterpriseSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Uint64())

	entAccount := app.EnterpriseKeeper.GetEnterpriseAccount(ctx)
	entAccountCoins := app.BankKeeper.GetAllBalances(ctx, entAccount.GetAddress())
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromAccK := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := app.BankKeeper.GetAllBalances(ctx, entAccFromAccK.GetAddress())
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestUnlockCoinsForFees(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalAmount := int64(0)

	testAddresses := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(simapp.RandInBetween(1000, 100000))
		amountToUnlock := int64(simapp.RandInBetween(1, 999))
		totalAmount = totalAmount + amountToMint

		toMint := sdk.NewInt64Coin(simapp.TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(simapp.TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalAmount = totalAmount - amountToUnlock

		expectedLocked := toMint.Sub(toUnlock)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(expectedLocked))
	}

	totalLocked := sdk.NewInt64Coin(simapp.TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := app.EnterpriseKeeper.GetEnterpriseSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Uint64())

	entAccount := app.EnterpriseKeeper.GetEnterpriseAccount(ctx)
	entAccountCoins := app.BankKeeper.GetAllBalances(ctx, entAccount.GetAddress())
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromAccK := app.AccountKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := app.BankKeeper.GetAllBalances(ctx, entAccFromAccK.GetAddress())
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestGetTotalSupplyWithLockedNundRemoved(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	// should be whatever the test app was initialised with - staked nund, account balances etc.
	initialBankSupply := app.BankKeeper.GetSupply(ctx, simapp.TestDenomination)
	totalSupply := sdk.NewCoins(initialBankSupply)
	totalMinted := sdk.NewInt64Coin(simapp.TestDenomination, 0)

	testAddresses := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(simapp.RandInBetween(1000, 100000))
		amountToUnlock := int64(simapp.RandInBetween(1, 999))

		toMint := sdk.NewInt64Coin(simapp.TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(simapp.TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalSupply = totalSupply.Add(toUnlock)
		totalMinted = totalMinted.Add(toMint)

		pageReq := &query.PageRequest{
			Limit: 10,
		}
		totalSupplyDb, _, _ := app.EnterpriseKeeper.GetTotalSupplyWithLockedNundRemoved(ctx, pageReq)
		require.True(t, totalSupplyDb.IsEqual(totalSupply))

		totalMintedDb := app.BankKeeper.GetSupply(ctx, simapp.TestDenomination)
		require.True(t, totalMintedDb.IsEqual(totalMinted.Add(initialBankSupply)))
	}
}

func TestUnlockCoinsForFeesAndUsedCounter(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalUsed := int64(0)

	testAddresses := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(simapp.RandInBetween(1000, 100000))
		amountToUnlock := int64(simapp.RandInBetween(1, 999))

		toMint := sdk.NewInt64Coin(simapp.TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(simapp.TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)
		totalUsed = totalUsed + amountToUnlock

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		usedDb := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)
		require.True(t, usedDb.Amount.IsEqual(toUnlock))
		require.True(t, usedDb.Owner == addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(simapp.TestDenomination, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.IsEqual(expectedTotalUsedCoin))
}

func TestUnlockCoinsForFeesAndUsedCounterWithHalfFunds(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)

	totalUsed := int64(0)

	testAddresses := simapp.AddTestAddrs(app, ctx, 100, sdk.NewInt(10000))

	for _, addr := range testAddresses {
		amountToMint := int64(simapp.RandInBetween(1, 999))
		// fee is more than minted, to test using account's normal fund supply in addition to minted efund
		feeToPay := amountToMint * 2
		// only minted will count as "used"
		totalUsed = totalUsed + amountToMint

		toMint := sdk.NewInt64Coin(simapp.TestDenomination, amountToMint)
		fee := sdk.NewInt64Coin(simapp.TestDenomination, feeToPay)
		feeCoins := sdk.NewCoins(fee)

		_ = app.EnterpriseKeeper.MintCoinsAndLock(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addr, feeCoins)
		require.NoError(t, err)

		usedDb := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)
		// fee is 2x what was minted. Only minted should count
		require.True(t, usedDb.Amount.IsEqual(toMint))
		require.True(t, usedDb.Owner == addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(simapp.TestDenomination, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.IsEqual(expectedTotalUsedCoin))
}
