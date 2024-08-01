package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

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
		GetCmdGetAllStreamsByReceiver(),
		GetCmdGetStream(),
		GetCmdGetStreamByReceiverSenderCurrentFlow(),
		GetCmdGetAllStreamsBySender(),
	)

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdCalculateFlowRate calculates a flow rate
func GetCmdCalculateFlowRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calculate_flow",
		Short: "Calculate the Flow Rate for given parameters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Calculate a flow rate given parameters:

Note - duration is the frequency of the payment, e.g. every 2 months

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

// GetCmdGetAllStreamsByReceiver queries a list of all streams for given receiver
func GetCmdGetAllStreamsByReceiver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receiver_streams",
		Short: "Query all streams for given receiver",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated streams for given receiver

Example:
$ %s query stream receiver_streams [receiver_address]
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			receiverAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())

			params := &types.QueryAllStreamsForReceiverRequest{
				ReceiverAddr: receiverAddr.String(),
				Pagination:   pageReq,
			}

			res, err := queryClient.AllStreamsForReceiver(context.Background(), params)
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

func GetCmdGetStream() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream",
		Short: "Query a stream by Receiver/Sender pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a stream by Receiver/Sender pair

Example:
$ %s query stream stream [receiver_addr] [sender_addr]
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			receiverAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			senderAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryStreamByReceiverSenderRequest{
				ReceiverAddr: receiverAddr.String(),
				SenderAddr:   senderAddr.String(),
			}

			res, err := queryClient.StreamByReceiverSender(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdGetStreamByReceiverSenderCurrentFlow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stream_receiver_sender_flow",
		Short: "Query a stream by Receiver/Sender pair",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a stream's current flow data by Receiver/Sender pair

Example:
$ %s query stream stream_receiver_sender_flow [receiver_addr] [sender_addr]
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			receiverAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			senderAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := &types.QueryStreamReceiverSenderCurrentFlowRequest{
				ReceiverAddr: receiverAddr.String(),
				SenderAddr:   senderAddr.String(),
			}

			res, err := queryClient.StreamReceiverSenderCurrentFlow(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdGetAllStreamsBySender queries a list of all streams for given receiver
func GetCmdGetAllStreamsBySender() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sender_streams",
		Short: "Query all streams for given sender",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated streams for given sender

Example:
$ %s query stream sender_streams [sender_address]
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			senderAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())

			params := &types.QueryAllStreamsForSenderRequest{
				SenderAddr: senderAddr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.AllStreamsForSender(context.Background(), params)
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
