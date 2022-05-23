package types

const MaxHashSubmissionsToExport = 20000

type BeaconExports []BeaconExport
type BeaconTimestampGenesisExports []BeaconTimestampGenesisExport

// BeaconTimestampLegacy is onlyused to support old style timestamp output for the legact REST endpoint
type BeaconTimestampLegacy struct {
	BeaconID    uint64 `json:"beacon_id"`
	TimestampID uint64 `json:"timestamp_id"`
	SubmitTime  uint64 `json:"submit_time"`
	Hash        string `json:"hash"`
	Owner       string `json:"owner"`
}

func NewBeacon(beaconId uint64, moniker string, name string, lastTimeStampId uint64, owner string) (Beacon, error) {
	b := Beacon{
		BeaconId:        beaconId,
		Moniker:         moniker,
		Name:            name,
		LastTimestampId: lastTimeStampId,
		Owner:           owner,
	}

	return b, nil
}

func NewBeaconTimestamp(timestampId uint64, submitTime uint64, hash string) (BeaconTimestamp, error) {
	bTs := BeaconTimestamp{
		TimestampId: timestampId,
		SubmitTime:  submitTime,
		Hash:        hash,
	}

	return bTs, nil
}
