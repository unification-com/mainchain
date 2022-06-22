package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

func GetQueryCmd() *cobra.Command {
	wrkchainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wrkchain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	wrkchainQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdWrkChain(),
		GetCmdSearchWrkChains(),
		GetCmdWrkChainBlock(),
	)
	return wrkchainQueryCmd
}

// GetCmdQueryParams implements a command to return the current WRKChain parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current WRKChain parameters",
		Args:  cobra.NoArgs,
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

// GetCmdWrkChain queries information about a wrkchain
func GetCmdWrkChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wrkchain [wrkchain id]",
		Short: "Query a WRKChain for given ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the beacon id is a uint
			wrkchainId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("beacon_id %s not a valid int, please input a valid beacon_id", args[0])
			}

			res, err := queryClient.WrkChain(context.Background(), &types.QueryWrkChainRequest{
				WrkchainId: wrkchainId,
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

// GetCmdSearchWrkChains runs a WRKChain search query with parameters
func GetCmdSearchWrkChains() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Query all WRKChains with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all paginated WRKChains that match optional filters:

Example:
$ %s query wrkchain search --moniker wrkchain1
$ %s query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ %s query wrkchain search --page=2 --limit=100
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			bechOwnerAddr, _ := cmd.Flags().GetString(FlagOwner)
			pageReq, _ := client.ReadPageRequest(cmd.Flags())

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryWrkChainsFilteredRequest{
				Pagination: pageReq,
			}

			if len(moniker) > 0 {
				params.Moniker = moniker
			}

			if len(bechOwnerAddr) > 0 {
				params.Owner = bechOwnerAddr
			}

			res, err := queryClient.WrkChainsFiltered(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagMoniker, "", "(optional) filter wrkchains by name")
	cmd.Flags().String(FlagOwner, "", "(optional) filter wrkchains by owner address")
	return cmd
}

// GetCmdWrkChainBlock queries information about a wrkchain's recorded block hashes
func GetCmdWrkChainBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [wrkchain_id] [height]",
		Short: "Query a WRKChain for given ID and block height to retrieve recorded hashes for that block",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the wrkchain_id is a uint
			wrkchainId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("wrkchain_id %s not a valid int, please input a valid wrkchain_id", args[0])
			}

			// validate that the height is a uint
			height, err := strconv.ParseUint(args[1], 10, 64)

			if err != nil {
				return fmt.Errorf("height %s not a valid int, please input a valid height", args[1])
			}

			res, err := queryClient.WrkChainBlock(context.Background(), &types.QueryWrkChainBlockRequest{
				WrkchainId: wrkchainId,
				Height:     height,
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
