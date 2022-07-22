package types

import (
	"fmt"
	undtypes "github.com/unification-com/mainchain/types"
)

const (
	// BEACON fees, in nano FUND
	RegFee             = 1000000000000                // 1000 FUND - used in init genesis
	RecordFee          = 1000000000                   // 1 FUND - used in init genesis
	PurchaseStorageFee = 5000000000                   // 5 FUND - used in init genesis
	FeeDenom           = undtypes.DefaultDenomination // used in init genesis

	DefaultStartingBeaconID uint64 = 1 // used in init genesis

	DefaultStorageLimit    uint64 = 50000  // used in init genesis
	DefaultMaxStorageLimit uint64 = 600000 // used in init genesis
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, startingBeaconID uint64, beacons BeaconExports) *GenesisState {
	return &GenesisState{
		Params:            params,
		StartingBeaconId:  startingBeaconID,
		RegisteredBeacons: beacons,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:            DefaultParams(),
		StartingBeaconId:  DefaultStartingBeaconID,
		RegisteredBeacons: nil,
	}
}

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
		if record.InStateLimit == 0 {
			return fmt.Errorf("invalid Beacon: InStateLimit: %d. Error: Missing InStateLimit", record.InStateLimit)
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
