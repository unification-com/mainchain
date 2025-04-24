package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	FlagMoniker       = "moniker"
	FlagOwner         = "owner"
	FlagTimestampHash = "hash"
	FlagName          = "name"
	FlagSubmitTime    = "subtime"
)

func GetTxCmd() *cobra.Command {
	beaconTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Beacon transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	beaconTxCmd.AddCommand(
		GetCmdRegisterBeacon(),
		GetCmdRecordBeaconTimestamp(),
		GetCmdPurchaseStorage(),
	)

	return beaconTxCmd
}

// GetCmdRegisterBeacon is the CLI command for sending a RegisterBeacon transaction
func GetCmdRegisterBeacon() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a new BEACON",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new BEACON, to enable timestamp hash submissions
Example:
$ %s tx %s register --moniker=MyBeacon --name="My WRKChain" --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// used for getting fees and checking beacon
			queryClient := types.NewQueryClient(clientCtx)

			from := clientCtx.GetFromAddress()

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			beaconName, _ := cmd.Flags().GetString(FlagName)

			if len(moniker) == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "please enter a moniker")
			}
			if len(beaconName) == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "please enter a name")
			}

			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			regFee := strconv.Itoa(int(params.Params.FeeRegister)) + params.Params.Denom

			msg := types.NewMsgRegisterBeacon(moniker, beaconName, from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, regFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	//flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagMoniker, "", "BEACON's moniker")
	cmd.Flags().String(FlagName, "", "(optional) BEACON's name")
	return cmd
}

// GetCmdRecordBeaconTimestamp is the CLI command for sending a RecordBeaconTimestamp transaction
func GetCmdRecordBeaconTimestamp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [beacon_id]",
		Short: "record a BEACON's timestamp hash",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a BEACON's' timestamp hash'
Example:
$ %s tx %s record 1 --hash=d04b98f48e8 --subtime=1234356 --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// used for getting fees and checking beacon
			queryClient := types.NewQueryClient(clientCtx)

			from := clientCtx.GetFromAddress()

			hash, _ := cmd.Flags().GetString(FlagTimestampHash)
			submitTime, _ := cmd.Flags().GetUint64(FlagSubmitTime)

			if len(hash) == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "BEACON timestamp must have a Hash submitted")
			}

			if submitTime == 0 {
				submitTime = uint64(time.Now().Unix())
			}

			beaconId, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if beaconId == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "beacon_id must be > 0")
			}

			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			recFee := strconv.Itoa(int(params.Params.FeeRecord)) + params.Params.Denom

			msg := types.NewMsgRecordBeaconTimestamp(uint64(beaconId), hash, submitTime, from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, recFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	//flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagTimestampHash, "", "BEACON's timestamp hash")
	cmd.Flags().Uint64(FlagSubmitTime, 0, "BEACON's timestamp submission time")
	return cmd

}

// GetCmdPurchaseStorage is the CLI command for sending a PurchaseStorageAction transaction
func GetCmdPurchaseStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase_storage [beacon_id] [num_slots]",
		Short: "purchase more in-state storage for a BEACON",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Purchase more in-state storage for a BEACON, allowing more
timestamps to be kept in-state

Example:
$ %s tx %s purchase_storage 1 100
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// used for getting fees and checking beacon
			queryClient := types.NewQueryClient(clientCtx)

			from := clientCtx.GetFromAddress()

			beaconId, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if beaconId == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "beacon_id must be > 0")
			}

			numToPurchase, err := strconv.Atoi(args[1])

			if err != nil {
				return err
			}

			if numToPurchase == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "num_slots must be > 0")
			}

			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			purchaseFee := strconv.Itoa(int(params.Params.FeePurchaseStorage)*numToPurchase) + params.Params.Denom

			msg := types.NewMsgPurchaseBeaconStateStorage(uint64(beaconId), uint64(numToPurchase), from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, purchaseFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	//flags.AddTxFlagsToCmd(cmd)
	return cmd
}
