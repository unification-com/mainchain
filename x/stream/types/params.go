package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// DefaultValidatorFee is set to 0%
var DefaultValidatorFee = "0.01"

// NewParams creates a new Params instance
func NewParams(validatorFee sdk.Dec) Params {
	return Params{
		ValidatorFee: validatorFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	defaultValidatorFee, _ := sdk.NewDecFromStr(DefaultValidatorFee)
	return NewParams(defaultValidatorFee)
}

// Validate validates the set of params
func (p Params) Validate() error {

	if err := validateBaseValidatorFee(p.ValidatorFee); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateBaseValidatorFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("validator fee cannot be nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("validator fee cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("validator fee cannot be greater than 100%% (1.00). Sent %s", v)
	}

	return nil
}
