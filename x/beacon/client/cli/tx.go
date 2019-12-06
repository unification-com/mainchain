package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
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

	beaconTxCmd.AddCommand(client.PostCommands(
		GetCmdRegisterBeacon(cdc),
		GetCmdRecordBeaconTimestamp(cdc),
	)...)

	return beaconTxCmd
}

// GetCmdRegisterBeacon is the CLI command for sending a RegisterBeacon transaction
func GetCmdRegisterBeacon(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "register [beacon moniker] [name]",
		Short: "register a new BEACON",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new BEACON, to enable timestamp hash submissions
Example:
$ %s tx %s register MyBeacon "My WRKChain" --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// first check if a BEACON exists with the same moniker.
			// The moniker should be a unique string identifier for the BEACON
			params := types.NewQueryBeaconParams(1, 1, args[0], sdk.AccAddress{})
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
				errMsg := fmt.Sprintf("beacon already registered with moniker '%s' - beacon id: %d, owner: %s", args[0], matchingBeacons[0].BeaconID, matchingBeacons[0].Owner)
				return types.ErrBeaconAlreadyRegistered(types.DefaultCodespace, errMsg)
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			paramsRetriever := keeper.NewParamsRetriever(cliCtx)
			beaconParams, err := paramsRetriever.GetParams()
			if err != nil {
				return err
			}

			txBldr = txBldr.WithFees(strconv.Itoa(int(beaconParams.FeeRegister)) + beaconParams.Denom)

			msg := types.NewMsgRegisterBeacon(args[0], args[1], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRecordBeaconTimestamp is the CLI command for sending a RecordBeaconTimestamp transaction
func GetCmdRecordBeaconTimestamp(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "record [beacon id] [hash] [submit time]",
		Short: "record a WRKChain's block hashes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a BEACON's' timestamp hash'
Example:
$ %s tx %s record 1 d04b98f48e8 1234356 --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			txBldr = txBldr.WithFees(strconv.Itoa(types.RecordFee) + types.FeeDenom)

			beaconID, err := strconv.Atoi(args[0])

			if err != nil {
				beaconID = 0
			}

			submitTime, err := strconv.Atoi(args[2])
			if err != nil {
				submitTime = 0
			}

			if submitTime == 0 {
				submitTime = int(time.Now().Unix())
			}

			msg := types.NewMsgRecordBeaconTimestamp(uint64(beaconID), args[1], uint64(submitTime), cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
