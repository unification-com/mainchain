package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	undtypes "github.com/unification-com/mainchain-cosmos/types"
	"os"
	"path"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/unification-com/mainchain-cosmos/app"

	"github.com/unification-com/mainchain-cosmos/x/enterprise"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "undcli",
		Short: "Unification Mainchain CLI",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of UND Mainchain node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		client.LineBreak,
		keys.Commands(),
		client.LineBreak,
		version.Cmd,
		client.NewCompletionCmd(rootCmd, true),
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
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
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
		client.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		client.LineBreak,
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
		client.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		client.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		client.LineBreak,
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
	if err := viper.BindPFlag(client.FlagChainID, cmd.PersistentFlags().Lookup(client.FlagChainID)); err != nil {
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

			_, _ = fmt.Fprintf(cliCtx.Output, "%s %s = %s\n", args[0], args[1], res)

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

			if lockedUnd.Amount.IsPositive() {
				// todo - this is a bit hackey
				accountWithLocked := undtypes.NewAccountWithLocked()
				entUnd := undtypes.NewEnterpriseUnd()

				entUnd.Locked = lockedUnd.Amount
				entUnd.Available = acc.GetCoins().Add(sdk.NewCoins(lockedUnd.Amount))

				accountWithLocked.Account = acc
				accountWithLocked.Enterprise = entUnd

				return cliCtx.PrintOutput(accountWithLocked)
			}

			return cliCtx.PrintOutput(acc)

		},
	}

	return flags.GetCommands(cmd)[0]
}
