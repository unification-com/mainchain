package types

// GenesisState - beacon state
type GenesisState struct {
	Params Params `json:"params" yaml:"params"` // beacon params
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
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
