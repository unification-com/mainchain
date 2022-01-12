package cmd

//func GetBankQueryOverrideCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:                        banktypes.ModuleName,
//		Short:                      "Querying commands for the bank module with enterprise FUND override",
//		DisableFlagParsing:         true,
//		SuggestionsMinimumDistance: 2,
//		RunE:                       client.ValidateCmd,
//	}
//
//	cmd.AddCommand(
//		bankcli.GetBalancesCmd(),
//		GetTotalSupplyOverride(),
//		bankcli.GetCmdDenomsMetadata(),
//	)
//
//	return cmd
//}
//
//// GetTotalSupplyOverride overrides SDK bank's default total supply query to subtract locked enterprise nund
//func GetTotalSupplyOverride() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "total",
//		Short: "Query the total supply of coins of the chain",
//		Long: strings.TrimSpace(
//			fmt.Sprintf(`Query total supply of coins that are held by accounts in the chain.
//
//NOTE: This overrides the default Cosmos SDK total supply query, to subtract locked enterprise FUND from the result
//
//Example:
//  $ %s query %s total
//To query for the total supply of a specific coin denomination use:
//  $ %s query %s total --denom=[denom]
//`,
//				version.AppName, types.ModuleName, version.AppName, types.ModuleName,
//			),
//		),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			clientCtx, err := client.GetClientQueryContext(cmd)
//			if err != nil {
//				return err
//			}
//			denom, err := cmd.Flags().GetString(bankcli.FlagDenom)
//			if err != nil {
//				return err
//			}
//
//			queryClient := types.NewQueryClient(clientCtx)
//
//			if denom == "" {
//				res, err := queryClient.TotalSupplyOverride(cmd.Context(), &types.QueryTotalSupplyOverrideRequest{})
//				if err != nil {
//					return err
//				}
//
//				return clientCtx.PrintProto(res)
//			}
//
//			res, err := queryClient.SupplyOfOverride(cmd.Context(), &types.QuerySupplyOfOverrideRequest{Denom: denom})
//			if err != nil {
//				return err
//			}
//
//			return clientCtx.PrintProto(&res.Amount)
//		},
//	}
//	cmd.Flags().String(bankcli.FlagDenom, "", "The specific balance denomination to query for")
//	flags.AddQueryFlagsToCmd(cmd)
//	return cmd
//}
