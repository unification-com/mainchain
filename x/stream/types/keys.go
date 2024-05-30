package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "stream"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stream"
)

var (
	// ParamsKey is the prefix for the params store
	ParamsKey = []byte{0x01}

	// HighestStreamIdKey key used to store the current highest Stream ID
	HighestStreamIdKey = []byte{0x02}

	// StreamKeyPrefix prefix for the Stream store
	StreamKeyPrefix = []byte{0x11}

	// StreamIdLookupKeyPrefix prefix for the StreamIdLookup store
	StreamIdLookupKeyPrefix = []byte{0x21}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetStreamKey creates the key for receiver bond with sender
func GetStreamKey(receiverAddr sdk.AccAddress, senderAddr sdk.AccAddress) []byte {
	return append(GetStreamsByReceiverKey(receiverAddr), address.MustLengthPrefix(senderAddr)...)
}

// GetStreamsByReceiverKey creates the prefix for a delegator for all validators
func GetStreamsByReceiverKey(receiverAddr sdk.AccAddress) []byte {
	return append(StreamKeyPrefix, address.MustLengthPrefix(receiverAddr)...)
}

// GetStreamIdBytes returns the byte representation of the streamId
// used for getting the highest Stream ID from the database
func GetStreamIdBytes(streamId uint64) (streamIdBz []byte) {
	streamIdBz = make([]byte, 8)
	binary.BigEndian.PutUint64(streamIdBz, streamId)
	return
}

// GetStreamIdFromBytes returns BeaconID in uint64 format from a byte array
// used for getting the highest Beacon ID from the database
func GetStreamIdFromBytes(bz []byte) (streamId uint64) {
	return binary.BigEndian.Uint64(bz)
}

// GetStreamIdLookupKey gets a specific stream lookup pair from the given stream ID
func GetStreamIdLookupKey(streamId uint64) []byte {
	return append(StreamIdLookupKeyPrefix, GetStreamIdBytes(streamId)...)
}
