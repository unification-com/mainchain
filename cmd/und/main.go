package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/store"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/viper"
	undtypes "github.com/unification-com/mainchain/types"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/unification-com/mainchain/app"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	config.SetCoinType(undtypes.CoinType)
	config.SetFullFundraiserPath(undtypes.HdWalletPath)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "und",
		Short:             "Unification Mainchain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, cdc),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		flags.NewCompletionCmd(rootCmd, true),
		debug.Cmd(cdc),
		DumpDataCmd(ctx, cdc, dumpBeaconOrWrkchainData),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "UND", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")

	rootCmd.PersistentFlags().IntSlice(undtypes.FlagExportIncludeWrkchainData, []int{},
		"Comma separated list of WRKChain IDs for which data will also be exported")
	rootCmd.PersistentFlags().IntSlice(undtypes.FlagExportIncludeBeaconData, []int{},
		"Comma separated list of BEACON IDs for which data will also be exported")

	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {

	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	return app.NewMainchainApp(logger, db, traceStore, true, invCheckPeriod,
		viper.GetString(flags.FlagHome),
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		undApp := app.NewMainchainApp(logger, db, traceStore, false, uint(1), viper.GetString(flags.FlagHome))
		err := undApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return undApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	undApp := app.NewMainchainApp(logger, db, traceStore, true, uint(1), viper.GetString(flags.FlagHome))

	return undApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func dumpBeaconOrWrkchainData(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, what string, id uint64,
) (json.RawMessage, error) {

	if height != -1 {
		undApp := app.NewMainchainApp(logger, db, traceStore, false, uint(1), viper.GetString(flags.FlagHome))
		err := undApp.LoadHeight(height)
		if err != nil {
			return nil, err
		}
		return undApp.DumpWrkchainOrBeaconData(what, id)
	}

	undApp := app.NewMainchainApp(logger, db, traceStore, true, uint(1), viper.GetString(flags.FlagHome))

	return undApp.DumpWrkchainOrBeaconData(what, id)
}
