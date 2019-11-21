package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEqualPurchaseOrderID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	state1.StartingPurchaseOrderID = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	state2.StartingPurchaseOrderID = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}
