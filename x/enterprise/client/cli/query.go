package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	enterpriseQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the enterprise module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	enterpriseQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdGetPurchaseOrders(storeKey, cdc),
		GetCmdGetPurchaseOrderByID(storeKey, cdc),
		GetCmdGetLockedUndByAddress(storeKey, cdc),
		GetCmdQueryTotalLocked(storeKey, cdc),
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
	return &cobra.Command{
		Use:   "get-all-pos",
		Short: "get all current raised Enterprise UND purchase orders",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryPurchaseOrders), nil)
			if err != nil {
				fmt.Printf("could not get query raised purchase orders\n")
				return nil
			}

			var out types.QueryResRaisedPurchaseOrders
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdGetPurchaseOrderByID queries a purchase order given an ID
func GetCmdGetPurchaseOrderByID(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [purchase_order_id]",
		Short: "get a purchase order by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryGetPurchaseOrder, args[0]), nil)
			if err != nil {
				fmt.Printf("could not get query raised purchase order: ID %s\n", args[0])
				return nil
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
				return nil
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
				return nil
			}

			var out sdk.Coin
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
