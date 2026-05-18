package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// BEACON fees, in nano FUND
	RegFee             = 1000000000000 // 1000 FUND - used in init genesis
	RecordFee          = 1000000000    // 1 FUND - used in init genesis
	PurchaseStorageFee = 5000000000    // 5 FUND - used in init genesis

	DefaultStartingBeaconID uint64 = 1 // used in init genesis

	DefaultStorageLimit    uint64 = 50000  // used in init genesis
	DefaultMaxStorageLimit uint64 = 600000 // used in init genesis
)

var (
	FeeDenom = sdk.DefaultBondDenom // used in init genesis
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
// expected state holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	seenBeaconIDs := make(map[uint64]struct{}, len(data.RegisteredBeacons))
	for i, record := range data.RegisteredBeacons {
		if record.Beacon.BeaconId == 0 {
			return fmt.Errorf("beacon[%d]: invalid BeaconId 0", i)
		}
		if _, dup := seenBeaconIDs[record.Beacon.BeaconId]; dup {
			return fmt.Errorf("beacon[%d]: duplicate BeaconId %d", i, record.Beacon.BeaconId)
		}
		seenBeaconIDs[record.Beacon.BeaconId] = struct{}{}
		if record.Beacon.BeaconId >= data.StartingBeaconId {
			return fmt.Errorf("beacon[%d]: BeaconId %d must be < StartingBeaconId %d",
				i, record.Beacon.BeaconId, data.StartingBeaconId)
		}
		if _, err := sdk.AccAddressFromBech32(record.Beacon.Owner); err != nil {
			return fmt.Errorf("beacon[%d]: invalid Owner %q: %w", i, record.Beacon.Owner, err)
		}
		if record.Beacon.Moniker == "" {
			return fmt.Errorf("beacon[%d]: empty Moniker", i)
		}
		if record.InStateLimit == 0 {
			return fmt.Errorf("beacon[%d]: zero InStateLimit", i)
		}
		seenTSIDs := make(map[uint64]struct{}, len(record.Timestamps))
		for j, timestamp := range record.Timestamps {
			if timestamp.Id == 0 {
				return fmt.Errorf("beacon[%d].timestamp[%d]: invalid TimestampID 0", i, j)
			}
			if _, dup := seenTSIDs[timestamp.Id]; dup {
				return fmt.Errorf("beacon[%d].timestamp[%d]: duplicate TimestampID %d", i, j, timestamp.Id)
			}
			seenTSIDs[timestamp.Id] = struct{}{}
			if timestamp.H == "" {
				return fmt.Errorf("beacon[%d].timestamp[%d]: empty Hash", i, j)
			}
			if timestamp.T == 0 {
				return fmt.Errorf("beacon[%d].timestamp[%d]: zero SubmitTime", i, j)
			}
		}
	}
	return nil
}
