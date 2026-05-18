package types

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaxAllowedStorageLimit is the absolute upper bound for the gov-tunable
// MaxStorageLimit parameter. 10 million slots × ~150 bytes per slot ≈ 1.5GB
// per WRKChain at the absolute limit; already excessive for any realistic
// anchor-tracking use case. Gov cannot push past this without a code change.
const MaxAllowedStorageLimit uint64 = 10_000_000

// MaxWrkChainsPerOwner caps the number of WRKChains any single owner address
// can hold. The per-registration FeeRegister is a deterrent but not a hard
// cap; a well-funded attacker could otherwise inflate state with many small
// WRKChains.
const MaxWrkChainsPerOwner = 100

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

func validateFeeDenom(v string) error {
	if strings.TrimSpace(v) == "" {
		return errors.New("fee denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}
	return nil
}

func validateFeeRegister(v uint64) error {
	if v == 0 {
		return fmt.Errorf("registration fee must be positive: %d", v)
	}
	return nil
}

func validateFeeRecord(v uint64) error {
	if v == 0 {
		return fmt.Errorf("record fee must be positive: %d", v)
	}
	return nil
}

func validateFeePurchaseStorage(v uint64) error {
	if v == 0 {
		return fmt.Errorf("purchase storage fee must be positive: %d", v)
	}
	return nil
}

func validateDefaultStorageLimit(v uint64) error {
	if v == 0 {
		return fmt.Errorf("default storage must be positive: %d", v)
	}
	return nil
}

func validateMaxStorageLimit(v uint64) error {
	if v == 0 {
		return fmt.Errorf("max storage must be positive: %d", v)
	}
	if v > MaxAllowedStorageLimit {
		return fmt.Errorf("max storage %d exceeds hard cap %d", v, MaxAllowedStorageLimit)
	}
	return nil
}
