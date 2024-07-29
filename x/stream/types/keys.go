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

// GetStreamKey creates the key for receiver bond with sender
func GetStreamKey(receiverAddr sdk.AccAddress, senderAddr sdk.AccAddress) []byte {
	return append(GetStreamsByReceiverKey(receiverAddr), address.MustLengthPrefix(senderAddr)...)
}

// GetStreamsByReceiverKey creates the prefix for a delegator for all validators
func GetStreamsByReceiverKey(receiverAddr sdk.AccAddress) []byte {
	return append(StreamKeyPrefix, address.MustLengthPrefix(receiverAddr)...)
}

// AddressesFromStreamKey returns a receiver and sender address from a stream prefix
// store key.
func AddressesFromStreamKey(key []byte) (sdk.AccAddress, sdk.AccAddress) {
	// key is of format:
	// 0x11<receiverAddressLen (1 Byte)><receiverAddress_Bytes><senderAddressLen (1 Byte)><senderAddress_Bytes>

	receiverAddrLen, receiverAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, 1, 1) // ignore key[0] since it is a prefix key
	receiverAddr, receiverAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrLenEndIndex+1, int(receiverAddrLen[0]))

	senderAddrLen, senderAddrLenEndIndex := sdk.ParseLengthPrefixedBytes(key, receiverAddrEndIndex+1, 1)
	senderAddr, senderAddrEndIndex := sdk.ParseLengthPrefixedBytes(key, senderAddrLenEndIndex+1, int(senderAddrLen[0]))

	kv.AssertKeyAtLeastLength(key, senderAddrEndIndex+1)
	return receiverAddr, senderAddr
}

// FirstAddressFromStreamStoreKey parses the first address only
func FirstAddressFromStreamStoreKey(key []byte) sdk.AccAddress {
	addrLen := key[0]
	return sdk.AccAddress(key[1 : 1+addrLen])
}
