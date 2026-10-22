package keeper_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

func RandomDecision() types.PurchaseOrderStatus {
	rnd := rand.Intn(100)
	if rnd >= 50 {
		return types.StatusAccepted
	}
	return types.StatusRejected
}

func RandomStatus() types.PurchaseOrderStatus {
	rnd := simapphelpers.RandInBetween(1, 5)
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

func AddressInDecisions(addr sdk.AccAddress, decisions types.PurchaseOrderDecisions) bool {
	for _, d := range decisions {
		if d.Signer == addr.String() {
			return true
		}
	}
	return false
}
