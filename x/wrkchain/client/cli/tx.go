package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	FlagMoniker     = "moniker"
	FlagOwner       = "owner"
	FlagBlockHash   = "block_hash"
	FlagParentHash  = "parent_hash"
	FlagHash1       = "hash1"
	FlagHash2       = "hash2"
	FlagHash3       = "hash3"
	FlagWcHeight    = "wc_height"
	FlagName        = "name"
	FlagBaseChain   = "base"
	FlagGenesisHash = "genesis"
)

func GetTxCmd() *cobra.Command {
	wrkchainTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "WRKChain transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	wrkchainTxCmd.AddCommand(
		GetCmdRegisterWrkChain(),
		GetCmdRecordWrkChainBlock(),
		GetCmdPurchaseStorage(),
	)

	return wrkchainTxCmd
}

// GetCmdRegisterWrkChain is the CLI command for sending a RegisterWrkChain transaction
func GetCmdRegisterWrkChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a new WRKChain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register a new WRKChain, to enable WRKChain hash submissions
Example:
$ %s tx %s register --moniker="MyWrkChain" --genesis="d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa" --name="My WRKChain" --base="geth" --from mykey
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

			// used for getting fees and checking wrkchain
			queryClient := types.NewQueryClient(clientCtx)

			from := clientCtx.GetFromAddress()

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			wrkchainName, _ := cmd.Flags().GetString(FlagName)
			wrkchainBase, _ := cmd.Flags().GetString(FlagBaseChain)
			wrkchainGenesisHash, _ := cmd.Flags().GetString(FlagGenesisHash)

			if len(moniker) == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "WRKChain must have a moniker")
			}

			if len(wrkchainName) == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "WRKChain must have a name")
			}

			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			if err != nil {
				return err
			}

			regFee := strconv.Itoa(int(params.Params.FeeRegister)) + params.Params.Denom

			msg := types.NewMsgRegisterWrkChain(moniker, wrkchainGenesisHash, wrkchainName, wrkchainBase, from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, regFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String(FlagMoniker, "", "WRKChain's moniker")
	cmd.Flags().String(FlagName, "", "(optional) WRKChain's name")
	cmd.Flags().String(FlagGenesisHash, "", "(optional) WRKChain's Genesis hash")
	cmd.Flags().String(FlagBaseChain, "", "(optional) WRKChain's chain type - geth, etc.")
	return cmd
}

// GetCmdRecordWrkChainBlock is the CLI command for sending a RecordWrkChainBlock transaction
func GetCmdRecordWrkChainBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record [wrkchain_id]",
		Short: "record a WRKChain's block hashes",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Record a new WRKChain block's hash(es)'
Example:
$ %s tx %s record 1 --wc_height=24 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --hash1="5ac050801cd6" --hash2="dcfd428fb5f9e" --hash3="65c4e16e7807340fa" --from mykey
$ %s tx %s record 1 --wc_height=25 --block_hash="d04b98f48e8" --from mykey
$ %s tx %s record 1 --wc_height=26 --block_hash="d04b98f48e8" --parent_hash="f8bcc15c6ae" --from mykey
`,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			// used for getting fees and checking wrkchain
			queryClient := types.NewQueryClient(clientCtx)

			height, _ := cmd.Flags().GetUint64(FlagWcHeight)
			blockHash, _ := cmd.Flags().GetString(FlagBlockHash)
			parentHash, _ := cmd.Flags().GetString(FlagParentHash)
			hash1, _ := cmd.Flags().GetString(FlagHash1)
			hash2, _ := cmd.Flags().GetString(FlagHash2)
			hash3, _ := cmd.Flags().GetString(FlagHash3)

			wrkchainId, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if wrkchainId == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wrkchain_id must be > 0")
			}

			if len(blockHash) == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "WRKChain block must have a Hash submitted")
			}

			if height == 0 {
				return errorsmod.Wrap(types.ErrMissingData, "WRKChain block hash submission must be for height > 0")
			}

			params, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{},
			)

			recFee := strconv.Itoa(int(params.Params.FeeRecord)) + params.Params.Denom

			msg := types.NewMsgRecordWrkChainBlock(uint64(wrkchainId), height, blockHash, parentHash, hash1, hash2, hash3, from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, recFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagWcHeight, 0, "WRKChain block's height/block number")
	cmd.Flags().String(FlagBlockHash, "", "WRKChain block's header (main) hash")
	cmd.Flags().String(FlagParentHash, "", "(optional) WRKChain block's parent hash")
	cmd.Flags().String(FlagHash1, "", "(optional) Additional WRKChain hash - e.g. State Merkle Root")
	cmd.Flags().String(FlagHash2, "", "(optional) Additional WRKChain hash - e.g. Tx Merkle Root")
	cmd.Flags().String(FlagHash3, "", "(optional) Additional WRKChain hash")

	return cmd
}

// GetCmdPurchaseStorage is the CLI command for sending a PurchaseStorageAction transaction
func GetCmdPurchaseStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purchase_storage [wrkchain_id] [num_slots]",
		Short: "purchase more in-state storage for a WrkChain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Purchase more in-state storage for a WrkChain, allowing more
hashes to be kept in-state

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

			// used for getting fees and checking wrkchain
			queryClient := types.NewQueryClient(clientCtx)

			from := clientCtx.GetFromAddress()

			wrkchainId, err := strconv.Atoi(args[0])

			if err != nil {
				return err
			}

			if wrkchainId == 0 {
				return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "wrkchain_id must be > 0")
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

			msg := types.NewMsgPurchaseWrkChainStateStorage(uint64(wrkchainId), uint64(numToPurchase), from)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := cmd.Flags().Set(flags.FlagFees, purchaseFee); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
