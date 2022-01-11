package keeper_test

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

const (
	TestDenomination = "nund"
)

var (
	TestAddrs  = createRandomAccounts(10)
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateRandomAccounts(num int) []sdk.AccAddress {
	return createRandomAccounts(num)
}

func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func ParamsEqual(paramsA, paramsB types.Params) bool {
	return paramsA == paramsB
}

func LockedUndEqual(lA, lB types.LockedUnd) bool {
	return lA == lB
}

func RandomDecision() types.PurchaseOrderStatus {
	rnd := rand.Intn(100)
	if rnd >= 50 {
		return types.StatusAccepted
	}
	return types.StatusRejected
}

func RandomStatus() types.PurchaseOrderStatus {
	rnd := test_helpers.RandInBetween(1, 5)
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
