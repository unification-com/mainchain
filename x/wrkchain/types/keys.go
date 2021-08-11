package types

import (
	"encoding/binary"
)

const (
	// module name
	ModuleName = "wrkchain"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for the wrkchain store.
	QuerierRoute = StoreKey
)

var (

	// key used to store the current highest WRKChain ID
	HighestWrkChainIDKey = []byte{0x20}

	// RegisteredWrkChainPrefix prefix for registered WRKChain store
	RegisteredWrkChainPrefix = []byte{0x01}

	// RegisteredWrkChainPrefix prefix for WRKChain Hashes store
	RecordedWrkChainBlockHashPrefix = []byte{0x02}
)

// GetWrkChainIDBytes returns the byte representation of the wrkChainID
// used for getting the highest WRKChain ID from the database
func GetWrkChainIDBytes(wrkChainID uint64) (wrkChainIDBz []byte) {
	wrkChainIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(wrkChainIDBz, wrkChainID)
	return
}

// GetWrkChainIDFromBytes returns wrkChainID in uint64 format from a byte array
// used for getting the highest WRKChain ID from the database
func GetWrkChainIDFromBytes(bz []byte) (wrkChainID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// WrkChainKey gets a specific purchase order ID key for use in the store
func WrkChainKey(wrkChainID uint64) []byte {
	return append(RegisteredWrkChainPrefix, GetWrkChainIDBytes(wrkChainID)...)
}

func WrkChainAllBlocksKey(wrkChainID uint64) []byte {
	return append(RecordedWrkChainBlockHashPrefix, GetWrkChainIDBytes(wrkChainID)...)
}

func WrkChainBlockKey(wrkChainID, height uint64) []byte {
	blocksKey := WrkChainAllBlocksKey(wrkChainID)
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, height)
	return append(blocksKey, heightBz...)
}
