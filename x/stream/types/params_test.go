package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
)

func TestParamsValidate(t *testing.T) {
	params1 := types.Params{ValidatorFee: sdk.MustNewDecFromStr("0.01")}
	err := params1.Validate()
	require.NoError(t, err)

	params2 := types.Params{ValidatorFee: sdk.MustNewDecFromStr("-0.01")}
	err = params2.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "base validator fee cannot be negative:")

	params3 := types.Params{ValidatorFee: sdk.MustNewDecFromStr("1.01")}
	err = params3.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "base validator fee cannot be greater than 100% (1.00). Sent")

	params4 := types.Params{ValidatorFee: sdk.Dec{}}
	err = params4.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "base validator fee cannot be nil")
}
