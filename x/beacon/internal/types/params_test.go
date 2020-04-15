package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewParams(t *testing.T) {
	params := NewParams(RegFee, RecordFee, FeeDenom)

	require.Equal(t, uint64(RegFee), params.FeeRegister)
	require.Equal(t, uint64(RecordFee), params.FeeRecord)
	require.Equal(t, FeeDenom, params.Denom)
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	require.Equal(t, uint64(RegFee), params.FeeRegister)
	require.Equal(t, uint64(RecordFee), params.FeeRecord)
	require.Equal(t, FeeDenom, params.Denom)
}

func TestParams_Validate(t *testing.T) {
	params1 := DefaultParams()
	err := params1.Validate()
	require.NoError(t, err)

	params2 := NewParams(RegFee, RecordFee, FeeDenom)
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
	require.NoError(t, err)
}
