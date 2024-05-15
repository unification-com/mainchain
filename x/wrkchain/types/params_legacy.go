package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

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
// Deprecated: Type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet
// Deprecated: Type declaration for parameters
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
