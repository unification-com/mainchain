package types

import (
	"crypto/sha256"
	"strconv"
)

const (
	// module name
	ModuleName = "wrkchain"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// WRKChain Recorded block hash delimeter
	Delimeter = "_"
)

var (
	// RegisteredWrkChainPrefix prefix for registered WRKChain store
	RegisteredWrkChainPrefix = []byte{0x01}

	// RegisteredWrkChainPrefix prefix for WRKChain Hashes store
	RecordedWrkChainBlockHashPrefix = []byte{0x02}
)

// GetWrkChainStoreKey turn an address to key used to get it from the account store
func GetWrkChainStoreKey(wrkchainId string) []byte {
	return append(RegisteredWrkChainPrefix, []byte(wrkchainId)...)
}

func GetWrkChainBlockHashStoreKey(wrkchainId string, height uint64) []byte {
	heightString := strconv.FormatUint(height, 10)
	wrkchainIdPrefix := GetWrkChainBlockHashStoreKeyPrefix(wrkchainId)
	heightSuffix := append([]byte(Delimeter), []byte(heightString)...)
	return append(wrkchainIdPrefix, heightSuffix...)
}

// GetWrkChainBlockHashStoreKeyPrefix returns a key for a WRKChain which can be used
// for iterating through the recorded hashes
func GetWrkChainBlockHashStoreKeyPrefix(wrkchainId string) []byte {
	h := sha256.New()

	// make WRKChainID into a hash to try and avoid potential clashes with ID + height concatenation
	// todo - get and handle err
	_, _ = h.Write([]byte(wrkchainId))

	return append(RecordedWrkChainBlockHashPrefix, h.Sum(nil)...)
}
