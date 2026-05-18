package types

import (
	"fmt"

	mathmod "cosmossdk.io/math"
)

// DefaultValidatorFee is the default 1% fee skimmed from each claim and sent
// to the fee collector. Gov can tune via MsgUpdateParams subject to the
// MaxValidatorFee cap.
var DefaultValidatorFee = mathmod.LegacyNewDecWithPrec(1, 2)

// MaxValidatorFee is the hard upper bound on the gov-tunable ValidatorFee.
// Set to 10% — anything higher would let gov effectively expropriate stream
// senders, and is out of scope for the module's intended fee policy.
var MaxValidatorFee = mathmod.LegacyNewDecWithPrec(10, 2)

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
	return validateValidatorFee(p.ValidatorFee)
}

func validateValidatorFee(v mathmod.LegacyDec) error {
	if v.IsNil() {
		return fmt.Errorf("validator fee cannot be nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("validator fee cannot be negative: %s", v)
	}
	if v.GT(MaxValidatorFee) {
		return fmt.Errorf("validator fee cannot exceed %s. Sent %s", MaxValidatorFee, v)
	}
	return nil
}
