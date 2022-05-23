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

	"github.com/unification-com/mainchain/x/beacon/types"
)

func GetQueryCmd() *cobra.Command {
	beaconQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the beacon module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	beaconQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdBeacon(),
		GetCmdBeaconTimestamp(),
		GetCmdSearchBeacons(),
	)
	return beaconQueryCmd
}

// GetCmdQueryParams implements a command to return the current Beacon parameters.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current Beacon parameters",
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

			return clientCtx.PrintObjectLegacy(params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdBeacon queries information about a BEACON
func GetCmdBeacon() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beacon [beacon_id]",
		Short: "Query a BEACON for given ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the beacon id is a uint
			beaconId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("beacon_id %s not a valid int, please input a valid beacon_id", args[0])
			}

			res, err := queryClient.Beacon(context.Background(), &types.QueryBeaconRequest{
				BeaconId: beaconId,
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

// GetCmdSearchBeacons runs a BEACON search query with parameters
func GetCmdSearchBeacons() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Query all BEACONs with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all paginated BEACONs that match optional filters:

Example:
$ %s query beacon search --moniker beacon1
$ %s query beacon search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ %s query beacon search --page=2 --limit=100
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

			params := &types.QueryBeaconsFilteredRequest{
				Pagination: pageReq,
			}

			if len(moniker) > 0 {
				params.Moniker = moniker
			}

			if len(bechOwnerAddr) > 0 {
				params.Owner = bechOwnerAddr
			}

			res, err := queryClient.BeaconsFiltered(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagMoniker, "", "(optional) filter beacons by name")
	cmd.Flags().String(FlagOwner, "", "(optional) filter beacons by owner address")
	return cmd
}

// GetCmdBeaconTimestamp queries information about a beacon's recorded timestamp
func GetCmdBeaconTimestamp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timestamp [beacon_id] [timestamp_id]",
		Short: "Query a BEACON for given ID and timestamp ID to retrieve recorded timestamp",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the beacon id is a uint
			beaconId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("beacon_id %s not a valid int, please input a valid beacon_id", args[0])
			}

			// validate that the timestamp id is a uint
			timestampId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("timestamp_id %s not a valid int, please input a valid timestamp_id", args[1])
			}

			res, err := queryClient.BeaconTimestamp(context.Background(), &types.QueryBeaconTimestampRequest{
				BeaconId:    beaconId,
				TimestampId: timestampId,
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
