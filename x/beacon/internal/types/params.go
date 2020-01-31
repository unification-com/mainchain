package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyFeeRegister = []byte("FeeRegister")
	KeyFeeRecord   = []byte("FeeRecord")
	KeyDenom       = []byte("Denom")
)

// beacon parameters
type Params struct {
	FeeRegister uint64 `json:"fee_register" yaml:"fee_register"` // Fee for registering a BEACON
	FeeRecord   uint64 `json:"fee_record" yaml:"fee_record"`     // Fee for recording timestamps for a BEACON
	Denom       string `json:"denom" yaml:"denom"`               // Fee denomination
}

// ParamTable for BEACON module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
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

func (p Params) String() string {
	return fmt.Sprintf(`WRKChain Params:
  Registration Fee: %d
  Recording Fee: %d
  Denomination: %s
`,
		p.FeeRegister, p.FeeRecord, p.Denom,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyFeeRegister, &p.FeeRegister, validateFeeRegister),
		params.NewParamSetPair(KeyFeeRecord, &p.FeeRecord, validateFeeRecord),
		params.NewParamSetPair(KeyDenom, &p.Denom, validateFeeDenom),
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
