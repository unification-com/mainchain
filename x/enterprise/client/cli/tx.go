package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

const (
	FlagPurchaseOrderStatus = "status"
	FlagPurchaser           = "purchaser"
)

func GetTxCmd() *cobra.Command {
	enterpriseTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Enterprise FUND transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	enterpriseTxCmd.AddCommand(
		GetCmdRaisePurchaseOrder(),
		GetCmdProcessPurchaseOrder(),
		GetCmdWhitelistAction(),
	)

	return enterpriseTxCmd
}

// GetCmdRegisterWrkChain is the CLI command for creating an Enterprise FUND purchase order
func GetCmdRaisePurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase [amount]",
		Short: "Raise a new Enterprise FUND purchase order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Raise a new Enterprise FUND purchase order
Example:
$ %s tx %s purchase 1000000000000%s --from wrktest
`,
				version.AppName, types.ModuleName, sdk.DefaultBondDenom,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			if amount.Denom != sdk.DefaultBondDenom {
				return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, fmt.Sprintf("denomination should be %s", sdk.DefaultBondDenom))
			}

			msg := types.NewMsgUndPurchaseOrder(from, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdProcessPurchaseOrder is the CLI command for processing an Enterprise FUND purchase order
func GetCmdProcessPurchaseOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process [purchase_order_id] [decision]",
		Short: "Process an Enterprise FUND purchase order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Process an Enterprise FUND purchase order
Example:
$ %s tx %s process 24 accept --from ent
$ %s tx %s process 24 reject --from ent
`,
				version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			// validate that the proposal id is a uint
			purchaseOrderId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("purchase_order_id %s not a valid int, please input a valid purchase_order_id", args[0])
			}

			decision, err := types.PurchaseOrderStatusFromString(args[1])
			if err != nil {
				return err
			}

			if !types.ValidPurchaseOrderAcceptRejectStatus(decision) {
				return errorsmod.Wrap(types.ErrInvalidDecision, "decision should be accept or reject")
			}

			msg := types.NewMsgProcessUndPurchaseOrder(purchaseOrderId, decision, from)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdWhitelistAction is the CLI command for adding/removing addresses from the purchase order whitelist
func GetCmdWhitelistAction() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist [action] [address]",
		Short: "Add/Remove an address from the enterprise purchase order whitelist",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Add/Remove an address from the enterprise purchase order whitelist
Example:
$ %s tx %s whitelist add und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --from ent
$ %s tx %s whitelist remove und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy --from ent
`,
				version.AppName, types.ModuleName, version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			from := clientCtx.GetFromAddress()

			address, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			action, err := types.WhitelistActionFromString(args[0])
			if err != nil {
				return err
			}

			if !types.ValidWhitelistAction(action) {
				return errorsmod.Wrap(types.ErrInvalidWhitelistAction, "action should be add or remove")
			}

			msg := types.NewMsgWhitelistAddress(address, action, from)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
