package types_test

import (
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"testing"
)

func TestRaisedQueueStoreKey(t *testing.T) {
	for i := uint64(1); i <= 999999; i++ {
		key := types.RaisedQueueStoreKey(i)
		poId := types.SplitRaisedQueueKey(key)
		require.Equal(t, i, poId)
	}
}

func TestAcceptedQueueStoreKey(t *testing.T) {
	for i := uint64(1); i <= 999999; i++ {
		key := types.AcceptedQueueStoreKey(i)
		poId := types.SplitAcceptedQueueKey(key)
		require.Equal(t, i, poId)
	}
}
