package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyEntSigners    = []byte("EntSigners")
	KeyDenom         = []byte("Denom")
	KeyMinAccepts    = []byte("MinAccepts")
	KeyDecisionLimit = []byte("DecisionLimit")
)

// enterprise UND parameters
type Params struct {
	EntSigners    []sdk.AccAddress `json:"ent_signers" yaml:"ent_signers"` // Accounts allowed to sign decisions on UND purchase orders
	Denom         string           `json:"denom" yaml:"denom"`
	MinAccepts    uint64           `json:"min_Accepts" yaml:"min_Accepts"` // must be <= len(EntSigners)
	DecisionLimit uint64           `json:"decision_time_limit" yaml:"decision_time_limit"` // num seconds elapsed before auto-reject
}

// ParamTable for enterprise UND module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(entSigners []sdk.AccAddress, denom string, minAccepts uint64, decisionLimit uint64) Params {

	return Params{
		EntSigners:    entSigners,
		Denom:         denom,
		MinAccepts:    minAccepts,
		DecisionLimit: decisionLimit,
	}
}

// default enterprise UND module parameters
func DefaultParams() Params {
	var entSigners []sdk.AccAddress
	return Params{
		EntSigners:    entSigners,
		Denom:         DefaultDenomination,
		MinAccepts:    1,
		DecisionLimit: 84600,
	}
}

// validate params
func ValidateParams(params Params) error {
	if len(params.EntSigners) == 0 {
		return fmt.Errorf("enterprise und source parameter is empty")
	}
	if len(params.Denom) == 0 {
		return fmt.Errorf("enterprise denomination parameter is empty")
	}
	if params.MinAccepts == 0 {
		return fmt.Errorf("enterprise minimum number of accets parameter must be > 0")
	}
	if params.DecisionLimit == 0 {
		return fmt.Errorf("enterprise decision time limit parameter must be > 0")
	}

	if len(params.EntSigners) < int(params.MinAccepts) {
		return fmt.Errorf("number of authorised accounts must be >= number of minimum accepts")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Enterprise UND Params:
  Source Address: %s
  Denomination: %s
  MinAccepts: %d
  DecisionLimit: %d
`,
		p.EntSigners, p.Denom, p.MinAccepts, p.DecisionLimit,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyEntSigners, Value: &p.EntSigners},
		{Key: KeyDenom, Value: &p.Denom},
		{Key: KeyMinAccepts, Value: &p.MinAccepts},
		{Key: KeyDecisionLimit, Value: &p.DecisionLimit},
	}
}
