package types

// GenesisState - minter state
type GenesisState struct {
	Params                  Params `json:"params" yaml:"params"` // enterprise params
	StartingPurchaseOrderID uint64 `json:"starting_purchase_order_id" yaml:"starting_purchase_order_id"` // should be 1
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingPurchaseOrderID uint64) GenesisState {
	return GenesisState{
		Params:                  params,
		StartingPurchaseOrderID: startingPurchaseOrderID,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:                  DefaultParams(),
		StartingPurchaseOrderID: DefaultStartingPurchaseOrderID,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := ValidateParams(data.Params)
	if err != nil {
		return err
	}

	return nil
}

