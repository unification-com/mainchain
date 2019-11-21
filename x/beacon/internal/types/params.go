package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyDenom = []byte("Denom")
)

// enterprise UND parameters
type Params struct {
	Denom string `json:"denom" yaml:"denom"` // Fee denomination
}

// ParamTable for enterprise UND module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(denom string) Params {
	return Params{
		Denom: denom,
	}
}

// default beacon UND module parameters
func DefaultParams() Params {
	return Params{
		Denom: FeeDenom,
	}
}

// validate params
func ValidateParams(params Params) error {

	if len(params.Denom) == 0 {
		return fmt.Errorf("beacon fee denomination parameter is empty ")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Beacon Params:
  Denomination: %s
`,
		p.Denom,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyDenom, Value: &p.Denom},
	}
}
