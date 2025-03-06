package app

import (
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	//"github.com/cosmos/ibc-go/v7/modules/core/exported"
	//ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	streamtypes "github.com/unification-com/mainchain/x/stream/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

// UpgradeName this will be changed with each new release that requires migrations
const UpgradeName = "5-pike"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// 5-pike updates Cosmos SDK to v0.47.x and IBC to v7.5.0

	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable() //nolint:staticcheck
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable() //nolint:staticcheck
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable() //nolint:staticcheck
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable() //nolint:staticcheck
		case beacontypes.ModuleName:
			keyTable = beacontypes.ParamKeyTable() //nolint:staticcheck
		case wrkchaintypes.ModuleName:
			keyTable = wrkchaintypes.ParamKeyTable() //nolint:staticcheck
		case enttypes.ModuleName:
			keyTable = enttypes.ParamKeyTable() //nolint:staticcheck
		case streamtypes.ModuleName:
			keyTable = streamtypes.ParamKeyTable() //nolint:staticcheck
		}

		if !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// Migrate Tendermint consensus parameters from x/params module to a dedicated x/consensus module.
			baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

			// Note: this migration is optional,
			// You can include x/gov proposal migration documented in [UPGRADING.md](https://github.com/cosmos/cosmos-sdk/blob/main/UPGRADING.md)

			// IBC
			// prune expired tendermint consensus states to save storage space
			// see https://github.com/cosmos/ibc-go/blob/main/docs/docs/05-migrations/08-v6-to-v7.md#chains
			//_, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, app.AppCodec(), app.IBCKeeper.ClientKeeper)
			//if err != nil {
			//	return nil, err
			//}

			// explicitly update the IBC 02-client params, adding the localhost client type
			// see https://github.com/cosmos/ibc-go/blob/main/docs/docs/05-migrations/09-v7-to-v7_1.md#chains
			params := app.IBCKeeper.ClientKeeper.GetParams(ctx)
			params.AllowedClients = append(params.AllowedClients, exported.Localhost)
			app.IBCKeeper.ClientKeeper.SetParams(ctx, params)

			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	// add new modules
	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				consensustypes.ModuleName,
				crisistypes.ModuleName,
				grouptypes.ModuleName,
				streamtypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

//func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
//	return app.ModuleManager.RunMigrations(ctx, app.configurator, vm)
//}
