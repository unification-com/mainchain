package types

import (
	"bytes"
	"fmt"
)

// GenesisState - beacon state
type GenesisState struct {
	Params           Params         `json:"params" yaml:"params"`                         // beacon params
	StartingBeaconID uint64         `json:"starting_beacon_id" yaml:"starting_beacon_id"` // should be 1
	Beacons          []BeaconExport `json:"registered_beacons" yaml:"registered_beacons"`
}

type BeaconExport struct {
	Beacon           Beacon
	BeaconTimestamps []BeaconTimestamp
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingBeaconID uint64) GenesisState {
	return GenesisState{
		Params:           params,
		StartingBeaconID: startingBeaconID,
		Beacons:          nil,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:           DefaultParams(),
		StartingBeaconID: DefaultStartingBeaconID,
	}
}

// Equal checks whether two beacon GenesisState structs are equivalent
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(data)
	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty
func (data GenesisState) IsEmpty() bool {
	return data.Equal(GenesisState{})
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.Beacons {
		if record.Beacon.Owner == nil {
			return fmt.Errorf("invalid Beacon: Owner: %s. Error: Missing Owner", record.Beacon.Owner)
		}
		if record.Beacon.BeaconID == 0 {
			return fmt.Errorf("invalid Beacon: Moniker: %d. Error: Missing ID", record.Beacon.BeaconID)
		}
		if record.Beacon.Name == "" {
			return fmt.Errorf("invalid Beacon: Name: %s. Error: Missing name", record.Beacon.Name)
		}
		for _, timestamp := range record.BeaconTimestamps {
			if timestamp.Owner == nil {
				return fmt.Errorf("invalid Beacon timestamp: Owner: %s. Error: Missing Owner", timestamp.Owner)
			}
			if timestamp.Hash == "" {
				return fmt.Errorf("invalid Beacon timestamp: Hash: %s. Error: Missing Hash", timestamp.Hash)
			}
			if timestamp.SubmitTime == 0 {
				return fmt.Errorf("invalid Beacon timestamp: SubmitTime: %d. Error: Missing SubmitTime", timestamp.SubmitTime)
			}
			if timestamp.TimestampID == 0 {
				return fmt.Errorf("invalid Beacon timestamp: TimestampID: %d. Error: Missing TimestampID", timestamp.TimestampID)
			}
			if timestamp.BeaconID == 0 {
				return fmt.Errorf("invalid Beacon timestamp: BeaconID: %d. Error: Missing BeaconID", timestamp.BeaconID)
			}
		}
	}
	return nil
}
