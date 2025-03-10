package types_test

import (
	"testing"
	"time"

	mathmod "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/stream/types"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		expErr   bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			expErr:   false,
		},
		{
			desc: "valid genesis state default params",
			genState: &types.GenesisState{
				Params:  types.DefaultParams(),
				Streams: []types.StreamExport{},
			},
			expErr: false,
		},
		{
			desc: "valid genesis state and params",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: mathmod.LegacyNewDecWithPrec(1, 2),
				},
				Streams: []types.StreamExport{
					{
						Sender:   "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
						Receiver: "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
						Stream: types.Stream{
							Deposit:         sdk.NewCoin("nund", mathmod.NewIntFromUint64(100000000)),
							FlowRate:        123,
							LastOutflowTime: time.Now(),
							DepositZeroTime: time.Now(),
							Cancellable:     true,
						},
					},
					{
						Sender:   "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
						Receiver: "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
						Stream: types.Stream{
							Deposit:         sdk.NewCoin("nund", mathmod.NewIntFromUint64(200000000)),
							FlowRate:        321,
							LastOutflowTime: time.Now(),
							DepositZeroTime: time.Now(),
							Cancellable:     false,
						},
					},
				},
			},
			expErr: false,
		},
		{
			desc: "invalid: empty genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: mathmod.LegacyDec{},
				},
				Streams: []types.StreamExport{},
			},
			expErr: true,
		},
		{
			desc: "invalid: genesis state fee > 100%",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: mathmod.LegacyNewDecWithPrec(101, 2),
				},
				Streams: []types.StreamExport{},
			},
			expErr: true,
		},
		{
			desc: "invalid: genesis state fee < 0%",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: mathmod.LegacyNewDecWithPrec(-1, 2),
				},
				Streams: []types.StreamExport{},
			},
			expErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
