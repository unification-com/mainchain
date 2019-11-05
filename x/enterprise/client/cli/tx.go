package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
	"strings"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	enterpriseTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Enterprise UND transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	enterpriseTxCmd.AddCommand(client.PostCommands(
		GetCmdRaisePurchaseOrder(cdc),
	)...)

	return enterpriseTxCmd
}


// GetCmdRegisterWrkChain is the CLI command for sending a RegisterWrkChain transaction
func GetCmdRaisePurchaseOrder(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "purchase [amount]",
		Short: "Raise a new Enterprise UND purchase order",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Raise a new Enterprise UND purchase order
Example:
$ %s tx %s purchase 1000000000000nund --from wrktest
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			amount, err := sdk.ParseCoin(args[0])
			if err != nil {
				return err
			}

			// todo - check denom is nund
			msg := types.NewMsgRaiseUndPurchaseOrder(cliCtx.GetFromAddress(), amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
