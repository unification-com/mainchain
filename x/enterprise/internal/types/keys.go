package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// module name
	ModuleName = "enterprise"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for the enterprise store.
	QuerierRoute = StoreKey
)

var (

	// key used to store the current highest purchase order ID
	HighestPurchaseOrderIDKey = []byte{0x20}

	// prefix used to store/retrieve an purchase order waiting to be processed from the store
	PurchaseOrderIDKeyPrefix = []byte{0x01}

	// LockedUndAddressKeyPrefix prefix for address keys - used to store locked UND for an account
	LockedUndAddressKeyPrefix = []byte{0x02}

	// WhitelistKeyPrefix is the prifix for whitelisted addresses
	WhitelistKeyPrefix = []byte{0x03}

	TotalLockedUndKey = []byte{0x99}
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

// AddressStoreKey turn an address to key used for enterprise und/locked data to get it from the store
func AddressStoreKey(acc sdk.AccAddress) []byte {
	return append(LockedUndAddressKeyPrefix, acc.Bytes()...)
}

// WhitelistAddressStoreKey turn an address to key used for the whitelist store
func WhitelistAddressStoreKey(acc sdk.AccAddress) []byte {
	return append(WhitelistKeyPrefix, acc.Bytes()...)
}
