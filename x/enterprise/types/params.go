package types

import (
	"errors"
	"fmt"
	undtypes "github.com/unification-com/mainchain/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewParams(denom string, minAccepts uint64, decisionLimit uint64, entSigners string) Params {
	return Params{
		EntSigners:        entSigners,
		Denom:             denom,
		MinAccepts:        minAccepts,
		DecisionTimeLimit: decisionLimit,
	}
}

// default enterprise FUND module parameters
func DefaultParams() Params {
	return Params{
		EntSigners:        "und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm", // default to 0000000000000000000000000000000000000000
		Denom:             undtypes.DefaultDenomination,
		MinAccepts:        1,
		DecisionTimeLimit: 84600,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateDenom(p.Denom); err != nil {
		return err
	}

	if err := validateMinAccepts(p.MinAccepts); err != nil {
		return err
	}

	if err := validateDecisionLimit(p.DecisionTimeLimit); err != nil {
		return err
	}

	if err := validateEntSigners(p.EntSigners); err != nil {
		return err
	}

	entSigners := strings.Split(p.EntSigners, ",")

	if len(entSigners) < int(p.MinAccepts) {
		return fmt.Errorf("number of authorised accounts must be >= number of minimum accepts")
	}

	return nil
}

func validateDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateMinAccepts(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("min accepts must be positive: %d", v)
	}

	return nil
}

func validateDecisionLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("decision limit must be positive: %d", v)
	}

	return nil
}

func validateEntSigners(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("must have at least one signer")
	}

	entSigners := strings.Split(v, ",")

	if len(entSigners) == 0 {
		return fmt.Errorf("must have at least one signer")
	}

	for _, authAddr := range entSigners {
		_, err := sdk.AccAddressFromBech32(authAddr)
		if err != nil {
			return fmt.Errorf("invalid address %s: %s", authAddr, err)
		}
	}

	return nil
}
