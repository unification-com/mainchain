package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// Parameter store keys
var (
	KeyEntSigners    = []byte("EntSigners")
	KeyDenom         = []byte("Denom")
	KeyMinAccepts    = []byte("MinAccepts")
	KeyDecisionLimit = []byte("DecisionLimit")
)

// ParamTable for enterprise FUND module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEntSigners, &p.EntSigners, validateEntSigners),
		paramtypes.NewParamSetPair(KeyDenom, &p.Denom, validateDenom),
		paramtypes.NewParamSetPair(KeyMinAccepts, &p.MinAccepts, validateMinAccepts),
		paramtypes.NewParamSetPair(KeyDecisionLimit, &p.DecisionTimeLimit, validateDecisionLimit),
	}
}
