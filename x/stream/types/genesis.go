package types

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		Streams: []StreamExport{},
	}
}

func NewGenesisState(streams []StreamExport, params Params) *GenesisState {
	return &GenesisState{
		Params:  params,
		Streams: streams,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
