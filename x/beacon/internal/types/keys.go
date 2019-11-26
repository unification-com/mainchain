package types

import "encoding/binary"

const (
	// module name
	ModuleName = "beacon"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	DefaultParamspace = ModuleName

	// QuerierRoute is the querier route for the BEACON store.
	QuerierRoute = StoreKey
)

var (

	// key used to store the current highest BEACON ID
	HighestBeaconIDKey = []byte{0x20}

	// RegisteredBeaconPrefix prefix for registered BEACON store
	RegisteredBeaconPrefix = []byte{0x01}

	// RecordedBeaconTimestampPrefix prefix for BEACON Timestamps store
	RecordedBeaconTimestampPrefix = []byte{0x02}
)

// GetBeaconIDBytes returns the byte representation of the BeaconID
// used for getting the highest Beacon ID from the database
func GetBeaconIDBytes(beaconID uint64) (beaconIDBz []byte) {
	beaconIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(beaconIDBz, beaconID)
	return
}

// GetBeaconIDFromBytes returns BeaconID in uint64 format from a byte array
// used for getting the highest Beacon ID from the database
func GetBeaconIDFromBytes(bz []byte) (beaconID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// BeaconKey gets a specific purchase order ID key for use in the store
func BeaconKey(beaconID uint64) []byte {
	return append(RegisteredBeaconPrefix, GetBeaconIDBytes(beaconID)...)
}

// BeaconAllTimestampsKey gets the key for a specific BEACON's timestamps
func BeaconAllTimestampsKey(beaconID uint64) []byte {
	return append(RecordedBeaconTimestampPrefix, GetBeaconIDBytes(beaconID)...)
}

// BeaconTimestampKey gets the key for a single BEACON's specific timestamp ID
func BeaconTimestampKey(beaconID, timestampID uint64) []byte {
	blocksKey := BeaconAllTimestampsKey(beaconID)
	timestampIdBz := GetTimestampIDBytes(timestampID)
	return append(blocksKey, timestampIdBz...)
}

func GetTimestampIDBytes(timestampID uint64) (timestampIDBz []byte) {
	timestampIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(timestampIDBz, timestampID)
	return
}

func GetTimestampIDFromBytes(bz []byte) (timestampID uint64) {
	return binary.BigEndian.Uint64(bz)
}
