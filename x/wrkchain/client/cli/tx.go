package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unification-com/mainchain/x/wrkchain/internal/keeper"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

const (
	FlagNumLimit    = "limit"
	FlagPage        = "page"
	FlagMoniker     = "moniker"
	FlagOwner       = "owner"
	FlagMinHeight   = "min"
	FlagMaxHeight   = "max"
	FlagMinDate     = "after"
	FlagMaxDate     = "before"
	FlagBlockHash   = "block_hash"
	FlagParentHash  = "parent_hash"
	FlagHash1       = "hash1"
	FlagHash2       = "hash2"
	FlagHash3       = "hash3"
	FlagHeight      = "wc_height"
	FlagWrkChainID  = "id"
	FlagName        = "name"
	FlagBaseChain   = "base"
	FlagGenesisHash = "genesis"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	wrkchainTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "WRKChain transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	wrkchainTxCmd.AddCommand(flags.PostCommands(
		GetCmdRegisterWrkChain(cdc),
		GetCmdRecordWrkChainBlock(cdc),
	)...)

	return wrkchainTxCmd
}

// GetCmdRegisterWrkChain is the CLI command for sending a RegisterWrkChain transaction
func GetCmdRegisterWrkChain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a new WRKChain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new WRKChain, to enable WRKChain hash submissions
Example:
$ %s tx %s register --moniker="MyWrkChain" --genesis="d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa" --name="My WRKChain" --base="geth" --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			moniker := viper.GetString(FlagMoniker)
			wrkchainName := viper.GetString(FlagName)
			wrkchainBase := viper.GetString(FlagBaseChain)
			wrkchainGenesisHash := viper.GetString(FlagGenesisHash)

			if len(moniker) == 0 {
				return sdk.ErrInternal("WRKChain must have a moniker")
			}

			// first check if a WRKChain exists with the same moniker.
			// The moniker should be a unique string identifier for the WRKChain
			params := types.NewQueryWrkChainParams(1, 1, moniker, sdk.AccAddress{})
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryWrkChainsFiltered), bz)
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
				errMsg := fmt.Sprintf("wrkchain already registered with moniker '%s' - wrkchain id: %d, owner: %s", moniker, matchingWrkChains[0].WrkChainID, matchingWrkChains[0].Owner)
				return types.ErrWrkChainAlreadyRegistered(types.DefaultCodespace, errMsg)
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			paramsRetriever := keeper.NewParamsRetriever(cliCtx)
			wrkchainParams, err := paramsRetriever.GetParams()
			if err != nil {
				return err
			}

			txBldr = txBldr.WithFees(strconv.Itoa(int(wrkchainParams.FeeRegister)) + wrkchainParams.Denom)

			msg := types.NewMsgRegisterWrkChain(moniker, wrkchainGenesisHash, wrkchainName, wrkchainBase, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(FlagMoniker, "", "WRKChain's moniker")
	cmd.Flags().String(FlagName, "", "(optional) WRKChain's name")
	cmd.Flags().String(FlagGenesisHash, "", "(optional) WRKChain's Genesis hash")
	cmd.Flags().String(FlagBaseChain, "", "(optional) WRKChain's chain type - geth, etc.")
	return cmd
}

// GetCmdRecordWrkChainBlock is the CLI command for sending a RecordWrkChainBlock transaction
func GetCmdRecordWrkChainBlock(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [wrkchain id]",
		Short: "record a WRKChain's block hashes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a new WRKChain block's hash(es)'
Example:
$ %s tx %s record 1 --wc_height=24 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --hash1="5ac050801cd6" --hash2="dcfd428fb5f9e" --hash3="65c4e16e7807340fa" --from mykey
$ %s tx %s record 1 --wc_height=25 --block_hash="d04b98f48e8" --from mykey
$ %s tx %s record 1 --wc_height=26 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --from mykey
`,
				version.ClientName, types.ModuleName,
				version.ClientName, types.ModuleName,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			height := viper.GetUint64(FlagHeight)
			blockHash := viper.GetString(FlagBlockHash)
			parentHash := viper.GetString(FlagParentHash)
			hash1 := viper.GetString(FlagHash1)
			hash2 := viper.GetString(FlagHash2)
			hash3 := viper.GetString(FlagHash3)

			if len(blockHash) == 0 {
				return sdk.ErrInternal("WRKChain block must have a Hash submitted")
			}

			if height == 0 {
				return sdk.ErrInternal("WRKChain block hash submission must be for height > 0")
			}

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// automatically apply fees
			txBldr = txBldr.WithFees(strconv.Itoa(types.RecordFee) + types.FeeDenom)

			wrkchainID, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if wrkchainID == 0 {
				return sdk.ErrInternal("WRKChain id must be > 0")
			}

			msg := types.NewMsgRecordWrkChainBlock(uint64(wrkchainID), height, blockHash, parentHash, hash1, hash2, hash3, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Uint64(FlagHeight, 0, "WRKChain block's height/block number")
	cmd.Flags().String(FlagBlockHash, "", "WRKChain block's header (main) hash")
	cmd.Flags().String(FlagParentHash, "", "(optional) WRKChain block's parent hash")
	cmd.Flags().String(FlagHash1, "", "(optional) Additional WRKChain hash - e.g. State Merkle Root")
	cmd.Flags().String(FlagHash2, "", "(optional) Additional WRKChain hash - e.g. Tx Merkle Root")
	cmd.Flags().String(FlagHash3, "", "(optional) Additional WRKChain hash")

	return cmd
}
