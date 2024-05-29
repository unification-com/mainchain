package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyBaseValidatorBonus = []byte("KeyBaseValidatorBonus")
)

// ParamKeyTable the param key table for launch module
// Deprecated: Type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBaseValidatorBonus, &p.BaseValidatorBonus, validateBaseValidatorBonus),
	}
}
