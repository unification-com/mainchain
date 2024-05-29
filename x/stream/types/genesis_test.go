package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
)

func TestGenesisState_Validate(t *testing.T) {
	invalidFee, _ := sdk.NewDecFromStr("1.01")

	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params:           types.DefaultParams(),
				StartingStreamId: 1,
				TotalDeposits:    types.TotalDeposits{},
				Streams:          []types.Stream{},
			},
			valid: true,
		},
		{
			desc: "invalid genesis state fee > 1.0",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: invalidFee,
				},
				StartingStreamId: 1,
				TotalDeposits:    types.TotalDeposits{},
				Streams:          []types.Stream{},
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
