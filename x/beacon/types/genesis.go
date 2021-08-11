package types

import (
	"fmt"
	undtypes "github.com/unification-com/mainchain/types"
)

const (
	// BEACON fees, in nano FUND
	RegFee    = 1000000000000                // 1000 FUND - used in init genesis
	RecordFee = 1000000000                   // 1 FUND - used in init genesis
	FeeDenom  = undtypes.DefaultDenomination // used in init genesis

	DefaultStartingBeaconID uint64 = 1 // used in init genesis
)

//// GenesisState - beacon state
//type GenesisState struct {
//	Params           Params         `json:"params" yaml:"params"`                         // beacon params
//	StartingBeaconID uint64         `json:"starting_beacon_id" yaml:"starting_beacon_id"` // should be 1
//	Beacons          []BeaconExport `json:"registered_beacons" yaml:"registered_beacons"`
//}
//
//type BeaconExport struct {
//	Beacon           Beacon                         `json:"beacon" yaml:"beacon"`
//	BeaconTimestamps []BeaconTimestampGenesisExport `json:"timestamps" yaml:"timestamps"`
//}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingBeaconID uint64) *GenesisState {
	return &GenesisState{
		Params:            params,
		StartingBeaconId:  startingBeaconID,
		RegisteredBeacons: nil,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:           DefaultParams(),
		StartingBeaconId: DefaultStartingBeaconID,
		RegisteredBeacons: nil,
	}
}

//// Equal checks whether two beacon GenesisState structs are equivalent
//func (data GenesisState) Equal(data2 GenesisState) bool {
//	b1 := ModuleCdc.MustMarshalBinaryBare(data)
//	b2 := ModuleCdc.MustMarshalBinaryBare(data2)
//	return bytes.Equal(b1, b2)
//}
//
//// IsEmpty returns true if a GenesisState is empty
//func (data GenesisState) IsEmpty() bool {
//	return data.Equal(GenesisState{})
//}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	for _, record := range data.RegisteredBeacons {
		if record.Beacon.BeaconId == 0 {
			return fmt.Errorf("invalid Beacon: ID: %d. Error: Missing ID", record.Beacon.BeaconId)
		}
		if record.Beacon.Owner == "" {
			return fmt.Errorf("invalid Beacon: Owner: %s. Error: Missing Owner", record.Beacon.Owner)
		}
		if record.Beacon.Moniker == "" {
			return fmt.Errorf("invalid Beacon: Moniker: %s. Error: Missing Moniker", record.Beacon.Moniker)
		}
		for _, timestamp := range record.Timestamps {
			if timestamp.Id == 0 {
				return fmt.Errorf("invalid Beacon timestamp: TimestampID: %d. Error: Missing TimestampID", timestamp.Id)
			}
			if timestamp.H == "" {
				return fmt.Errorf("invalid Beacon timestamp: Hash: %s. Error: Missing Hash", timestamp.H)
			}
			if timestamp.T == 0 {
				return fmt.Errorf("invalid Beacon timestamp: SubmitTime: %d. Error: Missing SubmitTime", timestamp.T)
			}
		}
	}
	return nil
}
