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

	require.True(t, lockedDb.Equal(locked))
	require.Equal(t, lockedDb.Denom, denom)
	require.Equal(t, lockedDb.Amount.Int64(), amount)
}

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

		require.Equal(t, locked.Owner, lockedDb.Owner)
		require.True(t, lockedDb.Amount.Equal(locked.Amount))

		lockedDbAmount := app.EnterpriseKeeper.GetLockedUndAmountForAccount(ctx, addr)
		require.True(t, lockedDbAmount.Equal(locked.Amount))
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
		require.True(t, lockedDb.Amount.Equal(expectedLocked))

		balanceAfter := app.BankKeeper.GetBalance(ctx, addr, sdk.DefaultBondDenom)
		require.Equal(t, expBalanceAfter, balanceAfter)
	}

	// compare total supply
	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	require.Equal(t, expTotalSupply, totalSupplyAfter)

	totalLocked := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalAmount)

	totalLockedDb := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.Equal(totalLocked))

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
		require.True(t, usedDb.Amount.Equal(toUnlock))
		require.Equal(t, usedDb.Owner, addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.Equal(expectedTotalUsedCoin))
}

func TestUnlockAndMintCoinsForFeesInsufficientFunds(t *testing.T) {
	// Branch C of UnlockAndMintCoinsForFees: locked + spendable < fee.
	// Expected behaviour: function returns nil silently, with NO state change.
	// The intent is to let the downstream DeductFeeDecorator reject the tx with a
	// proper "insufficient fee" error; the enterprise keeper should not partially
	// mint/unlock when the resulting balance still wouldn't cover the fee.

	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)

	// Random accounts have zero bank balance — guarantees spendable = 0.
	testAddresses := simapphelpers.GenerateRandomTestAccounts(10)

	denom := sdk.DefaultBondDenom
	lockedAmount := int64(100)
	feeAmount := int64(1000) // > locked (100) + spendable (0)

	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, denom)
	totalLockedBefore := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	totalSpentBefore := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)

	for _, addr := range testAddresses {
		// Lock a small amount — insufficient on its own to cover the fee.
		toLock := sdk.NewInt64Coin(denom, lockedAmount)
		err := app.EnterpriseKeeper.CreateAndLockEFUND(ctx, addr, toLock)
		require.NoError(t, err)

		balanceBefore := app.BankKeeper.GetBalance(ctx, addr, denom)
		require.True(t, balanceBefore.IsZero(), "test setup requires zero spendable balance")

		lockedBefore := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		spentBefore := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)

		// Fee exceeds locked + spendable. Branch C: no-op return.
		fee := sdk.NewCoins(sdk.NewInt64Coin(denom, feeAmount))
		err = app.EnterpriseKeeper.UnlockAndMintCoinsForFees(ctx, addr, fee)
		require.NoError(t, err, "Branch C should return nil silently")

		// Per-account state must be unchanged.
		lockedAfter := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedAfter.Amount.Equal(lockedBefore.Amount),
			"locked unchanged: before %s, after %s", lockedBefore.Amount, lockedAfter.Amount)

		spentAfter := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, addr)
		require.True(t, spentAfter.Amount.Equal(spentBefore.Amount),
			"spent unchanged: before %s, after %s", spentBefore.Amount, spentAfter.Amount)

		balanceAfter := app.BankKeeper.GetBalance(ctx, addr, denom)
		require.True(t, balanceAfter.Equal(balanceBefore),
			"bank balance unchanged: before %s, after %s", balanceBefore, balanceAfter)
	}

	// Module-wide totals: total locked should have grown by 10 * 100 = 1000
	// (from the CreateAndLockEFUND calls), but total spent and total supply
	// should be unchanged.
	expectedTotalLocked := totalLockedBefore.Add(sdk.NewInt64Coin(denom, lockedAmount*int64(len(testAddresses))))
	totalLockedAfter := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedAfter.Equal(expectedTotalLocked),
		"total locked: expected %s, got %s", expectedTotalLocked, totalLockedAfter)

	totalSpentAfter := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalSpentAfter.Equal(totalSpentBefore),
		"total spent must NOT change in Branch C: before %s, after %s", totalSpentBefore, totalSpentAfter)

	// Critically: no minting must have happened.
	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, denom)
	require.True(t, totalSupplyAfter.Equal(totalSupplyBefore),
		"total supply must NOT change in Branch C (no minting): before %s, after %s", totalSupplyBefore, totalSupplyAfter)
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
		require.True(t, usedDb.Amount.Equal(toMint))
		require.Equal(t, usedDb.Owner, addr.String())
	}

	expectedTotalUsedCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, totalUsed)

	totalUsedDb := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	require.True(t, totalUsedDb.Equal(expectedTotalUsedCoin))
}
