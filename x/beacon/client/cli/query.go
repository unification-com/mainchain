package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/viper"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
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
		GetCmdBeaconTimestamps(storeKey, cdc),
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

// GetCmdBeaconTimestamps queries information about a beacon's recorded timestamps
func GetCmdBeaconTimestamps(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timestamps [beacon id]",
		Short: "Query a BEACON for given ID to retrieve recorded timestamps",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all paginated hashes for a WRKChain that match optional filters:

Example:
$ %s query beacon timestamps 1 --before 1574871069 --after 1573481124
$ %s query beacon timestamps 1 --min 123 --max 456
$ %s query beacon timestamps 1 --page=2 --limit=100
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagNumLimit)
			subTime := viper.GetUint64(FlagSubmitTime)
			hash := viper.GetString(FlagTimestampHash)

			beaconID, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryBeaconTimestampParams(page, limit, uint64(beaconID), hash, subTime)

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/timestamps", queryRoute), bz)
			if err != nil {
				fmt.Printf("could not find beacon %d timestamps \n", beaconID)
				return nil
			}

			var out types.QueryResBeaconTimestampHashes
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().Int(FlagPage, 1, "pagination page of beacon timestamps to to query for")
	cmd.Flags().Int(FlagNumLimit, 100, "pagination limit of beacon timestamps to query for")
	cmd.Flags().Uint64(FlagSubmitTime, 0, "(optional) search by submit time")
	cmd.Flags().String(FlagTimestampHash, "", "(optional) search for a particular hash")
	return cmd
}

// Todo - query by params - sub time, hash & get all with pagination
