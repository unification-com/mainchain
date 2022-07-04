package v040

import "encoding/binary"

var (

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
