package types

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "enterprise"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for the enterprise store.
	QuerierRoute = StoreKey
)

var (

	// key used to store the current highest purchase order ID
	HighestPurchaseOrderIDKey = []byte{0x20}

	// prefix used to store/retrieve an purchase order waiting to be processed from the store
	PurchaseOrderIDKeyPrefix = []byte{0x01}

	// LockedUndAddressKeyPrefix prefix for address keys - used to store locked eFUND for an account
	LockedUndAddressKeyPrefix = []byte{0x02}

	// WhitelistKeyPrefix is the prefix for whitelisted addresses
	WhitelistKeyPrefix = []byte{0x03}

	// RaisedPoPrefix used to temporarily store currently raised purchase orders for the ABCI blocker
	RaisedPoPrefix = []byte{0x04}

	// AcceptedPoPrefix used to temporarily store currently accepted purchase orders for the ABCI blocker
	AcceptedPoPrefix = []byte{0x05}

	// SpentEFUNDAddressKeyPrefix prefix for address keys - used to store a tally of used eFUND for an account
	SpentEFUNDAddressKeyPrefix = []byte{0x06}

	ParamsKey = []byte{0x07}

	TotalSpentEFUNDKey = []byte{0x98}
	TotalLockedUndKey  = []byte{0x99}
)

// GetPurchaseOrderIDBytes returns the byte representation of the purchaseOrderID
// used for getting the highest purchase order ID from the database
func GetPurchaseOrderIDBytes(purchaseOrderID uint64) (purchaseOrderIDBz []byte) {
	purchaseOrderIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(purchaseOrderIDBz, purchaseOrderID)
	return
}

// GetPurchaseOrderIDFromBytes returns purchaseOrderID in uint64 format from a byte array
// used for getting the highest purchase order ID from the database
func GetPurchaseOrderIDFromBytes(bz []byte) (purchaseOrderID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// PurchaseOrderKey gets a specific purchase order ID key for use in the store
func PurchaseOrderKey(purchaseOrderID uint64) []byte {
	return append(PurchaseOrderIDKeyPrefix, GetPurchaseOrderIDBytes(purchaseOrderID)...)
}

// LockedUndAddressStoreKey turn an address to key used for enterprise und/locked data to get it from the store
func LockedUndAddressStoreKey(acc sdk.AccAddress) []byte {
	return append(LockedUndAddressKeyPrefix, acc.Bytes()...)
}

// SpentEFUNDAddressStoreKey turn an address to key used for spent eFUND data to get it from the store
func SpentEFUNDAddressStoreKey(acc sdk.AccAddress) []byte {
	return append(SpentEFUNDAddressKeyPrefix, acc.Bytes()...)
}

// WhitelistAddressStoreKey turn an address to key used for the whitelist store
func WhitelistAddressStoreKey(acc sdk.AccAddress) []byte {
	return append(WhitelistKeyPrefix, acc.Bytes()...)
}

// RaisedQueueStoreKey us used to temporarily store raised PO order IDs for the ABCI blocker
func RaisedQueueStoreKey(purchaseOrderID uint64) []byte {
	return append(RaisedPoPrefix, GetPurchaseOrderIDBytes(purchaseOrderID)...)
}

// SplitRaisedQueueKey is used to get the PO ID from the storekey
func SplitRaisedQueueKey(key []byte) (purchaseOrderId uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d != 8)", len(key[1:])))
	}
	return GetPurchaseOrderIDFromBytes(key[1:])
}

// AcceptedQueueStoreKey us used to temporarily store accepted PO order IDs for the ABCI blocker
func AcceptedQueueStoreKey(purchaseOrderID uint64) []byte {
	return append(AcceptedPoPrefix, GetPurchaseOrderIDBytes(purchaseOrderID)...)
}

// SplitAcceptedQueueKey is used to get the PO ID from the storekey
func SplitAcceptedQueueKey(key []byte) (purchaseOrderId uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d != 8)", len(key[1:])))
	}
	return GetPurchaseOrderIDFromBytes(key[1:])
}
