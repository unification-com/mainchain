package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"

	mathmod "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func TestSetGetTotalLockedUnd(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	denom := sdk.DefaultBondDenom
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	require.NoError(t, err)

	lockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)

	require.True(t, lockedDb.IsEqual(locked))
	require.True(t, lockedDb.Denom == denom)
	require.True(t, lockedDb.Amount.Int64() == amount)
}

//func TestGetTotalUnlocked(t *testing.T) {
//	app := simapphelpers.Setup(t)
//	ctx := app.BaseApp.NewContext(false)
//	simapphelpers.AddTestAddrs(app, ctx, 1, mathmod.NewInt(20000))
//
//	denom := sdk.DefaultBondDenom
//	amount := int64(1000)
//	locked := sdk.NewInt64Coin(denom, amount)
//
//	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
//	require.NoError(t, err)
//
//	totUnlocked := app.EnterpriseKeeper.GetTotalUnLockedUnd(ctx)
//	totalSupply := app.BankKeeper.GetSupply(ctx, denom)
//
//	diff := totalSupply.Sub(totUnlocked)
//
//	require.Equal(t, locked, diff)
//}

//func TestGetTotalUndSupply(t *testing.T) {
//	app := simapphelpers.Setup(t)
//	ctx := app.BaseApp.NewContext(false)
//	simapphelpers.AddTestAddrs(app, ctx, 1, mathmod.NewInt(20000))
//
//	totalSupply := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
//	totalSupplyFromEnt := app.EnterpriseKeeper.GetTotalUndSupply(ctx)
//	require.Equal(t, totalSupply, totalSupplyFromEnt)
//}

func TestSetGetLockedUndForAccount(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	testAddresses := simapphelpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		denom := sdk.DefaultBondDenom

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

	denom := sdk.DefaultBondDenom

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

func TestCreateAndLockEFUND(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	totalAmount := int64(0)

	testAddresses := simapphelpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		totalAmount = totalAmount + amount
		balanceBefore := app.BankKeeper.GetBalance(ctx, addr, sdk.DefaultBondDenom)

		toCreate := sdk.NewInt64Coin(sdk.DefaultBondDenom, amount)

		err := app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, toCreate)
		require.NoError(t, err)

		isLocked := app.EnterpriseKeeper.IsLocked(ctx, addr)
		require.True(t, isLocked)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.Equal(toCreate))

		balanceAfter := app.BankKeeper.GetBalance(ctx, addr, sdk.DefaultBondDenom)
		require.Equal(t, balanceBefore, balanceAfter)
	}

	totalLocked := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalAmount)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.Equal(totalLocked))

}

func TestUnlockAndMintCoinsForFees(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	totalAmount := int64(0)

	testAddresses := simapphelpers.GenerateRandomTestAccounts(100)

	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	expTotalSupply := totalSupplyBefore

	for _, addr := range testAddresses {
		amountToMint := int64(simapphelpers.RandInBetween(1000, 100000))
		amountToUnlock := int64(simapphelpers.RandInBetween(1, 999))
		totalAmount = totalAmount + amountToMint
		balanceBefore := app.BankKeeper.GetBalance(ctx, addr, sdk.DefaultBondDenom)

		toMint := sdk.NewInt64Coin(sdk.DefaultBondDenom, amountToMint)
		toUnlock := sdk.NewInt64Coin(sdk.DefaultBondDenom, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)
		expBalanceAfter := balanceBefore.Add(toUnlock)
		expTotalSupply = expTotalSupply.Add(toUnlock)

		_ = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalAmount = totalAmount - amountToUnlock

		expectedLocked := toMint.Sub(toUnlock)

		lockedDb := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(expectedLocked))

		balanceAfter := app.BankKeeper.GetBalance(ctx, addr, sdk.DefaultBondDenom)
		require.Equal(t, expBalanceAfter, balanceAfter)
	}

	// compare total supply
	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	require.Equal(t, expTotalSupply, totalSupplyAfter)

	totalLocked := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalAmount)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

}

func TestUnlockCoinsForFeesAndUsedCounter(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	totalUsed := int64(0)

	testAddresses := simapphelpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddresses {
		amountToMint := int64(simapphelpers.RandInBetween(1000, 100000))
		amountToUnlock := int64(simapphelpers.RandInBetween(1, 999))

		toMint := sdk.NewInt64Coin(sdk.DefaultBondDenom, amountToMint)
		toUnlock := sdk.NewInt64Coin(sdk.DefaultBondDenom, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)
		totalUsed = totalUsed + amountToUnlock

		_ = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		usedDb := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)
		require.True(t, usedDb.Amount.IsEqual(toUnlock))
		require.True(t, usedDb.Owner == addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.IsEqual(expectedTotalUsedCoin))
}

func TestUnlockCoinsForFeesAndUsedCounterWithHalfFunds(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	totalUsed := int64(0)

	testAddresses := simapphelpers.AddTestAddrs(app, ctx, 100, mathmod.NewInt(10000))

	for _, addr := range testAddresses {
		amountToMint := int64(simapphelpers.RandInBetween(1, 999))
		// fee is more than minted, to test using account's normal fund supply in addition to minted efund
		feeToPay := amountToMint * 2
		// only minted will count as "used"
		totalUsed = totalUsed + amountToMint

		toMint := sdk.NewInt64Coin(sdk.DefaultBondDenom, amountToMint)
		fee := sdk.NewInt64Coin(sdk.DefaultBondDenom, feeToPay)
		feeCoins := sdk.NewCoins(fee)

		_ = app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, toMint)

		err := app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addr, feeCoins)
		require.NoError(t, err)

		usedDb := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)
		// fee is 2x what was minted. Only minted should count
		require.True(t, usedDb.Amount.IsEqual(toMint))
		require.True(t, usedDb.Owner == addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.IsEqual(expectedTotalUsedCoin))
}
