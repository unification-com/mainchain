package simulation_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/wrkchain/simulation"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

var (
	bPk1   = ed25519.GenPrivKey().PubKey()
	bAddr1 = sdk.AccAddress(bPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	testApp := test_helpers.Setup(false)
	cdc := testApp.AppCodec()
	dec := simulation.NewDecodeStore(cdc)

	wc, err := types.NewWrkchain(1, "beacon1", "Test BEACON 1", "gen", "test", 0, 0, uint64(time.Now().Unix()), bAddr1.String())
	require.NoError(t, err)

	block, err := types.NewWrkchainBlock(1, 1, "bhash", "phash", "h1", "h2", "h3", uint64(time.Now().Unix()), bAddr1.String())
	require.NoError(t, err)

	wcBz, err := cdc.MarshalBinaryBare(&wc)
	require.NoError(t, err)

	blockBz, err := cdc.MarshalBinaryBare(&block)
	require.NoError(t, err)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.WrkChainKey(1), Value: wcBz},
			{Key: types.WrkChainBlockKey(1, 1), Value: blockBz},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"beacon", fmt.Sprintf("%v\n%v", wc, wc)},
		{"beacon timestamp", fmt.Sprintf("%v\n%v", block, block)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
