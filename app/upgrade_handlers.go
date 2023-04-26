package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpdateName this will be changed with each new release that requires migrations
const UpdateName = "3-keyleth"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// third upgrade 3-keyleth upgrades ibc-go to v5.2.0
	// and Cosmos SDK to v0.46.x
	app.UpgradeKeeper.SetUpgradeHandler(UpdateName, app.upgradeHandler)

	_, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
}

func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return app.mm.RunMigrations(ctx, app.configurator, vm)
}
