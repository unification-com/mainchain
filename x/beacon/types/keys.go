package types

import (
	"crypto/sha256"
	"encoding/binary"
)

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

	// BeaconStorageLimitPrefix prefix for BEACON storage limit store
	BeaconStorageLimitPrefix = []byte{0x03}

	ParamsKey = []byte{0x04}

	// RecordedBeaconHashIndexPrefix is the prefix for the (beacon, hash) -> timestampIds index (#129).
	// It is one-to-many: the same hash can be recorded many times (the dpv heartbeat re-stamps the same
	// root), so a (beacon, hash) maps to all the timestamp ids that recorded it.
	RecordedBeaconHashIndexPrefix = []byte{0x05}
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

// BeaconHashIndexHashKey is the prefix for every timestamp of a beacon that recorded a given hash:
// 0x05 | beaconID(8) | sha256(hash)(32). The (arbitrary-length, arbitrary-content) hash is sha256'd to a
// fixed 32 bytes so the key stays fixed-width + range-iterable and sidesteps delimiter issues.
func BeaconHashIndexHashKey(beaconID uint64, hash string) []byte {
	h := sha256.Sum256([]byte(hash))
	key := append(RecordedBeaconHashIndexPrefix, GetBeaconIDBytes(beaconID)...)
	return append(key, h[:]...)
}

// BeaconHashIndexKey is the full one-to-many index entry key: BeaconHashIndexHashKey | timestampID(8).
func BeaconHashIndexKey(beaconID uint64, hash string, timestampID uint64) []byte {
	return append(BeaconHashIndexHashKey(beaconID, hash), GetTimestampIDBytes(timestampID)...)
}

// BeaconStorageLimitKey gets the key for a single BEACON's specific storage limit
func BeaconStorageLimitKey(beaconID uint64) []byte {
	return append(BeaconStorageLimitPrefix, GetBeaconIDBytes(beaconID)...)
}
