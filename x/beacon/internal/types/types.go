package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain-cosmos/types"
)

const (
	// BEACON fees, in nano UND
	RegFee    = 1000000000000                // 1000 UND - used in init genesis
	RecordFee = 1000000000                   // 1 UND - used in init genesis
	FeeDenom  = undtypes.DefaultDenomination // used in init genesis

	DefaultStartingBeaconID uint64 = 1 // used in init genesis
)


// Beacons is an array of Beacon
type Beacons []Beacon

// String implements stringer interface
func (w Beacons) String() string {
	out := "ID - [Moniker] 'Name' {LastTimestampId} Owner\n"
	for _, b := range w {
		out += fmt.Sprintf("%d - [%s] '%s' {%d} %s\n",
			b.BeaconID, b.Moniker,
			b.Name, b.LastTimestampID, b.Owner)
	}
	return strings.TrimSpace(out)
}

// Beacon is a struct that contains all the metadata of a registered BEACON
type Beacon struct {
	BeaconID        uint64         `json:"beacon_id"`
	Moniker         string         `json:"moniker"`
	Name            string         `json:"name"`
	LastTimestampID uint64         `json:"last_timestamp_id"`
	Owner           sdk.AccAddress `json:"owner"`
}

// NewBeacon returns a new Beacon struct
func NewBeacon() Beacon {
	return Beacon{}
}

// implement fmt.Stringer
func (b Beacon) String() string {
	return strings.TrimSpace(fmt.Sprintf(`BeaconID: %d
Moniker: %s
Name: %s
LastTimestampID: %d
Owner: %s`, b.BeaconID, b.Moniker, b.Name, b.LastTimestampID, b.Owner))
}

// BeaconTimestamp is a struct that contains a BEACON's recorded timestamp hash
type BeaconTimestamp struct {
	BeaconID     uint64         `json:"beacon_id"`
	TimestampID  uint64         `json:"timestamp_id"`
	SubmitTime   int64          `json:"submit_time"`
	Hash         string         `json:"hash"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewBeaconTimestamp returns a new BeaconTimestamp struct
func NewBeaconTimestamp() BeaconTimestamp {
	return BeaconTimestamp{}
}

// implement fmt.Stringer
func (bts BeaconTimestamp) String() string {
	return strings.TrimSpace(fmt.Sprintf(`BeaconID: %d
TimestampID: %d
SubmitTime: %d
Hash: %s
Owner: %s`, bts.BeaconID, bts.TimestampID, bts.SubmitTime, bts.Hash, bts.Owner))
}
