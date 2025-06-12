package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/stretchr/testify/require"

	appparams "github.com/unification-com/mainchain/app/params"
	"github.com/unification-com/mainchain/x/stream/types"
)

func TestAddressFromStreamsStore(t *testing.T) {
	appparams.SetAddressPrefixes()
	receiverAddr, err := sdk.AccAddressFromBech32("und1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc6xaj4c")
	require.NoError(t, err)
	receiverAddrLen := len(receiverAddr)
	require.Equal(t, 20, receiverAddrLen)

	senderAddr, err := sdk.AccAddressFromBech32("und139f7kncmglres2nf3h4hc4tade85ekfrq4mkyj")
	require.NoError(t, err)
	senderAddrLen := len(senderAddr)
	require.Equal(t, 20, senderAddrLen)

	key := types.GetStreamKey(receiverAddr, senderAddr)

	require.Len(t, key, len(types.StreamKeyPrefix)+len(address.MustLengthPrefix(receiverAddr))+len(address.MustLengthPrefix(senderAddr)))

	r, s := types.AddressesFromStreamKey(key)

	require.Equal(t, receiverAddr, r)
	require.Equal(t, senderAddr, s)
}
