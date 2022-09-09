package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	entutils "github.com/unification-com/mainchain/x/enterprise/client/utils"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	enterpriseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the enterprise module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	enterpriseQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdGetPurchaseOrders(),
		GetCmdGetPurchaseOrderByID(),
		GetCmdGetLockedUndByAddress(),
		GetCmdQueryTotalLocked(),
		GetCmdQueryTotalUnlocked(),
		GetCmdGetWhitelistedAddresses(),
		GetCmdGetAddresIsWhitelisted(),
		GetCmdGetEnterpriseUserAccount(),
		GetCmdGetEnterpriseSupply(),
		GetCmdGetSpentEFUNDByAddress(),
		GetCmdQueryTotalSpentEFUND(),
	)

	return enterpriseQueryCmd
}

// GetCmdQueryParams implements a command to return the current enterprise FUND
// parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current enterprise FUND parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all the current enterprise FUND parameters.

Example:
$ %s query enterprise params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query store for all params
			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdGetPurchaseOrders queries a list of all purchase orders
func GetCmdGetPurchaseOrders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orders",
		Short: "Query Enterprise FUND purchase orders with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated Enterprise FUND purchase orders that match optional filters:

Example:
$ %s query enterprise orders --status (raised|accept|reject|complete)
$ %s query enterprise orders --purchaser und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ %s query enterprise orders --page=2 --limit=100
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			strProposalStatus, _ := cmd.Flags().GetString(FlagPurchaseOrderStatus)
			purchaserAddr, _ := cmd.Flags().GetString(FlagPurchaser)

			statusNorm := entutils.NormalisePurchaseOrderStatus(strProposalStatus)
			status, err := types.PurchaseOrderStatusFromString(statusNorm)

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())

			params := &types.QueryEnterpriseUndPurchaseOrdersRequest{
				Pagination: pageReq,
			}

			if status != types.StatusNil {
				params.Status = status
			}

			if len(purchaserAddr) > 0 {
				params.Purchaser = purchaserAddr
			}

			res, err := queryClient.EnterpriseUndPurchaseOrders(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagPurchaseOrderStatus, "", "(optional) filter purchase orders by status, status: raised/accept/reject/complete")
	cmd.Flags().String(FlagPurchaser, "", "(optional) filter purchase orders raised by address")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "purchase orders")
	return cmd
}

// GetCmdGetPurchaseOrderByID queries a purchase order given an ID
func GetCmdGetPurchaseOrderByID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order [purchase_order_id]",
		Short: "get a purchase order by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			purchaseOrderId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("purchase_order_id %s not a valid int, please input a valid purchase_order_id", args[0])
			}

			res, err := queryClient.EnterpriseUndPurchaseOrder(context.Background(), &types.QueryEnterpriseUndPurchaseOrderRequest{
				PurchaseOrderId: purchaseOrderId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetLockedUndByAddress queries locked FUND for a given address
func GetCmdGetLockedUndByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locked [address]",
		Short: "get locked FUND for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			purchaser, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.LockedUndByAddress(context.Background(), &types.QueryLockedUndByAddressRequest{
				Owner: purchaser.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTotalLocked implements a command to return the current total locked enterprise und
func GetCmdQueryTotalLocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-locked",
		Short: "Query the current total locked enterprise FUND",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalLocked(context.Background(), &types.QueryTotalLockedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTotalUnlocked implements a command to return the current total locked enterprise und
func GetCmdQueryTotalUnlocked() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-unlocked",
		Short: "Query the current total unlocked und in circulation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalUnlocked(context.Background(), &types.QueryTotalUnlockedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetWhitelistedAddresses queries all addresses whitelisted for raising enterprise und purchase orders
func GetCmdGetWhitelistedAddresses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist",
		Short: "get addresses whitelisted for raising enterprise purchase orders",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Whitelist(context.Background(), &types.QueryWhitelistRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetLockedUndByAddress queries locked FUND for a given address
func GetCmdGetAddresIsWhitelisted() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelisted [address]",
		Short: "check if given address is whitelested for purchase orders",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.Whitelisted(context.Background(), &types.QueryWhitelistedRequest{
				Address: address.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetLockedUndByAddress queries locked FUND for a given address
func GetCmdGetEnterpriseUserAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "get data about an address - locked, unlocked and total FUND",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.EnterpriseAccount(context.Background(), &types.QueryEnterpriseAccountRequest{
				Address: address.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetEnterpriseSupply queries eFUND data, including locked, unlocked and total chain supply
func GetCmdGetEnterpriseSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ent-supply",
		Short: "get eFUND data, including locked, unlocked and chain total supply",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.EnterpriseSupply(context.Background(), &types.QueryEnterpriseSupplyRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetSpentEFUNDByAddress queries spent eFUND for a given address
func GetCmdGetSpentEFUNDByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spent [address]",
		Short: "get spent eFUND for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			purchaser, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, err := queryClient.SpentEFUNDByAddress(context.Background(), &types.QuerySpentEFUNDByAddressRequest{
				Address: purchaser.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryTotalSpentEFUND implements a command to return the current total locked enterprise und
func GetCmdQueryTotalSpentEFUND() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "total-spent",
		Short: "Query the current total spent eFUND",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalSpentEFUND(context.Background(), &types.QueryTotalSpentEFUNDRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
