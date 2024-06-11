package simulation_test

import (
	"fmt"
	"github.com/unification-com/mainchain/x/stream"
	"github.com/unification-com/mainchain/x/stream/simulation"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

func TestDecodeStore(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(stream.AppModuleBasic{})

	dec := simulation.NewDecodeStore(encCfg.Codec)

	now := time.Now().UTC()

	newStream := types.Stream{
		Deposit:         sdk.NewInt64Coin("stake", 1000),
		FlowRate:        1,
		LastOutflowTime: now,
		DepositZeroTime: now,
		Cancellable:     true,
	}

	streamBz, err := encCfg.Codec.Marshal(&newStream)
	require.NoError(t, err)
	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: []byte(types.StreamKeyPrefix), Value: streamBz},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectErr   bool
		expectedLog string
	}{
		{"Stream", false, fmt.Sprintf("%v\n%v", newStream, newStream)},
		{"other", true, ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
