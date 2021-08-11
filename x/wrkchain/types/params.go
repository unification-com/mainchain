package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyFeeRegister = []byte("FeeRegister")
	KeyFeeRecord   = []byte("FeeRecord")
	KeyDenom       = []byte("Denom")
)

// ParamTable for BEACON module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(feeReg, feeRec uint64, denom string) Params {
	return Params{
		FeeRegister: feeReg,
		FeeRecord:   feeRec,
		Denom:       denom,
	}
}

// default BEACON module parameters
func DefaultParams() Params {
	return Params{
		FeeRegister: RegFee,
		FeeRecord:   RecordFee,
		Denom:       FeeDenom,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateFeeDenom(p.Denom); err != nil {
		return err
	}

	if err := validateFeeRegister(p.FeeRegister); err != nil {
		return err
	}

	if err := validateFeeRecord(p.FeeRecord); err != nil {
		return err
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeRegister, &p.FeeRegister, validateFeeRegister),
		paramtypes.NewParamSetPair(KeyFeeRecord, &p.FeeRecord, validateFeeRecord),
		paramtypes.NewParamSetPair(KeyDenom, &p.Denom, validateFeeDenom),
	}
}

func validateFeeDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("fee denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateFeeRegister(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("registration fee must be positive: %d", v)
	}

	return nil
}

func validateFeeRecord(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("record fee must be positive: %d", v)
	}

	return nil
}
