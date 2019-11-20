package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	purchaseOrder := types.NewEnterpriseUndPurchaseOrder()
	purchaseOrder.Purchaser = delAddr1
	purchaseOrder.Status = types.StatusRaised
	purchaseOrder.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 100000000)
	purchaseOrder.PurchaseOrderID = 1

	lockedUnd := types.NewLockedUnd(delAddr1, types.DefaultDenomination)
	lockedUnd.Amount = sdk.NewInt64Coin(types.DefaultDenomination, 100000000)

	totalLocked := sdk.NewInt64Coin(types.DefaultDenomination, 100000000)

	kvPairs := cmn.KVPairs{
		cmn.KVPair{Key: types.PurchaseOrderKey(1), Value: cdc.MustMarshalBinaryLengthPrefixed(purchaseOrder)},
		cmn.KVPair{Key: types.AddressStoreKey(delAddr1), Value: cdc.MustMarshalBinaryLengthPrefixed(lockedUnd)},
		cmn.KVPair{Key: types.TotalLockedUndKey, Value: cdc.MustMarshalBinaryLengthPrefixed(totalLocked)},
		cmn.KVPair{Key: []byte{0x99}, Value: []byte{0x99}},
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
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
