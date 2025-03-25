package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/stream/types"
)

func TestParamsValidate(t *testing.T) {
	params1 := types.Params{ValidatorFee: "0.01"}
	err := params1.Validate()
	require.NoError(t, err)

	params2 := types.Params{ValidatorFee: "-0.01"}
	err = params2.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator fee cannot be negative:")

	params3 := types.Params{ValidatorFee: "1.01"}
	err = params3.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator fee cannot be greater than 100% (1.00). Sent")

	params4 := types.Params{ValidatorFee: ""}
	err = params4.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "decimal string cannot be empty")

}
