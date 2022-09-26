package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpdateNameTn this will be changed with each new release that requires migrations
// "1-ibc" was used for FUND-TestNet-2
const UpdateNameTn = "1-ibc"

// UpdateNameMn this will be changed with each new release that requires migrations
// "1-ibc_mn" will be used for FUND-MainNet-2, as this version of und onwards contains a fix for the
// enterprise module's PO decision handler that requires all validators and nodes to upgrade to.
const UpdateNameMn = "1-init_ibc"

// see https://docs.cosmos.network/v0.45/migrations/chain-upgrade-guide-044.html
func (app *App) registerUpgradeHandlers() {

	// first upgrade 1-ibc, integrates IBC - FUND-TestNet-2 & DevNets
	app.UpgradeKeeper.SetUpgradeHandler(UpdateNameTn, app.upgradeHandler)
	// first upgrade 1-ibc, integrates IBC - FUND-MainNet-2
	app.UpgradeKeeper.SetUpgradeHandler(UpdateNameMn, app.upgradeHandler)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	// add new modules in 1-ibc upgrade for both FUND-TestNet-2/DevNets and FUND-MainNet-2
	if (upgradeInfo.Name == UpdateNameTn || upgradeInfo.Name == UpdateNameMn) && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{"authz", "feegrant", "capability", "ibc", "transfer"},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

func (app *App) upgradeHandler(ctx sdk.Context, plan upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
	// 1st-time running in-store migrations, using 1 as fromVersion to
	// avoid running InitGenesis.
	fromVM := map[string]uint64{
		"auth":         1,
		"bank":         1,
		"crisis":       1,
		"distribution": 1,
		"evidence":     1,
		"gov":          1,
		"params":       1,
		"slashing":     1,
		"staking":      1,
		"upgrade":      1,
		"vesting":      1,
		"genutil":      1,
		"enterprise":   1,
		"beacon":       1,
		"wrkchain":     1,
	}

	// Staking params - set HistoricalEntries to 10,000
	// This is required for IBC relayers to work
	stParams := app.StakingKeeper.GetParams(ctx)
	stParams.HistoricalEntries = 10000
	app.StakingKeeper.SetParams(ctx, stParams)

	// set BankKeeper's denom metadata for nund/FUND
	nundMeta := banktypes.Metadata{
		Description: "The native token of FUND mainchain.",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "nund",
				Exponent: 0,
				Aliases:  []string{"nanofund"},
			},
			{
				Denom:    "fund",
				Exponent: 9,
				Aliases:  []string{"FUND"},
			},
		},
		Base:    "nund",
		Display: "fund",
		Symbol:  "FUND",
		Name:    "FUND",
	}

	app.BankKeeper.SetDenomMetaData(ctx, nundMeta)

	return app.mm.RunMigrations(ctx, app.configurator, fromVM)
}
