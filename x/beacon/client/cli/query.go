package cli

import (
	"fmt"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
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

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", queryRoute, beaconID), nil)
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

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get-ts/%s/%s", queryRoute, beaconID, timestampID), nil)
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

// Todo - query by params - sub time, hash & get all with pagination
