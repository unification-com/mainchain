package types

import (
	"bytes"
	"fmt"
)

// GenesisState - wrkchain state
type GenesisState struct {
	Params             Params           `json:"params" yaml:"params"`                             // wrkchain params
	StartingWrkChainID uint64           `json:"starting_wrkchain_id" yaml:"starting_wrkchain_id"` // should be 1
	WrkChains          []WrkChainExport `json:"registered_wrkchains" yaml:"registered_wrkchains"`
}

type WrkChainExport struct {
	WrkChain       WrkChain        `json:"wrkchain" yaml:"wrkchain"`
	WrkChainBlocks []WrkChainBlock `json:"blocks" yaml:"blocks"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingWrkChainID uint64) GenesisState {
	return GenesisState{
		Params:             params,
		StartingWrkChainID: startingWrkChainID,
		WrkChains:          nil,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:             DefaultParams(),
		StartingWrkChainID: DefaultStartingWrkChainID,
	}
}

// Equal checks whether two wrkchain GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.WrkChains {
		if record.WrkChain.WrkChainID == 0 {
			return fmt.Errorf("invalid WrkChain: ID: %d. Error: Missing ID", record.WrkChain.WrkChainID)
		}
		if record.WrkChain.Owner == nil {
			return fmt.Errorf("invalid WrkChain: Owner: %s. Error: Missing Owner", record.WrkChain.Owner)
		}
		if record.WrkChain.Moniker == "" {
			return fmt.Errorf("invalid WrkChain: Moniker: %s. Error: Missing Moniker", record.WrkChain.Moniker)
		}
		if record.WrkChain.BaseType == "" {
			return fmt.Errorf("invalid WrkChain: BaseType: %s. Error: Missing BaseType", record.WrkChain.BaseType)
		}
		for _, block := range record.WrkChainBlocks {
			if block.BlockHash == "" {
				return fmt.Errorf("invalid WrkChain block: BlockHash: %s. Error: Missing BlockHash", block.BlockHash)
			}
			if block.Height == 0 {
				return fmt.Errorf("invalid WrkChain block: Height: %d. Error: Missing Height", block.Height)
			}
		}
	}
	return nil
}
