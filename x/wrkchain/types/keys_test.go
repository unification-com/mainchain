package types

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrkChainKey(t *testing.T) {
	// key wrkchain
	wcID := uint64(24)
	key := WrkChainKey(wcID)
	require.Equal(t, len(key[1:]), 8)
	require.True(t, bytes.Equal(key[:1], RegisteredWrkChainPrefix))
	wcBz := key[1:]
	wcFromBz := GetWrkChainIDFromBytes(wcBz)
	require.Equal(t, wcFromBz, wcID)
	require.True(t, bytes.Equal(wcBz, GetWrkChainIDBytes(wcID)))
}

func TestWrkChainAllBlocksKey(t *testing.T) {
	wcID := uint64(24)
	key := WrkChainAllBlocksKey(wcID)
	require.Equal(t, len(key[1:]), 8)
	require.True(t, bytes.Equal(key[:1], RecordedWrkChainBlockHashPrefix))
	wcBz := key[1:]
	wcFromBz := GetWrkChainIDFromBytes(wcBz)
	require.Equal(t, wcFromBz, wcID)
	require.True(t, bytes.Equal(wcBz, GetWrkChainIDBytes(wcID)))

}

func TestWrkChainBlockKey(t *testing.T) {
	wcID := uint64(24)
	height := uint64(12345)
	key := WrkChainBlockKey(wcID, height)
	require.Equal(t, len(key[1:]), 16)
	require.True(t, bytes.Equal(key[:1], RecordedWrkChainBlockHashPrefix))

	wcIDbz := key[1:9]
	heightBz := key[9:]

	wcIdFromBz := binary.BigEndian.Uint64(wcIDbz)
	heightFromBz := binary.BigEndian.Uint64(heightBz)

	require.Equal(t, wcIdFromBz, wcID)
	require.Equal(t, heightFromBz, height)
}

func TestWrkChainStorageLimitKey(t *testing.T) {
	wcID := uint64(24)
	key := WrkChainStorageLimitKey(wcID)
	require.Equal(t, len(key[1:]), 8)
	require.True(t, bytes.Equal(key[:1], WrkChainStorageLimitPrefix))
	wcBz := key[1:]
	wcFromBz := GetWrkChainIDFromBytes(wcBz)
	require.Equal(t, wcFromBz, wcID)
	require.True(t, bytes.Equal(wcBz, GetWrkChainIDBytes(wcID)))
}
