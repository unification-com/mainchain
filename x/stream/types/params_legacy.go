package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

var (
	KeyValidatorFee = []byte("KeyValidatorFee")
)

// ParamKeyTable the param key table for launch module
// Deprecated: Type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyValidatorFee, &p.ValidatorFee, validateBaseValidatorFee),
	}
}
