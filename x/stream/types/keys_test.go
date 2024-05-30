package types_test

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
)

func TestGetStreamIdLookupKey(t *testing.T) {
	sId := uint64(246973)
	key := types.GetStreamIdLookupKey(sId)
	require.True(t, len(key[1:]) == 8)
	require.True(t, bytes.Equal(key[:1], types.StreamIdLookupKeyPrefix))
	bBz := key[1:]
	bFromBz := types.GetStreamIdFromBytes(bBz)
	require.True(t, bFromBz == sId)
	require.True(t, bytes.Equal(bBz, types.GetStreamIdBytes(sId)))
}
