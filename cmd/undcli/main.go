package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/unification-com/mainchain/client/lcd"
	undtypes "github.com/unification-com/mainchain/types"
	entrest "github.com/unification-com/mainchain/x/enterprise/client/rest"

	"github.com/unification-com/mainchain/app"

	"github.com/unification-com/mainchain/x/enterprise"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	config.SetCoinType(undtypes.CoinType)
	config.SetFullFundraiserPath(undtypes.HdWalletPath)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "undcli",
		Short: "Unification Mainchain CLI",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of UND Mainchain node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
		denomConversion(cdc),
	)

	executor := cli.PrepareMainCmd(rootCmd, "UND", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	entrest.RegisterAuthAccountOverride(rs.CliCtx, rs.Mux)
	entrest.RegisterTotalSupplyOverride(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	RegisterQueryRestApiEndpoints(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		//authcmd.GetAccountCmd(cdc),
		GetAccountWithLockedCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
		GetTotalSupplyWithLockedCmd(cdc),
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

func denomConversion(cdc *amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "convert [amount] [from_denom] [to_denom]",
		Short: "convert between UND denominations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`convert between UND denominations'
Example:
$ %s convert 24 und nund
`,
				version.ClientName,
			),
		),

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := undtypes.ConvertUndDenomination(args[0], args[1], args[2])

			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cliCtx.Output, "%s%s = %s\n", args[0], args[1], res)

			return nil
		},
	}
}

func GetAccountWithLockedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			accGetter := auth.NewAccountRetriever(cliCtx)

			lockedUndGetter := enterprise.NewLockedUndRetriever(cliCtx)

			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			acc, err := accGetter.GetAccount(key)
			if err != nil {
				return err
			}

			lockedUnd, err := lockedUndGetter.GetLockedUnd(key)
			if err != nil {
				return err
			}

			// todo - this is a bit hackey
			accountWithLocked := undtypes.NewAccountWithLocked()
			entUnd := undtypes.NewEnterpriseUnd()

			entUnd.Locked = lockedUnd.Amount
			entUnd.Available = acc.GetCoins().Add(sdk.NewCoins(lockedUnd.Amount))

			accountWithLocked.Account = acc
			accountWithLocked.Enterprise = entUnd

			return cliCtx.PrintOutput(accountWithLocked)
		},
	}

	return flags.GetCommands(cmd)[0]
}

func GetTotalSupplyWithLockedCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply",
		Short: "Query total supply including locked enterprise UND",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total UND supply, including locked and unlocked

Returns three values:

locked
------
total UND locked through Enterprise purchases.
This UND is only available to pay WRKChain/BEACON fees
and cannot be used for transfers or staking/delegation

amount
--------
Liquid UND in active circulation, which can be used for 
transfers, staking etc. It is the
LOCKED amount subtracted from TOTAL_SUPPLY

total_supply
------------
The total amount of UND currently on the chain, including locked UND

Example:
$ %s query supply
`,
				version.ClientName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			totalSupplyGetter := enterprise.NewTotalSupplyRetriever(cliCtx)

			totalSupply, err := totalSupplyGetter.GetTotalSupply()
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(totalSupply)

		},
	}

	return flags.GetCommands(cmd)[0]
}
