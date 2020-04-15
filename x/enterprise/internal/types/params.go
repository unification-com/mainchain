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
	KeyEntSigners    = []byte("EntSigners")
	KeyDenom         = []byte("Denom")
	KeyMinAccepts    = []byte("MinAccepts")
	KeyDecisionLimit = []byte("DecisionLimit")
)

// enterprise UND parameters
type Params struct {
	EntSigners    string `json:"ent_signers" yaml:"ent_signers"` // Accounts allowed to sign decisions on UND purchase orders
	Denom         string `json:"denom" yaml:"denom"`
	MinAccepts    uint64 `json:"min_Accepts" yaml:"min_Accepts"`                 // must be <= len(EntSigners)
	DecisionLimit uint64 `json:"decision_time_limit" yaml:"decision_time_limit"` // num seconds elapsed before auto-reject
}

// ParamTable for enterprise UND module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(denom string, minAccepts uint64, decisionLimit uint64, entSigners string) Params {
	return Params{
		EntSigners:    entSigners,
		Denom:         denom,
		MinAccepts:    minAccepts,
		DecisionLimit: decisionLimit,
	}
}

// default enterprise UND module parameters
func DefaultParams() Params {
	return Params{
		EntSigners:    "und1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq5x8kpm", // default to 0000000000000000000000000000000000000000
		Denom:         DefaultDenomination,
		MinAccepts:    1,
		DecisionLimit: 84600,
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

	if err := validateDecisionLimit(p.DecisionLimit); err != nil {
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
		params.NewParamSetPair(KeyEntSigners, &p.EntSigners, validateEntSigners),
		params.NewParamSetPair(KeyDenom, &p.Denom, validateDenom),
		params.NewParamSetPair(KeyMinAccepts, &p.MinAccepts, validateMinAccepts),
		params.NewParamSetPair(KeyDecisionLimit, &p.DecisionLimit, validateDecisionLimit),
	}
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
