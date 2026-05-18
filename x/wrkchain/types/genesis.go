package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// WRKChain fees, in nano FUND
	RegFee             = 1000000000000 // 1000 FUND - used in init genesis
	RecordFee          = 1000000000    // 1 FUND - used in init genesis
	PurchaseStorageFee = 5000000000    // 5 FUND - used in init genesis

	DefaultStartingWrkChainID uint64 = 1 // used in init genesis

	DefaultStorageLimit    uint64 = 50000  // used in init genesis
	DefaultMaxStorageLimit uint64 = 600000 // used in init genesis
)

var (
	FeeDenom = sdk.DefaultBondDenom // used in init genesis
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingWrkChainID uint64, wrkChains WrkChainExports) *GenesisState {
	return &GenesisState{
		Params:              params,
		StartingWrkchainId:  startingWrkChainID,
		RegisteredWrkchains: wrkChains,
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
// expected state holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	seenIDs := make(map[uint64]struct{}, len(data.RegisteredWrkchains))
	for i, record := range data.RegisteredWrkchains {
		if record.Wrkchain.WrkchainId == 0 {
			return fmt.Errorf("wrkchain[%d]: invalid WrkchainId 0", i)
		}
		if _, dup := seenIDs[record.Wrkchain.WrkchainId]; dup {
			return fmt.Errorf("wrkchain[%d]: duplicate WrkchainId %d", i, record.Wrkchain.WrkchainId)
		}
		seenIDs[record.Wrkchain.WrkchainId] = struct{}{}
		if record.Wrkchain.WrkchainId >= data.StartingWrkchainId {
			return fmt.Errorf("wrkchain[%d]: WrkchainId %d must be < StartingWrkchainId %d",
				i, record.Wrkchain.WrkchainId, data.StartingWrkchainId)
		}
		if _, err := sdk.AccAddressFromBech32(record.Wrkchain.Owner); err != nil {
			return fmt.Errorf("wrkchain[%d]: invalid Owner %q: %w", i, record.Wrkchain.Owner, err)
		}
		if record.Wrkchain.Moniker == "" {
			return fmt.Errorf("wrkchain[%d]: empty Moniker", i)
		}
		if record.Wrkchain.BaseType == "" {
			return fmt.Errorf("wrkchain[%d]: empty BaseType", i)
		}
		if record.InStateLimit == 0 {
			return fmt.Errorf("wrkchain[%d]: zero InStateLimit", i)
		}
		seenHeights := make(map[uint64]struct{}, len(record.Blocks))
		for j, block := range record.Blocks {
			if block.Bh == "" {
				return fmt.Errorf("wrkchain[%d].block[%d]: empty BlockHash", i, j)
			}
			if block.He == 0 {
				return fmt.Errorf("wrkchain[%d].block[%d]: zero Height", i, j)
			}
			if _, dup := seenHeights[block.He]; dup {
				return fmt.Errorf("wrkchain[%d].block[%d]: duplicate Height %d", i, j, block.He)
			}
			seenHeights[block.He] = struct{}{}
		}
	}
	return nil
}
