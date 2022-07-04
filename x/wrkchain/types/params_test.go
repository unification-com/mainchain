package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewParams(t *testing.T) {
	params := NewParams(RegFee, RecordFee, PurchaseStorageFee, FeeDenom, DefaultStorageLimit, DefaultMaxStorageLimit)

	require.Equal(t, uint64(RegFee), params.FeeRegister)
	require.Equal(t, uint64(RecordFee), params.FeeRecord)
	require.Equal(t, uint64(PurchaseStorageFee), params.FeePurchaseStorage)
	require.Equal(t, uint64(DefaultStorageLimit), params.DefaultStorageLimit)
	require.Equal(t, uint64(DefaultMaxStorageLimit), params.MaxStorageLimit)
	require.Equal(t, FeeDenom, params.Denom)
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	require.Equal(t, uint64(RegFee), params.FeeRegister)
	require.Equal(t, uint64(RecordFee), params.FeeRecord)
	require.Equal(t, uint64(PurchaseStorageFee), params.FeePurchaseStorage)
	require.Equal(t, uint64(DefaultStorageLimit), params.DefaultStorageLimit)
	require.Equal(t, uint64(DefaultMaxStorageLimit), params.MaxStorageLimit)
	require.Equal(t, FeeDenom, params.Denom)
}

func TestParams_Validate(t *testing.T) {
	params1 := DefaultParams()
	err := params1.Validate()
	require.NoError(t, err)

	params2 := NewParams(RegFee, RecordFee, PurchaseStorageFee, FeeDenom, DefaultStorageLimit, DefaultMaxStorageLimit)
	err = params2.Validate()
	require.NoError(t, err)

	params3 := Params{}
	err = params3.Validate()
	require.Equal(t, "fee denom cannot be blank", err.Error())

	params3.Denom = "test"
	err = params3.Validate()
	require.Equal(t, "registration fee must be positive: 0", err.Error())

	params3.FeeRegister = 1
	err = params3.Validate()
	require.Equal(t, "record fee must be positive: 0", err.Error())

	params3.FeeRecord = 1
	err = params3.Validate()
	require.Equal(t, "purchase storage fee must be positive: 0", err.Error())

	params3.FeePurchaseStorage = 1
	err = params3.Validate()
	require.Equal(t, "default storage must be positive: 0", err.Error())

	params3.DefaultStorageLimit = 100
	err = params3.Validate()
	require.Equal(t, "max storage must be positive: 0", err.Error())

	params3.MaxStorageLimit = 1
	err = params3.Validate()
	require.Equal(t, "default storage 100 > max storage 1", err.Error())

	params3.MaxStorageLimit = 101
	err = params3.Validate()
	require.NoError(t, err)

	params3.MaxStorageLimit = 100
	err = params3.Validate()
	require.NoError(t, err)
}
