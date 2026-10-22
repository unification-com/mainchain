package app

import (
	"context"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName this will be changed with each new release that requires migrations
const UpgradeName = "8-vaxildan"

// see https://docs.cosmos.network/main/build/migrations/chain-upgrade-guide-044
func (app *App) registerUpgradeHandlers() {

	// 8-vaxildan
	// 1. updates Cosmos SDK to v0.54.x, IBC-go to v11.x, CometBFT to v0.39.x.
	// 2. removes x/group module wiring (mainnet verified empty).
	// 3. removes x/circuit module wiring (mainnet verified empty).
	// 4. widens x/stream keys with a denom coordinate; registers v1→v2 store migration.

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// RunMigrations picks up the stream module's v1→v2 migration
			// (registered in x/stream/module.go via cfg.RegisterMigration) — no
			// explicit call needed here; it fires automatically when
			// fromVM[streamtypes.ModuleName] == 1.
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// Store keys for modules removed in 8-vaxildan. Hardcoded as raw strings
		// because the module packages were deleted in stage 1 of vaxildan, so the
		// `group.StoreKey` / `circuittypes.StoreKey` constants are no longer importable.
		// The string values match the upstream constants at the time of removal.
		storeUpgrades := storetypes.StoreUpgrades{
			Added:   []string{},
			Deleted: []string{"group", "circuit"},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
