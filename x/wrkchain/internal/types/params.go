package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyFeeRegister = []byte("FeeRegister")
	KeyFeeRecord   = []byte("FeeRecord")
	KeyDenom       = []byte("Denom")
)

// wrkchain parameters
type Params struct {
	FeeRegister uint64 `json:"fee_register" yaml:"fee_register"` // Fee for registering a WRKChain
	FeeRecord   uint64 `json:"fee_record" yaml:"fee_record"`     // Fee for recording hashes for a WRKChain
	Denom       string `json:"denom" yaml:"denom"`               // Fee denomination
}

// ParamTable for wrkchain module.
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

// default wrkchain module parameters
func DefaultParams() Params {
	return Params{
		FeeRegister: RegFee,
		FeeRecord:   RecordFee,
		Denom:       FeeDenom,
	}
}

// validate params
func ValidateParams(params Params) error {

	if len(params.Denom) == 0 {
		return fmt.Errorf("wrkchain fee denomination parameter is empty ")
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
		{Key: KeyFeeRegister, Value: &p.FeeRegister},
		{Key: KeyFeeRecord, Value: &p.FeeRecord},
		{Key: KeyDenom, Value: &p.Denom},
	}
}
