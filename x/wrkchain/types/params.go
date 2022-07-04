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
	KeyFeeRegister         = []byte("FeeRegister")
	KeyFeeRecord           = []byte("FeeRecord")
	KeyFeePurchaseStorage  = []byte("FeePurchaseStorage")
	KeyDenom               = []byte("Denom")
	KeyDefaultStorageLimit = []byte("DefaultStorageLimit")
	KeyMaxStorageLimit     = []byte("MaxStorageLimit")
)

// ParamTable for BEACON module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(feeReg, feeRec, feePurchase uint64, denom string, defaultStorage, maxStorage uint64) Params {
	return Params{
		FeeRegister:         feeReg,
		FeeRecord:           feeRec,
		FeePurchaseStorage:  feePurchase,
		Denom:               denom,
		DefaultStorageLimit: defaultStorage,
		MaxStorageLimit:     maxStorage,
	}
}

// default BEACON module parameters
func DefaultParams() Params {
	return Params{
		FeeRegister:         RegFee,
		FeeRecord:           RecordFee,
		FeePurchaseStorage:  PurchaseStorageFee,
		Denom:               FeeDenom,
		DefaultStorageLimit: DefaultStorageLimit,
		MaxStorageLimit:     DefaultMaxStorageLimit,
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

	if err := validateFeePurchaseStorage(p.FeePurchaseStorage); err != nil {
		return err
	}

	if err := validateDefaultStorageLimit(p.DefaultStorageLimit); err != nil {
		return err
	}

	if err := validateMaxStorageLimit(p.MaxStorageLimit); err != nil {
		return err
	}

	if p.DefaultStorageLimit > p.MaxStorageLimit {
		return fmt.Errorf("default storage %d > max storage %d", p.DefaultStorageLimit, p.MaxStorageLimit)
	}

	return nil
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyFeeRegister, &p.FeeRegister, validateFeeRegister),
		paramtypes.NewParamSetPair(KeyFeeRecord, &p.FeeRecord, validateFeeRecord),
		paramtypes.NewParamSetPair(KeyFeePurchaseStorage, &p.FeePurchaseStorage, validateFeePurchaseStorage),
		paramtypes.NewParamSetPair(KeyDenom, &p.Denom, validateFeeDenom),
		paramtypes.NewParamSetPair(KeyDefaultStorageLimit, &p.DefaultStorageLimit, validateDefaultStorageLimit),
		paramtypes.NewParamSetPair(KeyMaxStorageLimit, &p.MaxStorageLimit, validateMaxStorageLimit),
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

func validateFeePurchaseStorage(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("purchase storage fee must be positive: %d", v)
	}

	return nil
}

func validateDefaultStorageLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("default storage must be positive: %d", v)
	}

	return nil
}

func validateMaxStorageLimit(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max storage must be positive: %d", v)
	}

	return nil
}
