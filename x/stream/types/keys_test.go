package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
)

func TestAddressFromStreamsStore(t *testing.T) {
	receiverAddr, err := sdk.AccAddressFromBech32("cosmos1n88uc38xhjgxzw9nwre4ep2c8ga4fjxcar6mn7")
	require.NoError(t, err)
	receiverAddrLen := len(receiverAddr)
	require.Equal(t, 20, receiverAddrLen)

	senderAddr, err := sdk.AccAddressFromBech32("cosmos139f7kncmglres2nf3h4hc4tade85ekfr8sulz5")
	require.NoError(t, err)
	senderAddrLen := len(senderAddr)
	require.Equal(t, 20, senderAddrLen)

	key := types.GetStreamKey(receiverAddr, senderAddr)

	require.Len(t, key, len(types.StreamKeyPrefix)+len(address.MustLengthPrefix(receiverAddr))+len(address.MustLengthPrefix(senderAddr)))

	r, s := types.AddressesFromStreamKey(key)

	require.Equal(t, receiverAddr, r)
	require.Equal(t, senderAddr, s)
}
