package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpdateName this will be changed with each new release that requires migrations
const UpdateName = "2-grog"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// second upgrade 2-grog upgrades ibc-go to v3.4.0
	// See https://github.com/cosmos/ibc-go/blob/v3.4.0/docs/migrations/support-denoms-with-slashes.md
	// Note: Cosmos SDK upgrade to v0.45.13 does not have any migrations
	app.UpgradeKeeper.SetUpgradeHandler(UpdateName, app.upgradeHandler)

	_, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
}

func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return app.mm.RunMigrations(ctx, app.configurator, vm)
}
