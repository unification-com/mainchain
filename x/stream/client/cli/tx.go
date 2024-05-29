package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/unification-com/mainchain/x/stream/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdCreateStream(),
		GetCmdClaimStream(),
	)

	return cmd
}

// GetCmdCreateStream is the CLI command for creating a new stream
func GetCmdCreateStream() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [receiver] [deposit] [flow_rate]",
		Short: "Create a new payment stream between you and the receiver wallet",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new payment stream
Example:
$ %s tx %s create und173qnkw458p646fahmd53xa45vqqvga7kyu6ryy 777000000000nund 299768 --from t1
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress()

			receiver, err := sdk.AccAddressFromBech32(args[0])

			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			flowRate, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateStream(deposit, flowRate, receiver, sender)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdClaimStream is the CLI command for claiming funds held in a stream
func GetCmdClaimStream() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [stream_id]",
		Short: "Claim funds held in a stream",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Claim funds held in a stream
Example:
$ %s tx %s claim 1 --from t1
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			receiver := clientCtx.GetFromAddress()

			streamId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimStream(streamId, receiver)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
