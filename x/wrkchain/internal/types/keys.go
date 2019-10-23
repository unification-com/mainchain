package types

const (
	// module name
	ModuleName          = "wrkchain"
	BlockStoreKeySuffix = "block"

	// StoreKey to be used when creating the KVStore
	StoreKey      = ModuleName
	StoreKeyBlock = StoreKey + BlockStoreKeySuffix
)
