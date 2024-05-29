package types

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// DefaultMinCommissionRate is set to 0%
var DefaultBaseValidatorBonus = "0.01"

// NewParams creates a new Params instance
func NewParams(baseValidatorBonus sdk.Dec) Params {
	return Params{
		BaseValidatorBonus: baseValidatorBonus,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	defaultValidatorBonus, _ := sdk.NewDecFromStr(DefaultBaseValidatorBonus)
	return NewParams(defaultValidatorBonus)
}

// Validate validates the set of params
func (p Params) Validate() error {

	if err := validateBaseValidatorBonus(p.BaseValidatorBonus); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateBaseValidatorBonus(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("base validator bonus cannot be negative: %s", v)
	}
	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("base validator bonus cannot be greater than 100%%: %s", v)
	}

	return nil
}
