package types

import (
	"fmt"

	mathmod "cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// DefaultValidatorFee is set to 0%
var DefaultValidatorFee = "0.01"

// NewParams creates a new Params instance
func NewParams(validatorFee string) Params {
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
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	vDec, err := mathmod.LegacyNewDecFromStr(v)
	if err != nil {
		return fmt.Errorf("invalid validator fee string: %w", err)
	}

	if vDec.IsNil() {
		return fmt.Errorf("validator fee cannot be nil")
	}

	if vDec.IsNegative() {
		return fmt.Errorf("validator fee cannot be negative: %s", v)
	}

	if vDec.GT(mathmod.LegacyOneDec()) {
		return fmt.Errorf("validator fee cannot be greater than 100%% (1.00). Sent %s", v)
	}

	return nil
}
