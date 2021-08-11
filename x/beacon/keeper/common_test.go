package keeper_test

import (
	"github.com/unification-com/mainchain/x/beacon/types"
)

// BeaconEqual checks if two Beacons are equal
func BeaconEqual(wcA types.Beacon, wcB types.Beacon) bool {
	return wcA == wcB
}
//
//// ParamsEqual checks params are equal
//func ParamsEqual(paramsA, paramsB types.Params) bool {
//	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(paramsA),
//		types.ModuleCdc.MustMarshalBinaryBare(paramsB))
//}
//
// BeaconTimestampEqual checks if two BeaconTimestamps are equal
func BeaconTimestampEqual(lA, lB types.BeaconTimestamp) bool {
	return lA == lB
}