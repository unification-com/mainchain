package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"
	"time"

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
				Params:           types.DefaultParams(),
				StartingStreamId: 1,
				Streams:          []types.Stream{},
			},
			expErr: false,
		},
		{
			desc: "valid genesis state and params",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: sdk.MustNewDecFromStr("0.01"),
				},
				StartingStreamId: 3,
				Streams: []types.Stream{
					{
						StreamId:        1,
						Sender:          "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
						Receiver:        "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
						Deposit:         sdk.NewCoin("nund", sdk.NewIntFromUint64(100000000)),
						FlowRate:        123,
						CreateTime:      time.Now(),
						LastUpdatedTime: time.Now(),
						LastOutflowTime: time.Now(),
						DepositZeroTime: time.Now(),
						TotalStreamed:   sdk.NewCoin("nund", sdk.NewIntFromUint64(1234)),
						Cancellable:     true,
					},
					{
						StreamId:        2,
						Sender:          "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
						Receiver:        "und1x8pl6wzqf9atkm77ymc5vn5dnpl5xytmn200xy",
						Deposit:         sdk.NewCoin("nund", sdk.NewIntFromUint64(200000000)),
						FlowRate:        321,
						CreateTime:      time.Now(),
						LastUpdatedTime: time.Now(),
						LastOutflowTime: time.Now(),
						DepositZeroTime: time.Now(),
						TotalStreamed:   sdk.NewCoin("nund", sdk.NewIntFromUint64(4321)),
						Cancellable:     false,
					},
				},
			},
			expErr: false,
		},
		{
			desc: "invalid: empty genesis state",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: sdk.Dec{},
				},
				StartingStreamId: 0,
				Streams:          []types.Stream{},
			},
			expErr: true,
		},
		{
			desc: "invalid: genesis state fee > 100%",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: sdk.MustNewDecFromStr("1.01"),
				},
				StartingStreamId: 1,
				Streams:          []types.Stream{},
			},
			expErr: true,
		},
		{
			desc: "invalid: genesis state fee < 0%",
			genState: &types.GenesisState{
				Params: types.Params{
					ValidatorFee: sdk.MustNewDecFromStr("-0.01"),
				},
				StartingStreamId: 1,
				Streams:          []types.Stream{},
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
