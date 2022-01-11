package types

import (
	"fmt"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingWrkChainID uint64) *GenesisState {
	return &GenesisState{
		Params:              params,
		StartingWrkchainId:  startingWrkChainID,
		RegisteredWrkchains: nil,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		StartingWrkchainId:  DefaultStartingWrkChainID,
		RegisteredWrkchains: nil,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.RegisteredWrkchains {
		if record.Wrkchain.WrkchainId == 0 {
			return fmt.Errorf("invalid WrkChain: ID: %d. Error: Missing ID", record.Wrkchain.WrkchainId)
		}
		if record.Wrkchain.Owner == "" {
			return fmt.Errorf("invalid WrkChain: Owner: %s. Error: Missing Owner", record.Wrkchain.Owner)
		}
		if record.Wrkchain.Moniker == "" {
			return fmt.Errorf("invalid WrkChain: Moniker: %s. Error: Missing Moniker", record.Wrkchain.Moniker)
		}
		if record.Wrkchain.Type == "" {
			return fmt.Errorf("invalid WrkChain: BaseType: %s. Error: Missing BaseType", record.Wrkchain.Type)
		}
		for _, block := range record.Blocks {
			if block.Bh == "" {
				return fmt.Errorf("invalid WrkChain block: BlockHash: %s. Error: Missing BlockHash", block.Bh)
			}
			if block.He == 0 {
				return fmt.Errorf("invalid WrkChain block: Height: %d. Error: Missing Height", block.He)
			}
		}
	}
	return nil
}
