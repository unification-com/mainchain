package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	undtypes "github.com/unification-com/mainchain/types"
	"strings"
)

func GetDenomConversionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "convert [amount] [from_denom] [to_denom]",
		Short: "convert between FUND denominations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`convert between FUND denominations'
Example:
$ %s convert 24 fund nund
`,
				version.AppName,
			),
		),

		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := undtypes.ConvertUndDenomination(args[0], args[1], args[2])

			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprintf("%s%s = %s\n", args[0], args[1], res))
		},
	}
}
