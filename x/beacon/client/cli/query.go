package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/viper"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	beaconQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the beacon module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	beaconQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdBeacon(storeKey, cdc),
		GetCmdBeaconTimestamp(storeKey, cdc),
		GetCmdSearchBeacons(storeKey, cdc),
	)...)
	return beaconQueryCmd
}

// GetCmdQueryParams implements a command to return the current Beacon parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current Beacon parameters",
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

// GetCmdBeacon queries information about a BEACON
func GetCmdBeacon(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "beacon [beacon id]",
		Short: "Query a BEACON for given ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			beaconID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/beacon/%s", queryRoute, beaconID), nil)
			if err != nil {
				fmt.Printf("could not find beacon - %s \n", beaconID)
				return nil
			}

			var out types.Beacon
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdSearchBeacons runs a BEACON search query with parameters
func GetCmdSearchBeacons(queryRoute string, cdc *codec.Codec) *cobra.Command {
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
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			moniker := viper.GetString(FlagMoniker)
			bechOwnerAddr := viper.GetString(FlagOwner)
			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagNumLimit)

			var ownerAddr sdk.AccAddress

			params := types.NewQueryBeaconParams(page, limit, moniker, ownerAddr)
			if len(bechOwnerAddr) != 0 {
				ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
				if err != nil {
					return err
				}
				params.Owner = ownerAddr
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryBeacons), bz)
			if err != nil {
				return err
			}

			var out types.QueryResBeacons
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().Int(FlagPage, 1, "pagination page of beacons to to query for")
	cmd.Flags().Int(FlagNumLimit, 100, "pagination limit of beacons to query for")
	cmd.Flags().String(FlagMoniker, "", "(optional) filter beacons by name")
	cmd.Flags().String(FlagOwner, "", "(optional) filter beacons by owner address")
	return cmd
}

// GetCmdBeaconTimestamp queries information about a beacon's recorded timestamp
func GetCmdBeaconTimestamp(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "timestamp [beacon id] [timestamp id]",
		Short: "Query a BEACON for given ID and timestamp ID to retrieve recorded timestamp",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			beaconID := args[0]
			timestampID := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/timestamp/%s/%s", queryRoute, beaconID, timestampID), nil)
			if err != nil {
				fmt.Printf("could not find beacon %s timestamp %s \n", beaconID, timestampID)
				return nil
			}

			var out types.BeaconTimestamp
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
