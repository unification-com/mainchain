package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"math/rand"
)

//// PurchaseOrderEqual checks if two purchase orders are equal
//func PurchaseOrderEqual(poA types.EnterpriseUndPurchaseOrder, poB types.EnterpriseUndPurchaseOrder) bool {
//	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(poA),
//		types.ModuleCdc.MustMarshalBinaryBare(poB))
//}
//
//func ParamsEqual(paramsA, paramsB types.Params) bool {
//	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(paramsA),
//		types.ModuleCdc.MustMarshalBinaryBare(paramsB))
//}
//
//func LockedUndEqual(lA, lB types.LockedUnd) bool {
//	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(lA),
//		types.ModuleCdc.MustMarshalBinaryBare(lB))
//}
//
func RandomDecision() types.PurchaseOrderStatus {
	rnd := rand.Intn(100)
	if rnd >= 50 {
		return types.StatusAccepted
	}
	return types.StatusRejected
}

func RandomStatus() types.PurchaseOrderStatus {
	rnd := test_helpers.RandInBetween(1, 4)
	switch rnd {
	case 1:
		return types.StatusRaised
	case 2:
		return types.StatusAccepted
	case 3:
		return types.StatusRejected
	case 4:
		return types.StatusCompleted
	default:
		return types.StatusRaised
	}
}

func AddressInDecisions(addr sdk.AccAddress, decisions []*types.PurchaseOrderDecision) bool {
	for _, d := range decisions {
		if d.Signer == addr.String() {
			return true
		}
	}
	return false
}
