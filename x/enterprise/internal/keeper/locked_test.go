package keeper

import (
	"math/rand"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

func TestSetGetTotalLockedUnd(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	denom := TestDenomination
	amount := int64(1000)
	locked := sdk.NewInt64Coin(denom, amount)

	err := keeper.SetTotalLockedUnd(ctx, locked)
	require.NoError(t, err)

	lockedDb := keeper.GetTotalLockedUnd(ctx)

	require.True(t, lockedDb.IsEqual(locked))
	require.True(t, lockedDb.Denom == denom)
	require.True(t, lockedDb.Amount.Int64() == amount)
}

func TestSetGetLockedUndForAccount(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	testAddresses := GenerateRandomAddresses(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		denom := TestDenomination

		locked := types.NewLockedUnd(addr, denom)
		locked.Amount = sdk.NewInt64Coin(denom, amount)

		err := keeper.SetLockedUndForAccount(ctx, locked)
		require.NoError(t, err)

		lockedDb := keeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, LockedUndEqual(locked, lockedDb))

		lockedDbAmount := keeper.GetLockedUndAmountForAccount(ctx, addr)
		require.True(t, lockedDbAmount.IsEqual(locked.Amount))
	}
}

func TestIsLocked(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	denom := TestDenomination
	l0 := types.NewLockedUnd(TestAddrs[0], denom)
	l0.Amount = sdk.NewInt64Coin(denom, 0)

	l1 := types.NewLockedUnd(TestAddrs[1], denom)
	l1.Amount = sdk.NewInt64Coin(denom, 1)

	l2 := types.NewLockedUnd(TestAddrs[2], denom)
	l2.Amount = sdk.NewInt64Coin(denom, 0)

	l3 := types.NewLockedUnd(TestAddrs[3], denom)
	l3.Amount = sdk.NewInt64Coin(denom, 100)

	testCases := []struct {
		l           types.LockedUnd
		expIsLocked bool
	}{
		{l0, false},
		{l1, true},
		{l2, false},
		{l3, true},
	}

	for _, tc := range testCases {
		err := keeper.SetLockedUndForAccount(ctx, tc.l)
		require.NoError(t, err)

		isLocked := keeper.IsLocked(ctx, tc.l.Owner)
		require.Equal(t, tc.expIsLocked, isLocked)
	}
}

func TestMintCoinsAndLock(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper := createTestInput(t, false, 100)

	totalAmount := int64(0)

	testAddresses := GenerateRandomAddresses(100)

	for _, addr := range testAddresses {
		amount := int64(rand.Intn(10000) + 1)
		totalAmount = totalAmount + amount

		toMint := sdk.NewInt64Coin(TestDenomination, amount)

		err := keeper.MintCoinsAndLock(ctx, addr, toMint)
		require.NoError(t, err)

		isLocked := keeper.IsLocked(ctx, addr)
		require.True(t, isLocked)

		lockedDb := keeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(toMint))
	}

	totalLocked := sdk.NewInt64Coin(TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := keeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := keeper.GetTotalSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Int64())

	entAccount := keeper.GetEnterpriseAccount(ctx)
	entAccountCoins := entAccount.GetCoins()
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromSk := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := entAccFromSk.GetCoins()
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestUnlockCoinsForFees(t *testing.T) {
	ctx, _, keeper, _, supplyKeeper := createTestInput(t, false, 100)
	totalAmount := int64(0)

	testAddresses := GenerateRandomAddresses(100)

	for _, addr := range testAddresses {
		amountToMint := int64(RandInBetween(1000, 100000))
		amountToUnlock := int64(RandInBetween(1, 999))
		totalAmount = totalAmount + amountToMint

		toMint := sdk.NewInt64Coin(TestDenomination, amountToMint)
		toUnlock := sdk.NewInt64Coin(TestDenomination, amountToUnlock)
		toUnlockCoins := sdk.NewCoins(toUnlock)

		_ = keeper.MintCoinsAndLock(ctx, addr, toMint)

		err := keeper.UnlockCoinsForFees(ctx, addr, toUnlockCoins)
		require.NoError(t, err)

		totalAmount = totalAmount - amountToUnlock

		expectedLocked := toMint.Sub(toUnlock)

		lockedDb := keeper.GetLockedUndForAccount(ctx, addr)
		require.True(t, lockedDb.Amount.IsEqual(expectedLocked))
	}

	totalLocked := sdk.NewInt64Coin(TestDenomination, totalAmount)
	totalLockedCoins := sdk.NewCoins(totalLocked)

	totalLockedDb := keeper.GetTotalLockedUnd(ctx)
	require.True(t, totalLockedDb.IsEqual(totalLocked))

	totalSupplyDb := keeper.GetTotalSupplyIncludingLockedUnd(ctx)
	require.True(t, totalSupplyDb.Locked == totalLocked.Amount.Int64())

	entAccount := keeper.GetEnterpriseAccount(ctx)
	entAccountCoins := entAccount.GetCoins()
	require.True(t, entAccountCoins.IsEqual(totalLockedCoins))

	entAccFromSk := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	entAccFromSkCoins := entAccFromSk.GetCoins()
	require.True(t, entAccFromSkCoins.IsEqual(totalLockedCoins))
}

func TestIncrementLockedUnd(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	testAddresses := GenerateRandomAddresses(100)
	for _, addr := range testAddresses {
		totalIncr := int64(0)
		for i := 0; i < 100; i++ {
			amountToIncrement := int64(RandInBetween(1000, 100000))
			totalIncr = totalIncr + amountToIncrement
			beforeIncr := keeper.GetLockedUndAmountForAccount(ctx, addr)

			amountToIncrementCoin := sdk.NewInt64Coin(TestDenomination, amountToIncrement)
			err := keeper.incrementLockedUnd(ctx, addr, amountToIncrementCoin)
			require.NoError(t, err)

			afterInc := keeper.GetLockedUndAmountForAccount(ctx, addr)
			require.True(t, afterInc.IsEqual(beforeIncr.Add(amountToIncrementCoin)))
		}

		totalIncrCoin := sdk.NewInt64Coin(TestDenomination, totalIncr)
		totalIncrForAddr := keeper.GetLockedUndAmountForAccount(ctx, addr)
		require.True(t, totalIncrForAddr.IsEqual(totalIncrCoin))
	}
}

func TestDecrementLockedUnd(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	testAddresses := GenerateRandomAddresses(100)
	for _, addr := range testAddresses {
		totDecr := int64(0)
		initLocked := int64(RandInBetween(1000, 100000))
		initCoin := sdk.NewInt64Coin(TestDenomination, initLocked)
		locked := types.NewLockedUnd(addr, TestDenomination)
		locked.Amount = initCoin
		_ = keeper.SetLockedUndForAccount(ctx, locked)
		for i := 0; i < 100; i++ {
			amountToDecr := int64(RandInBetween(1, int(initLocked)))
			totDecr = totDecr + amountToDecr
			beforeDecr := keeper.GetLockedUndAmountForAccount(ctx, addr)
			amountToDecrCoin := sdk.NewInt64Coin(TestDenomination, amountToDecr)
			err := keeper.DecrementLockedUnd(ctx, addr, amountToDecrCoin)
			require.NoError(t, err)

			afterDecr := keeper.GetLockedUndAmountForAccount(ctx, addr)
			if !afterDecr.IsZero() {
				require.True(t, afterDecr.IsEqual(beforeDecr.Sub(amountToDecrCoin)))
			} else {
				require.True(t, afterDecr.IsZero())
			}
		}
		totalDecrCoin := sdk.NewInt64Coin(TestDenomination, totDecr)
		totalDecrForAddr := keeper.GetLockedUndAmountForAccount(ctx, addr)
		if initCoin.IsLT(totalDecrCoin) {
			require.True(t, totalDecrForAddr.IsZero())
		} else {
			require.True(t, totalDecrForAddr.IsEqual(initCoin.Sub(totalDecrCoin)))
		}
	}
}
