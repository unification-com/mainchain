package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

const flagGenesisTime = "genesis-time"
const chainUpgradeGuide = "https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html"

// MigrateGenesisCmd returns a command to execute genesis state migration.
func MigrateGenesisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [genesis-file]",
		Short: "Migrate 0.38 genesis to the latest version",
		Long: fmt.Sprintf(`Migrate the source genesis into the latest version and print to STDOUT.

Example:
$ %s migrate /path/to/genesis.json --chain-id=FUND-Main-2 --genesis-time=2022-06-01T10:00:00Z
`, version.AppName),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// no migration script
			return nil
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flags.FlagChainID, "", "override chain_id with this flag")

	return cmd
}

// validateGenDoc reads a genesis file and validates that it is a correct
// Tendermint GenesisDoc. This function does not do any cosmos-related
// validation.
func validateGenDoc(importGenesisFile string) (*tmtypes.GenesisDoc, error) {
	genDoc, err := tmtypes.GenesisDocFromFile(importGenesisFile)
	if err != nil {
		return nil, fmt.Errorf("%s. Make sure that"+
			" you have correctly migrated all Tendermint consensus params, please see the"+
			" chain migration guide at %s for more info",
			err.Error(), chainUpgradeGuide,
		)
	}

	return genDoc, nil
}
