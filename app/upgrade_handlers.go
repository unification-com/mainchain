package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpdateName this will be changed with each new release that requires migrations
const UpdateName = "4-percival"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// 4-percival is a hotfix and updates Cosmos SDK to v0.46.16
	app.UpgradeKeeper.SetUpgradeHandler(UpdateName, app.upgradeHandler)

	_, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}
}

func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
	return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
}
