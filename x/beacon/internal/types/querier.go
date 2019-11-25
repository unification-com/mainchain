package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryResBeaconTimestampHashes Queries Result Payload for a Beacon timestamp Hashes query
type QueryResBeaconTimestampHashes []BeaconTimestamp

// implement fmt.Stringer
func (h QueryResBeaconTimestampHashes) String() (out string) {
	for _, val := range h {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// QueryResBeacons Queries BEACONs
type QueryResBeacons []Beacon

// implement fmt.Stringer
func (wc QueryResBeacons) String() (out string) {
	for _, val := range wc {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// QueryBeaconParams Params for query 'custom/beacon/registered'
type QueryBeaconParams struct {
	Page    int
	Limit   int
	Moniker string
	Owner   sdk.AccAddress
}

// NewQueryBeaconParams creates a new instance of QueryBeaconParams
func NewQueryBeaconParams(page, limit int, monkker string, owner sdk.AccAddress) QueryBeaconParams {
	return QueryBeaconParams{
		Page:    page,
		Limit:   limit,
		Moniker: monkker,
		Owner:   owner,
	}
}
