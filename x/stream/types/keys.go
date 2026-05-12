package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
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

	// StreamKeyPrefix prefix for the Stream store
	StreamKeyPrefix = []byte{0x11}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetStreamKey creates the key for the (receiver, sender, denom) triple
// 0x11 | len(receiver) | receiver | len(sender) | sender | len(denom) | denom
func GetStreamKey(receiverAddr sdk.AccAddress, senderAddr sdk.AccAddress, denom string) []byte {
	return append(GetStreamsByPairKey(receiverAddr, senderAddr), lengthPrefixedDenom(denom)...)
}

// GetStreamsByReceiverKey is the prefix for "all streams to receiver"
// 0x11 | len(receiver) | receiver
func GetStreamsByReceiverKey(receiverAddr sdk.AccAddress) []byte {
	return append(StreamKeyPrefix, address.MustLengthPrefix(receiverAddr)...)
}

// GetStreamsByPairKey is the prefix for "all streams between this (receiver, sender) pair, one entry per denom"
// 0x11 | len(receiver) | receiver | len(sender) | sender
func GetStreamsByPairKey(receiverAddr sdk.AccAddress, senderAddr sdk.AccAddress) []byte {
	return append(GetStreamsByReceiverKey(receiverAddr), address.MustLengthPrefix(senderAddr)...)
}

// AddressesFromStreamKey returns (receiver, sender, denom) from a full stream store key.
func AddressesFromStreamKey(key []byte) (sdk.AccAddress, sdk.AccAddress, string) {
	// key is of format:
	// 0x11<receiverAddrLen (1 Byte)><receiverAddr><senderAddrLen (1 Byte)><senderAddr><denomLen (1 Byte)><denom>

	receiverAddrLen, receiverAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, 1, 1) // ignore key[0] since it is a prefix key
	receiverAddr, receiverAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrLenEndIndex+1, int(receiverAddrLen[0]))

	senderAddrLen, senderAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrEndIndex+1, 1)
	senderAddr, senderAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, senderAddrLenEndIndex+1, int(senderAddrLen[0]))

	denomLen, denomLenEndIndex := sdk.ParseLengthPrefixedBytes(key, senderAddrEndIndex+1, 1)
	denomBytes, denomEndIndex := sdk.ParseLengthPrefixedBytes(key, denomLenEndIndex+1, int(denomLen[0]))

	kv.AssertKeyAtLeastLength(key, denomEndIndex+1)
	return receiverAddr, senderAddr, string(denomBytes)
}

// FirstAddressFromStreamStoreKey parses the first address only
func FirstAddressFromStreamStoreKey(key []byte) sdk.AccAddress {
	addrLen := key[0]
	return sdk.AccAddress(key[1 : 1+addrLen])
}

// lengthPrefixedDenom returns a single-byte length prefix followed by the denom bytes.
// The denom length is bounded by sdk.Coin validation (≤128 chars) so a uint8 prefix is sufficient.
func lengthPrefixedDenom(denom string) []byte {
	b := []byte(denom)
	if len(b) > 255 {
		panic("denom length must be <= 255")
	}
	return append([]byte{byte(len(b))}, b...)
}
