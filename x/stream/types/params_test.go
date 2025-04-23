package types_test

import (
	"testing"

	mathmod "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/stream/types"
)

func TestParamsValidate(t *testing.T) {
	params1 := types.Params{ValidatorFee: mathmod.LegacyNewDecWithPrec(1, 2)}
	err := params1.Validate()
	require.NoError(t, err)

	params2 := types.Params{ValidatorFee: mathmod.LegacyNewDecWithPrec(-1, 2)}
	err = params2.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator fee cannot be negative:")

	params3 := types.Params{ValidatorFee: mathmod.LegacyNewDecWithPrec(101, 2)}
	err = params3.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator fee cannot be greater than 100% (1.00). Sent")

	params4 := types.Params{ValidatorFee: mathmod.LegacyDec{}}
	err = params4.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "validator fee cannot be nil")

}
