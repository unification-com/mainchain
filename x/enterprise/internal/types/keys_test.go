package types

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestPurchaseOrderKeys(t *testing.T) {
	// key purchase order
	poID := uint64(24)
	key := PurchaseOrderKey(poID)
	require.True(t, len(key[1:]) == 8)
	require.True(t, bytes.Equal(key[:1], PurchaseOrderIDKeyPrefix))
	poBz := key[1:]
	poFromBz := GetPurchaseOrderIDFromBytes(poBz)
	require.True(t, poFromBz == poID)
	require.True(t, bytes.Equal(poBz, GetPurchaseOrderIDBytes(poID)))
}

func TestAddressKeys(t *testing.T) {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	key := AddressStoreKey(addr)
	require.True(t, len(key[1:]) == sdk.AddrLen)
	require.True(t, bytes.Equal(key[:1], LockedUndAddressKeyPrefix))
	addrBz := key[1:]
	addrFromBytes := sdk.AccAddress(addrBz)
	require.True(t, addrFromBytes.String() == addr.String())
}
