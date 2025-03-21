package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	"strings"
)

const FlagDenom = "denom"

// GetTotalSupplyCmd used in place of Cosmos SDK 'bank total' command
func GetTotalSupplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply",
		Short: "Query the total supply of coins of the chain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total supply of coins that are held by accounts in the chain.

NOTE: Use instead of the default Cosmos SDK total supply query, to subtract locked enterprise FUND from the result

Example:
 $ %s query supply
To query for the total supply of a specific coin denomination use:
 $ %s query supply --denom=[denom]
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return err
			}

			queryClient := enttypes.NewQueryClient(clientCtx)

			if denom == "" {
				res, err := queryClient.TotalSupply(cmd.Context(), &enttypes.QueryTotalSupplyRequest{})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.SupplyOf(cmd.Context(), &enttypes.QuerySupplyOfRequest{Denom: denom})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Amount)
		},
	}
	cmd.Flags().String(FlagDenom, "", "The specific balance denomination to query for")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryTotalSupplyOverrideBankDefault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply",
		Short: "Query the total supply of coins of the chain (x/enterprise override)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total supply of coins that are held by accounts in the chain.

Example:
  $ %s query bank total
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := enttypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.TotalSupply(ctx, &enttypes.QueryTotalSupplyRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagDenom, "", "The specific balance denomination to query for")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all supply totals")

	return cmd
}

func GetCmdQueryTotalSupplyOfOverrideBankDefault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-supply-of [denom]",
		Short: "Query the supply of a single coin denom (x/enterprise override)",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query total supply of coins that are held by accounts in the chain.

Example:
  $ %s query bank total %s
`,
				version.AppName, sdk.DefaultBondDenom,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			denom := args[0]

			queryClient := enttypes.NewQueryClient(clientCtx)
			ctx := cmd.Context()

			if err != nil {
				return err
			}
			if denom == "" {
				return fmt.Errorf("must provide denomination")
			}

			res, err := queryClient.SupplyOf(ctx, &enttypes.QuerySupplyOfRequest{Denom: denom})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Amount)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
