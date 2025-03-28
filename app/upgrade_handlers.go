package app

import (
	"context"
	storetypes "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
)

// UpgradeName this will be changed with each new release that requires migrations
const UpgradeName = "6-scanlan"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// 6-scanlan
	// 1. updates Cosmos SDK to v0.50.x and IBC to v8
	// 2. upgrades enterprise module to consensus v4 (see notes for BurnEnterpriseAccCoins method)

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Starting module migrations...")
			versionMap, err := app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)

			if err != nil {
				return nil, err
			}

			err = app.BurnEnterpriseAccCoins(ctx)

			if err != nil {
				return nil, err
			}

			return versionMap, err
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				circuittypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// BurnEnterpriseAccCoins migrates the enterprise module from the previous mint/delegate strategy
// to the v4 mints at the point of use strategy.
//
// The legacy (<v4) method would mint all eFUND in a processed purchase order, send it to the eFUND owner account
// then delegate the entire amount to the enterprise module account. When eFUND was used to pay for WrkChain/BEACON fees
// the enterprise anter handler would then undelegate the specified amount from the module account to the eFUND owner
// account, so that it could be used as fee payment. However, this meant that eFUND was included in the Bank module's
// calculation, since it was indeed minted (albeit allocated to the enterprise module account, and therefore not
// available in the general supply). This meant that the enterprise module needed to override any queries to the
// bank module and subtract eFUND from the supply total returned in order to accuratley reflect the actual total supply.
//
// With v4, however, minting FUND from eFUND is handled directly by the ante, and therefore at the point of usage.
// If the fee for a // BEACON/WrkChain is 1 FUND, then only 1 eFUND is minted as FUND in order for the fee to be paid
// (as long as the enterprise account owner has suffiient eFUND to mint).
// This means that bank query overrides are no longer required, and the vanilla bank query for Supply and SupplyOf can
// be used without modification.
//
// For the migration from v3 to v4, all minted FUND delegated to the enterprise module account is burned
func (app *App) BurnEnterpriseAccCoins(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// burn enterprise module account balance
	burnModule := govtypes.ModuleName // need to use another module with burn permissions
	sdkCtx.Logger().Info("Burn enterprise module account balance")
	totalLockedBefore := app.EnterpriseKeeper.GetTotalLockedUnd(sdkCtx)
	totalSupplyBefore := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	legacyActualTotalSupply := totalSupplyBefore.Sub(totalLockedBefore)
	entModAccBalanceBefore := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(enttypes.ModuleName), sdk.DefaultBondDenom)
	burnModAccBalanceBefore := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(burnModule), sdk.DefaultBondDenom)

	sdkCtx.Logger().Info(fmt.Sprintf("totalLockedBefore %s", totalLockedBefore.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("totalSupplyBefore %s", totalSupplyBefore.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("legacyActualTotalSupply %s", legacyActualTotalSupply.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("entModAccBalanceBefore %s", entModAccBalanceBefore.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("burnModAccBalanceBefore %s", burnModAccBalanceBefore.String()))

	lockedEfund := app.EnterpriseKeeper.GetAllLockedUnds(sdkCtx)

	// 1. loop through enterprise accounts
	for _, l := range lockedEfund {
		// 2. get locked amount for account
		lBalance := l.Amount
		lOwner := l.Owner
		ownerAcc, _ := sdk.AccAddressFromBech32(lOwner)

		// 3. undelegate from module to account
		lockedCoins := sdk.NewCoins(lBalance)
		err := app.BankKeeper.UndelegateCoinsFromModuleToAccount(ctx, enttypes.ModuleName, ownerAcc, lockedCoins)
		if err != nil {
			return err
		}
		// 4. send from account to module
		err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, ownerAcc, burnModule, lockedCoins)
		if err != nil {
			return err
		}
		// 5. burn
		err = app.BankKeeper.BurnCoins(ctx, burnModule, lockedCoins)
		if err != nil {
			return err
		}
	}

	totalLockedAfter := app.EnterpriseKeeper.GetTotalLockedUnd(sdkCtx)
	ActualTotalSupplyAfter := app.BankKeeper.GetSupply(ctx, sdk.DefaultBondDenom)
	entModAccBalanceAfter := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(enttypes.ModuleName), sdk.DefaultBondDenom)
	burnModAccBalanceAfter := app.BankKeeper.GetBalance(ctx, app.AccountKeeper.GetModuleAddress(burnModule), sdk.DefaultBondDenom)
	sdkCtx.Logger().Info(fmt.Sprintf("totalLockedAfter %s", totalLockedAfter.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("ActualTotalSupplyAfter %s", ActualTotalSupplyAfter.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("entModAccBalanceAfter %s", entModAccBalanceAfter.String()))
	sdkCtx.Logger().Info(fmt.Sprintf("burnModAccBalanceAfter %s", burnModAccBalanceAfter.String()))

	return nil
}

//func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
//	return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
//}
