package types

import (
	"fmt"

	mathmod "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// DefaultValidatorFee is set to 0%
var DefaultValidatorFee = mathmod.LegacyNewDecWithPrec(1, 2)

// NewParams creates a new Params instance
func NewParams(validatorFee mathmod.LegacyDec) Params {
	return Params{
		ValidatorFee: validatorFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultValidatorFee)
}

// Validate validates the set of params
func (p Params) Validate() error {

	if err := validateBaseValidatorFee(p.ValidatorFee); err != nil {
		return err
	}

	return nil
}

func validateBaseValidatorFee(i interface{}) error {
	v, ok := i.(mathmod.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("validator fee cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("validator fee cannot be negative: %s", v)
	}

	if v.GT(mathmod.LegacyOneDec()) {
		return fmt.Errorf("validator fee cannot be greater than 100%% (1.00). Sent %s", v)
	}

	return nil
}
