package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	entutils "github.com/unification-com/mainchain/x/enterprise/client/utils"
	"github.com/unification-com/mainchain/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	enterpriseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the enterprise module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	enterpriseQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdGetPurchaseOrders(storeKey, cdc),
		GetCmdGetPurchaseOrderByID(storeKey, cdc),
		GetCmdGetLockedUndByAddress(storeKey, cdc),
		GetCmdQueryTotalLocked(storeKey, cdc),
		GetCmdQueryTotalUnlocked(storeKey, cdc),
	)...)
	return enterpriseQueryCmd
}

// GetCmdQueryParams implements a command to return the current enterprise und
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current enterprise UND parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdGetPurchaseOrders queries a list of all purchase orders
func GetCmdGetPurchaseOrders(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orders",
		Short: "Query Enterprise UND purchase orders with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated Enterprise UND purchase orders that match optional filters:

Example:
$ %s query enterprise orders --status (raised|accept|reject|complete)
$ %s query enterprise orders --purchaser und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ %s query enterprise orders --page=2 --limit=100
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			strProposalStatus := viper.GetString(FlagPurchaseOrderStatus)
			bechPurchaserAddr := viper.GetString(FlagPurchaser)
			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagNumLimit)

			var purchaseOrderStatus types.PurchaseOrderStatus
			var purchaserAddr sdk.AccAddress

			params := types.NewQueryPurchaseOrdersParams(page, limit, purchaseOrderStatus, purchaserAddr)

			if len(strProposalStatus) != 0 {
				purchaseOrderStatus, err := types.PurchaseOrderStatusFromString(entutils.NormalisePurchaseOrderStatus(strProposalStatus))
				if err != nil {
					return err
				}
				params.PurchaseOrderStatus = purchaseOrderStatus
			}

			if len(bechPurchaserAddr) != 0 {
				purchaserAddr, err := sdk.AccAddressFromBech32(bechPurchaserAddr)
				if err != nil {
					return err
				}
				params.Purchaser = purchaserAddr
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryPurchaseOrders), bz)
			if err != nil {
				return err
			}

			var matchingOrders types.QueryResPurchaseOrders
			err = cdc.UnmarshalJSON(res, &matchingOrders)

			if err != nil {
				return err
			}

			if len(matchingOrders) == 0 {
				return fmt.Errorf("no matching purchase orders found")
			}

			return cliCtx.PrintOutput(matchingOrders)
		},
	}

	cmd.Flags().Int(FlagPage, 1, "pagination page of purchase orders to to query for")
	cmd.Flags().Int(FlagNumLimit, 100, "pagination limit of purchase orders to query for")
	cmd.Flags().String(FlagPurchaseOrderStatus, "", "(optional) filter purchase orders by status, status: raised/accept/reject/complete")
	cmd.Flags().String(FlagPurchaser, "", "(optional) filter purchase orders raised by address")
	return cmd
}

// GetCmdGetPurchaseOrderByID queries a purchase order given an ID
func GetCmdGetPurchaseOrderByID(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "order [purchase_order_id]",
		Short: "get a purchase order by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryGetPurchaseOrder, args[0]), nil)
			if err != nil {
				fmt.Printf("could not get query purchase order: ID %s\n", args[0])
				return err
			}

			var out types.EnterpriseUndPurchaseOrder
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdGetLockedUndByAddress queries locked UND for a given address
func GetCmdGetLockedUndByAddress(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "locked [address]",
		Short: "get locked UND for an address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryGetLocked, args[0]), nil)
			if err != nil {
				fmt.Printf("could not get query locked UND for address %s\n", args[0])
				return err
			}

			var out types.LockedUnd
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryTotalLocked implements a command to return the current total locked enterprise und
func GetCmdQueryTotalLocked(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "total-locked",
		Short: "Query the current total locked enterprise UND",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryTotalLocked), nil)
			if err != nil {
				fmt.Printf("could not get query total locked enterprise UND\n")
				return err
			}

			var out sdk.Coin
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryTotalUnlocked implements a command to return the current total locked enterprise und
func GetCmdQueryTotalUnlocked(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "total-unlocked",
		Short: "Query the current total unlocked und in circulation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryTotalUnlocked), nil)
			if err != nil {
				fmt.Printf("could not get query total unlocked UND\n")
				return err
			}

			var out sdk.Coin
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
