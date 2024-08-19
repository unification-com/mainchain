package simulation_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/enterprise/simulation"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	testApp := simapp.Setup(t, false)
	cdc := testApp.AppCodec()
	dec := simulation.NewDecodeStore(cdc)

	denom := simapp.TestDenomination

	purchaseOrder, err := types.NewEnterpriseUndPurchaseOrder(1, delAddr1.String(), sdk.NewInt64Coin(denom, 100000000),
		types.StatusRaised, 1234, 5678)
	require.NoError(t, err)

	purchaseOrderBz, err := cdc.Marshal(&purchaseOrder)
	require.NoError(t, err)

	lockedUnd, err := types.NewLockedUnd(delAddr1.String(), sdk.NewInt64Coin(denom, 100000000))
	require.NoError(t, err)

	lockedUndBz, err := cdc.Marshal(&lockedUnd)
	require.NoError(t, err)

	totalLocked := sdk.NewInt64Coin(denom, 100000000)
	totalLockedBz, err := cdc.Marshal(&totalLocked)
	require.NoError(t, err)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.PurchaseOrderKey(1), Value: purchaseOrderBz},
			{Key: types.LockedUndAddressStoreKey(delAddr1), Value: lockedUndBz},
			{Key: types.TotalLockedUndKey, Value: totalLockedBz},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"purchase orders", fmt.Sprintf("%v\n%v", purchaseOrder, purchaseOrder)},
		{"locked unds", fmt.Sprintf("%v\n%v", lockedUnd, lockedUnd)},
		{"total locked", fmt.Sprintf("%v\n%v", totalLocked, totalLocked)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
