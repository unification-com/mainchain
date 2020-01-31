package cli

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain/x/beacon/internal/types"
)

const (
	FlagNumLimit      = "limit"
	FlagPage          = "page"
	FlagMoniker       = "moniker"
	FlagOwner         = "owner"
	FlagTimestampHash = "hash"
	FlagBeaconID      = "id"
	FlagName          = "name"
	FlagSubmitTime    = "subtime"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	beaconTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Beacon transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	beaconTxCmd.AddCommand(flags.PostCommands(
		GetCmdRegisterBeacon(cdc),
		GetCmdRecordBeaconTimestamp(cdc),
	)...)

	return beaconTxCmd
}

// GetCmdRegisterBeacon is the CLI command for sending a RegisterBeacon transaction
func GetCmdRegisterBeacon(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a new BEACON",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new BEACON, to enable timestamp hash submissions
Example:
$ %s tx %s register --moniker=MyBeacon --name="My WRKChain" --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(FlagMoniker)
			beaconName := viper.GetString(FlagName)

			// first check if a BEACON exists with the same moniker.
			// The moniker should be a unique string identifier for the BEACON
			params := types.NewQueryBeaconParams(1, 1, moniker, sdk.AccAddress{})
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryBeacons), bz)
			if err != nil {
				return err
			}
			var matchingBeacons types.QueryResBeacons
			err = cdc.UnmarshalJSON(res, &matchingBeacons)

			if err != nil {
				return err
			}

			// BEACON already registered with same moniker - output an error instead of broadcasting
			// the Tx and therefore charging reg fees
			if (len(matchingBeacons)) > 0 {
				errMsg := fmt.Sprintf("beacon already registered with moniker '%s' - beacon id: %d, owner: %s", moniker, matchingBeacons[0].BeaconID, matchingBeacons[0].Owner)
				return sdkerrors.Wrap(types.ErrBeaconAlreadyRegistered, errMsg)
			}

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			paramsRetriever := keeper.NewParamsRetriever(cliCtx)
			beaconParams, err := paramsRetriever.GetParams()
			if err != nil {
				return err
			}

			txBldr = txBldr.WithFees(strconv.Itoa(int(beaconParams.FeeRegister)) + beaconParams.Denom)

			msg := types.NewMsgRegisterBeacon(moniker, beaconName, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(FlagMoniker, "", "BEACON's moniker")
	cmd.Flags().String(FlagName, "", "(optional) BEACON's name")
	return cmd
}

// GetCmdRecordBeaconTimestamp is the CLI command for sending a RecordBeaconTimestamp transaction
func GetCmdRecordBeaconTimestamp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [beacon id]",
		Short: "record a BEACON's timestamp hash",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a BEACON's' timestamp hash'
Example:
$ %s tx %s record 1 --hash=d04b98f48e8 --subtime=1234356 --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			hash := viper.GetString(FlagTimestampHash)
			submitTime := viper.GetUint64(FlagSubmitTime)

			if len(hash) == 0 {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "BEACON timestamp must have a Hash submitted")
			}

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			txBldr = txBldr.WithFees(strconv.Itoa(types.RecordFee) + types.FeeDenom)

			beaconID, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if beaconID == 0 {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "BEACON id must be > 0")
			}

			if submitTime == 0 {
				submitTime = uint64(time.Now().Unix())
			}

			msg := types.NewMsgRecordBeaconTimestamp(uint64(beaconID), hash, submitTime, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTimestampHash, "", "BEACON's timestamp hash")
	cmd.Flags().Uint64(FlagSubmitTime, 0, "BEACON's timestamp submission time")
	return cmd

}
