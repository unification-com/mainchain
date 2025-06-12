package types

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBeaconKey(t *testing.T) {
	// key beacon
	bID := uint64(24)
	key := BeaconKey(bID)
	require.True(t, len(key[1:]) == 8)
	require.True(t, bytes.Equal(key[:1], RegisteredBeaconPrefix))
	bBz := key[1:]
	bFromBz := GetBeaconIDFromBytes(bBz)
	require.True(t, bFromBz == bID)
	require.True(t, bytes.Equal(bBz, GetBeaconIDBytes(bID)))
}

func TestBeaconAllTimestampsKey(t *testing.T) {
	bID := uint64(24)
	key := BeaconAllTimestampsKey(bID)
	require.True(t, len(key[1:]) == 8)
	require.True(t, bytes.Equal(key[:1], RecordedBeaconTimestampPrefix))
	bBz := key[1:]
	bFromBz := GetBeaconIDFromBytes(bBz)
	require.True(t, bFromBz == bID)
	require.True(t, bytes.Equal(bBz, GetBeaconIDBytes(bID)))

}

func TestBeaconTimestampKey(t *testing.T) {
	bID := uint64(24)
	tsID := uint64(12345)
	key := BeaconTimestampKey(bID, tsID)
	require.True(t, len(key[1:]) == 16)
	require.True(t, bytes.Equal(key[:1], RecordedBeaconTimestampPrefix))

	bIDbz := key[1:9]
	tsIdBz := key[9:]

	wcIdFromBz := binary.BigEndian.Uint64(bIDbz)
	tsIdFromBz := binary.BigEndian.Uint64(tsIdBz)

	require.True(t, wcIdFromBz == bID)
	require.True(t, tsIdFromBz == tsID)
}

func TestBeaconStorageLimitKey(t *testing.T) {
	bID := uint64(24)
	key := BeaconStorageLimitKey(bID)
	require.True(t, len(key[1:]) == 8)
	require.True(t, bytes.Equal(key[:1], BeaconStorageLimitPrefix))

	bBz := key[1:]
	bFromBz := GetBeaconIDFromBytes(bBz)
	require.True(t, bFromBz == bID)
	require.True(t, bytes.Equal(bBz, GetBeaconIDBytes(bID)))
}
