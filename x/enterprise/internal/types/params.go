package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyEntSource = []byte("EntSource")
)

// enterprise UND parameters
type Params struct {
	EntSource sdk.AccAddress `json:"ent_source" yaml:"ent_source"` // Acc allowed to sign and raise UND purchase orders
}

// ParamTable for enterprise UND module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(entSource sdk.AccAddress) Params {

	return Params{
		EntSource: entSource,
	}
}

// default enterprise UND module parameters
func DefaultParams() Params {
	return Params{
		EntSource: sdk.AccAddress{},
	}
}

// validate params
func ValidateParams(params Params) error {
	if params.EntSource.Empty() {
		return fmt.Errorf("enterprise und source parameter is empty ")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Enterprise UND Params:
  Source Address: %s
`,
		p.EntSource,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyEntSource, Value: &p.EntSource},
	}
}
