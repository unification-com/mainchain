package types

const (
	// module name
	ModuleName = "beacon"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for the enterprise store.
	QuerierRoute = StoreKey
)
