package cmd

import (
	"cosmossdk.io/client/v2/autocli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtxconfig "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/ibc-go/v8/testing/simapp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/unification-com/mainchain/app"
)

//var ChainID string

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() *cobra.Command {
	// we "pre"-instantiate the application for getting the injected/configured encoding configuration
	// note, this is not necessary when using app wiring, as depinject can be directly used (see root_v2.go)
	tempApp := app.NewApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome))
	//encodingConfig := params.EncodingConfig{
	//	InterfaceRegistry: tempApp.InterfaceRegistry(),
	//	Codec:             tempApp.AppCodec(),
	//	TxConfig:          tempApp.TxConfig(),
	//	Amino:             tempApp.LegacyAmino(),
	//}

	initClientCtx := client.Context{}.
		WithCodec(tempApp.AppCodec()).
		WithInterfaceRegistry(tempApp.InterfaceRegistry()).
		//WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(tempApp.LegacyAmino()).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(simapp.DefaultNodeHome).
		WithViper("") // In simapp, we don't use any prefix for env variables.

	rootCmd := &cobra.Command{
		Use:           app.Name,
		Short:         "Unification Mainchain App",
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx = initClientCtx.WithCmdContext(cmd.Context())
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				enabledSignModes := append(tx.DefaultSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: authtxconfig.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfig, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}

				initClientCtx = initClientCtx.WithTxConfig(txConfig)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customCMTConfig := initCometBFTConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customCMTConfig)
		},
	}

	initRootCmd(rootCmd, tempApp.BasicModuleManager, tempApp.GetTxConfig())

	autoCliOpts, err := enrichAutoCliOpts(tempApp.AutoCliOpts(), initClientCtx)
	if err != nil {
		panic(err)
	}

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	// ToDo - need to remove this once enterprise module's minting process is modified
	rootCmd = overrideBankQuery(rootCmd)

	return rootCmd
}

func enrichAutoCliOpts(autoCliOpts autocli.AppOptions, clientCtx client.Context) (autocli.AppOptions, error) {
	autoCliOpts.AddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	autoCliOpts.ValidatorAddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	autoCliOpts.ConsensusAddressCodec = addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix())

	var err error
	clientCtx, err = config.ReadFromClientConfig(clientCtx)
	if err != nil {
		return autocli.AppOptions{}, err
	}

	autoCliOpts.ClientCtx = clientCtx

	return autoCliOpts, nil
}

//func NewRootCmdOLD() (*cobra.Command, params.EncodingConfig) {
//	encodingConfig := app.MakeEncodingConfig()
//	initClientCtx := client.Context{}.
//		WithCodec(encodingConfig.Codec).
//		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
//		WithTxConfig(encodingConfig.TxConfig).
//		WithLegacyAmino(encodingConfig.Amino).
//		WithInput(os.Stdin).
//		WithAccountRetriever(types.AccountRetriever{}).
//		WithHomeDir(app.DefaultNodeHome).
//		WithViper("")
//
//	rootCmd := &cobra.Command{
//		Use:   app.Name,
//		Short: "Unification Mainchain App",
//		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
//			// set the default command outputs
//			cmd.SetOut(cmd.OutOrStdout())
//			cmd.SetErr(cmd.ErrOrStderr())
//
//			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
//			if err != nil {
//				return err
//			}
//
//			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
//			if err != nil {
//				return err
//			}
//
//			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
//				return err
//			}
//
//			customAppTemplate, customAppConfig := initAppConfig()
//			customTMConfig := initTendermintConfig()
//
//			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
//		},
//	}
//
//	initRootCmd(rootCmd, encodingConfig)
//	overwriteFlagDefaults(rootCmd, map[string]string{
//		flags.FlagChainID: ChainID,
//	})
//
//	return rootCmd, encodingConfig
//}

// initTendermintConfig helps to override default Tendermint Config values.
// return tmcfg.DefaultConfig if no custom configuration is required for the application.
//func initTendermintConfig() *tmcfg.Config {
//	cfg := tmcfg.DefaultConfig()
//
//	// these values put a higher strain on node memory
//	// cfg.P2P.MaxNumInboundPeers = 100
//	// cfg.P2P.MaxNumOutboundPeers = 40
//
//	return cfg
//}

//func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
//
//	// Set config for prefixes
//	app.SetConfig()
//
//	rootCmd.AddCommand(
//		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
//		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
//		debug.Cmd(),
//		config.Cmd(),
//		pruning.Cmd(newApp, app.DefaultNodeHome),
//		snapshot.Cmd(newApp),
//	)
//
//	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)
//
//	// add keybase, auxiliary RPC, query, genesis, and tx child commands
//	rootCmd.AddCommand(
//		rpc.StatusCommand(),
//		genesisCommand(encodingConfig),
//		queryCommand(),
//		txCommand(),
//		keys.Commands(app.DefaultNodeHome),
//		GetDenomConversionCmd(),
//	)
//
//}

//func addModuleInitFlags(startCmd *cobra.Command) {
//	crisis.AddModuleInitFlags(startCmd)
//}

// genesisCommand builds genesis-related `simd genesis` command. Users may provide application specific commands as a parameter
//func genesisCommand(encodingConfig params.EncodingConfig, cmds ...*cobra.Command) *cobra.Command {
//	cmd := genutilcli.GenesisCoreCommand(encodingConfig.TxConfig, app.ModuleBasics, app.DefaultNodeHome)
//
//	for _, sub_cmd := range cmds {
//		cmd.AddCommand(sub_cmd)
//	}
//	return cmd
//}

//func queryCommand() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:                        "query",
//		Aliases:                    []string{"q"},
//		Short:                      "Querying subcommands",
//		DisableFlagParsing:         true,
//		SuggestionsMinimumDistance: 2,
//		RunE:                       client.ValidateCmd,
//	}
//
//	cmd.AddCommand(
//		rpc.QueryEventForTxCmd(),
//		authcmd.GetAccountCmd(),
//		rpc.ValidatorCommand(),
//		rpc.BlockCommand(),
//		authcmd.QueryTxsByEventsCmd(),
//		authcmd.QueryTxCmd(),
//	)
//
//	app.ModuleBasics.AddQueryCommands(cmd)
//
//	cmd.AddCommand(
//		GetTotalSupplyCmd(),
//	)
//
//	// replace bank total command with Enterprise version to get correct total supply
//	// since the Enterprise module's balance is not part of the circulating supply until
//	// it is used to pay for BEACON/WrkChain Tx fees
//	origBankCmd, _, _ := cmd.Find([]string{"bank"})
//	origTotalSupplyCmd, _, _ := origBankCmd.Find([]string{"total"})
//
//	// remove "bank" command from "query" command
//	cmd.RemoveCommand(origBankCmd)
//	// remove "total" command from "bank" cmd
//	origBankCmd.RemoveCommand(origTotalSupplyCmd)
//	// add Enterprise version of "total" command to "bank" cmd
//	origBankCmd.AddCommand(GetCmdQueryTotalSupplyOverrideBankDefault())
//	// re-add "bank" command to "query"
//	cmd.AddCommand(origBankCmd)
//
//	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")
//
//	return cmd
//}

//func txCommand() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:                        "tx",
//		Short:                      "Transactions subcommands",
//		DisableFlagParsing:         false,
//		SuggestionsMinimumDistance: 2,
//		RunE:                       client.ValidateCmd,
//	}
//
//	cmd.AddCommand(
//		authcmd.GetSignCommand(),
//		authcmd.GetSignBatchCommand(),
//		authcmd.GetMultiSignCommand(),
//		authcmd.GetMultiSignBatchCmd(),
//		authcmd.GetValidateSignaturesCommand(),
//		authcmd.GetBroadcastCommand(),
//		authcmd.GetEncodeCommand(),
//		authcmd.GetDecodeCommand(),
//		authcmd.GetAuxToFeeCommand(),
//	)
//
//	app.ModuleBasics.AddTxCommands(cmd)
//	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")
//
//	return cmd
//}

//type appCreator struct {
//	encCfg params.EncodingConfig
//}

func overwriteFlagDefaults(c *cobra.Command, defaults map[string]string) {
	set := func(s *pflag.FlagSet, key, val string) {
		if f := s.Lookup(key); f != nil {
			f.DefValue = val
			f.Value.Set(val)
		}
	}
	for key, val := range defaults {
		set(c.Flags(), key, val)
		set(c.PersistentFlags(), key, val)
	}
	for _, c := range c.Commands() {
		overwriteFlagDefaults(c, defaults)
	}
}
