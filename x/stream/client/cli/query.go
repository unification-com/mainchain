package cli

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/unification-com/mainchain/x/stream/types"
)

const (
	FlagStreamCoin     = "coin"
	FlagStreamDuration = "duration"
	FlagStreamPeriod   = "period"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group stream queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryParams(),
		GetCmdCalculateFlowRate(),
		GetCmdGetAllStreams(),
	)

	return cmd
}

// GetCmdCalculateFlowRate calculates a flow rate
func GetCmdCalculateFlowRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate_flow",
		Short: "Calculate the Flow Rate for vigen parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Calculate a flow rate given parameters:

Example:
$ %s query stream calculate_flow --coin 1000000000nund --period month --duration 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			strPeriod, _ := cmd.Flags().GetString(FlagStreamPeriod)
			strCoin, _ := cmd.Flags().GetString(FlagStreamCoin)
			duration, _ := cmd.Flags().GetUint64(FlagStreamDuration)
			period := types.PeriodEnumFromString(strPeriod)

			if period == types.StreamPeriodUnspecified {
				return fmt.Errorf("period %s not valid. Use second/hour/day/week/month/year", strPeriod)
			}

			coin, err := sdk.ParseCoinNormalized(strCoin)

			if err != nil {
				return err
			}

			if duration < 1 {
				return fmt.Errorf("duration cannot be zero")
			}

			clientCtx, err := client.GetClientQueryContext(cmd)

			totalDuration, _, flowRateInt64 := types.CalculateFlowRateForCoin(coin, period, duration)

			res := &types.QueryCalculateFlowRateResponse{
				Coin:     coin,
				Period:   period,
				Duration: duration,
				Seconds:  totalDuration,
				FlowRate: flowRateInt64,
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(FlagStreamCoin, "", "coin e.g. 1000000000nund")
	cmd.Flags().Uint64(FlagStreamDuration, 0, "number of periods e.g. 1")
	cmd.Flags().String(FlagStreamPeriod, "", "second/hour/day/week/month/year")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetAllStreams queries a list of all streams
func GetCmdGetAllStreams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streams",
		Short: "Query all streams",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated streams

Example:
$ %s query stream streams
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())

			params := &types.QueryStreamsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Streams(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pagination")
	return cmd
}
