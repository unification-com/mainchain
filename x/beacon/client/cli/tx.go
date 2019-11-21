package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
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
	)...)

	return beaconTxCmd
}
