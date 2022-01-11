package types

const MaxHashSubmissionsToExport = 20000

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

func NewBeaconTimestamp(beaconId uint64, timestampId uint64, submitTime uint64, hash string, owner string) (BeaconTimestamp, error) {
	bTs := BeaconTimestamp{
		BeaconId:    beaconId,
		TimestampId: timestampId,
		SubmitTime:  submitTime,
		Hash:        hash,
		Owner:       owner,
	}

	return bTs, nil
}
