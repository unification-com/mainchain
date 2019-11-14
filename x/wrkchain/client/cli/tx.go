package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	wrkchainTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "WRKChain transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	wrkchainTxCmd.AddCommand(client.PostCommands(
		GetCmdRegisterWrkChain(cdc),
		GetCmdRecordWrkChainBlock(cdc),
	)...)

	return wrkchainTxCmd
}

// GetCmdRegisterWrkChain is the CLI command for sending a RegisterWrkChain transaction
func GetCmdRegisterWrkChain(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "register [wrkchain moniker] [genesis hash] [name]",
		Short: "register a new WRKChain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new WRKChain, to enable WRKChain hash submissions
Example:
$ %s tx %s register MyWrkChain d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa "My WRKChain" --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// first check if a WRKChain exists with the same moniker.
			// The moniker should be a unique string identifier for the WRKChain
			params := types.NewQueryWrkChainParams(1, 1, args[0], sdk.AccAddress{})
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryWrkChains), bz)
			if err != nil {
				return err
			}
			var matchingWrkChains types.QueryResWrkChains
			err = cdc.UnmarshalJSON(res, &matchingWrkChains)

			if err != nil {
				return err
			}

			// WRKchain already registered with same moniker - output an error instead of broadcasting
			// the Tx and therefore charging reg fees
			if (len(matchingWrkChains)) > 0 {
				errMsg := fmt.Sprintf("wrkchain already registered with moniker '%s' - wrkchain id: %d, owner: %s", args[0], matchingWrkChains[0].WrkChainID, matchingWrkChains[0].Owner)
				return types.ErrWrkChainAlreadyRegistered(types.DefaultCodespace, errMsg)
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			txBldr = txBldr.WithFees(strconv.Itoa(types.RegFee) + types.FeeDenom)

			msg := types.NewMsgRegisterWrkChain(args[0], args[1], args[2], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdRecordWrkChainBlock is the CLI command for sending a RecordWrkChainBlock transaction
func GetCmdRecordWrkChainBlock(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "record [wrkchain id] [height] [block hash] [parent hash] [hash 1] [hash 2] [hash 3] --fees 1und",
		Short: "record a WRKChain's block hashes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a new WRKChain block's hashes'
Example:
$ %s tx %s record 1 24 d04b98f48e8 f8bcc15c6ae 5ac050801cd6 dcfd428fb5f9e 65c4e16e7807340fa --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		// todo - make parent hash, and hashes 1 - 3 optional
		Args: cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			txBldr = txBldr.WithFees(strconv.Itoa(types.RecordFee) + types.FeeDenom)

			height, err := strconv.Atoi(args[1])

			if err != nil {
				height = 0
			}

			wrkchainID, err := strconv.Atoi(args[0])

			if err != nil {
				wrkchainID = 0
			}

			msg := types.NewMsgRecordWrkChainBlock(uint64(wrkchainID), uint64(height), args[2], args[3], args[4], args[5], args[6], cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
