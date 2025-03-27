package app_test

import (
	mathmod "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	"math/rand"
	"testing"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
)

func TestUpgradeHandlerBurn(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAccs := simapphelpers.AddTestAddrsIncremental(app, ctx, 100, mathmod.NewInt(1000000000000))

	for _, acc := range testAccs {
		// simulate previous method for creating eFUND (previously MintCoinsAndLock)
		rnd := rand.Int63n(1000000000000) + 500000000000
		toLock := sdk.NewInt64Coin(sdk.DefaultBondDenom, rnd)
		toLockCoins := sdk.NewCoins(toLock)
		app.BankKeeper.MintCoins(ctx, enttypes.ModuleName, toLockCoins)
		app.BankKeeper.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, acc, toLockCoins)
		app.BankKeeper.DelegateCoinsFromAccountToModule(ctx, acc, enttypes.ModuleName, toLockCoins)
		lockedUnd := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, acc)
		lockedUnd.Amount = lockedUnd.Amount.Add(toLock)

		_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

		totalLocked := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
		totalLockedAdd := totalLocked.Add(toLock)
		_ = app.EnterpriseKeeper.SetTotalLockedUnd(ctx, totalLockedAdd)

	}

	for _, acc := range testAccs {
		// simulate previous method for spending eFUND (previously UnlockCoinsForFees)
		rnd := rand.Int63n(500000000000)
		toUnLock := sdk.NewInt64Coin(sdk.DefaultBondDenom, rnd)
		toUnlockCoins := sdk.NewCoins(toUnLock)
		_ = app.BankKeeper.UndelegateCoinsFromModuleToAccount(ctx, enttypes.ModuleName, acc, toUnlockCoins)

		lockedUnd := app.EnterpriseKeeper.GetLockedUndForAccount(ctx, acc)
		lockedUnd.Amount = lockedUnd.Amount.Sub(toUnLock)
		_ = app.EnterpriseKeeper.SetLockedUndForAccount(ctx, lockedUnd)

		totalLocked := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
		totalLockedSub := totalLocked.Sub(toUnLock)
		_ = app.EnterpriseKeeper.SetTotalLockedUnd(ctx, totalLockedSub)

		spentEFUND := app.EnterpriseKeeper.GetSpentEFUNDForAccount(ctx, acc)
		spentEFUND.Amount = spentEFUND.Amount.Add(toUnLock)
		_ = app.EnterpriseKeeper.SetSpentEFUNDForAccount(ctx, spentEFUND)

		totalSpent := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
		newTotalUsed := totalSpent.Add(toUnLock)
		_ = app.EnterpriseKeeper.SetTotalSpentEFUND(ctx, newTotalUsed)
	}

	totalLockedBefore := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	legacyTotalSupply := totalSupplyBefore.Sub(totalLockedBefore) // simulate old method for overriding bank/supply query
	modAccBalanceBefore := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(enttypes.ModuleName), sdk.DefaultBondDenom)
	totalSpentBefore := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	fmt.Println("totalLockedBefore   :", totalLockedBefore.String())
	fmt.Println("totalSupplyBefore   :", totalSupplyBefore.String())
	fmt.Println("legacyTotalSupply   :", legacyTotalSupply.String())
	fmt.Println("totalSpentBefore    :", totalSpentBefore.String())
	fmt.Println("modAccBalanceBefore :", modAccBalanceBefore.String())

	app.BurnEnterpriseAccCoins(ctx)

	totalLockedAfter := app.EnterpriseKeeper.GetTotalLockedUnd(ctx)
	totalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	modAccBalanceAfter := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(enttypes.ModuleName), sdk.DefaultBondDenom)
	totalSpentAfter := app.EnterpriseKeeper.GetTotalSpentEFUND(ctx)
	fmt.Println("totalLockedAfter    :", totalLockedAfter.String())
	fmt.Println("totalSupplyAfter    :", totalSupplyAfter.String())
	fmt.Println("totalSpentAfter     :", totalSpentAfter.String())
	fmt.Println("modAccBalanceAfter  :", modAccBalanceAfter.String())

	require.Equal(t, totalLockedBefore, totalLockedAfter)                           // should be no change
	require.Equal(t, totalSpentBefore, totalSpentAfter)                             // should be no change
	require.Equal(t, legacyTotalSupply, totalSupplyAfter)                           // should be equal
	require.Equal(t, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), modAccBalanceAfter) // should be zero
}
