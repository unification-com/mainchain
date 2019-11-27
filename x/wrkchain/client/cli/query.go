package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	wrkchainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wrkchain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	wrkchainQueryCmd.AddCommand(client.GetCommands(
		GetCmdQueryParams(cdc),
		GetCmdWrkChain(storeKey, cdc),
		GetCmdWrkChains(storeKey, cdc),
		GetCmdWrkChainBlock(storeKey, cdc),
		GetCmdWrkChainBlockHashes(storeKey, cdc),
		GetCmdSearchWrkChainBlockHashes(storeKey, cdc),
	)...)
	return wrkchainQueryCmd
}

// GetCmdQueryParams implements a command to return the current WRKChain parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current WRKChain parameters",
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

// GetCmdWrkChain queries information about a wrkchain
func GetCmdWrkChain(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "wrkchain [wrkchain id]",
		Short: "Query a WRKChain for given ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			wrkchainId := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryWrkChain, wrkchainId), nil)
			if err != nil {
				fmt.Printf("could not find WRKChain - %s \n", wrkchainId)
				return nil
			}

			var out types.WrkChain
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdWrkChain queries information about a wrkchain
func GetCmdWrkChains(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Query all WRKChains with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all paginated WRKChains that match optional filters:

Example:
$ %s query wrkchain search --moniker wrkchain1
$ %s query wrkchain search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ %s query wrkchain search --page=2 --limit=100
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {

			moniker := viper.GetString(FlagMoniker)
			bechOwnerAddr := viper.GetString(FlagOwner)
			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagNumLimit)

			var ownerAddr sdk.AccAddress

			params := types.NewQueryWrkChainParams(page, limit, moniker, ownerAddr)
			if len(bechOwnerAddr) != 0 {
				ownerAddr, err := sdk.AccAddressFromBech32(bechOwnerAddr)
				if err != nil {
					return err
				}
				params.Owner = ownerAddr
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryWrkChainsFiltered), bz)
			if err != nil {
				return err
			}

			var out types.QueryResWrkChains
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().Int(FlagPage, 1, "pagination page of wrkchain hashes to to query for")
	cmd.Flags().Int(FlagNumLimit, 100, "pagination limit of wrkchain hashes to query for")
	cmd.Flags().String(FlagMoniker, "", "(optional) filter wrkchains by name")
	cmd.Flags().String(FlagOwner, "", "(optional) filter wrkchains by owner address")
	return cmd
}

// GetCmdWrkChainBlock queries information about a wrkchain's recorded block hashes
func GetCmdWrkChainBlock(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "block [wrkchain id] [height]",
		Short: "Query a WRKChain for given ID and block height to retrieve recorded hashes for that block",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			wrkchainId := args[0]
			height := args[1]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", queryRoute, keeper.QueryWrkChainBlock, wrkchainId, height), nil)
			if err != nil {
				fmt.Printf("could not find WRKChain %s block hashes at height %s \n", wrkchainId, height)
				return nil
			}

			var out types.WrkChainBlock
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdWrkChainBlockHashes queries a list of all recorded block hashes for a WRKChain
func GetCmdWrkChainBlockHashes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "blocks [wrkchain id]",
		Short: "Query a WRKChain for all hashes recorded to date",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			wrkchainId := args[0]
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryWrkChainBlockHashes, wrkchainId), nil)
			if err != nil {
				fmt.Printf("could not get query block hashes\n")
				return nil
			}

			var out types.QueryResWrkChainBlockHashes
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdWrkChainBlockHashes queries a list of all recorded block hashes for a WRKChain
func GetCmdSearchWrkChainBlockHashes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search-blocks [wrkchain id]",
		Short: "Query a WRKChain for all hashes recorded to date",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for all paginated hashes for a WRKChain that match optional filters:

Example:
$ %s query wrkchain search-blocks 1 --before 1574871069 --after 1573481124
$ %s query wrkchain search-blocks 1 --min 123 --max 456
$ %s query wrkchain search-blocks 1 --page=2 --limit=100
`,
				version.ClientName, version.ClientName, version.ClientName,
			),
		),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			minHeight := viper.GetUint64(FlagMinHeight)
			maxHeight := viper.GetUint64(FlagMaxHeight)
			minDate := viper.GetUint64(FlagMinDate)
			maxDate := viper.GetUint64(FlagMaxDate)
			page := viper.GetInt(FlagPage)
			limit := viper.GetInt(FlagNumLimit)

			params := types.NewQueryWrkChainBlockParams(page, limit, minHeight, maxHeight, minDate, maxDate, "")

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			wrkchainId := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, keeper.QueryWrkChainHashesFiltered, wrkchainId), bz)
			if err != nil {
				fmt.Printf("could not get query block hashes\n")
				return nil
			}

			var out types.QueryResWrkChainBlockHashes
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
	cmd.Flags().Int(FlagPage, 1, "pagination page of wrkchain hashes to to query for")
	cmd.Flags().Int(FlagNumLimit, 100, "pagination limit of wrkchain hashes to query for")
	cmd.Flags().Uint64(FlagMinHeight, 0, "(optional) filter blocks by minimum height")
	cmd.Flags().Uint64(FlagMaxHeight, 0, "(optional) filter blocks by maximum height")
	cmd.Flags().Uint64(FlagMinDate, 0, "(optional) filter blocks submitted after this date")
	cmd.Flags().Uint64(FlagMaxDate, 0, "(optional) filter blocks submitted before this date")
	return cmd
}
