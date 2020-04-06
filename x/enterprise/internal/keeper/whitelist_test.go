package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
	"testing"
)

func TestAddAddressesToWhitelist(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	numTests := 100
	testAddrs := GenerateRandomAddresses(numTests)

	for _, addr := range testAddrs {
		err := keeper.AddAddressToWhitelist(ctx, addr)
		require.NoError(t, err)
	}
}

func TestRemoveAddressesToWhitelist(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	numTests := 100
	testAddrs := GenerateRandomAddresses(numTests)

	for _, addr := range testAddrs {
		_ = keeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		err := keeper.RemoveAddressFromWhitelist(ctx, addr)
		require.NoError(t, err)
	}
}

func TestOnlyAuthotisedAddressesCanModify(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	numTests := 100
	testAddrs := GenerateRandomAddresses(numTests)

	expectedErr := sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist")

	for _, addr := range testAddrs {
		err := keeper.ProcessWhitelistAction(ctx, addr, types.WhitelistActionAdd, TestAddrs[1])
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())

		_ = keeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		err := keeper.ProcessWhitelistAction(ctx, addr, types.WhitelistActionRemove, TestAddrs[1])
		require.Equal(t, expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
	}
}

func TestAddressIsWhitelisted(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)
	numTests := 100
	testAddrs := GenerateRandomAddresses(numTests)

	for _, addr := range testAddrs {
		isWhitelisted := keeper.AddressIsWhitelisted(ctx, addr)
		require.False(t, isWhitelisted)
	}

	for _, addr := range testAddrs {
		_ = keeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		isWhitelisted := keeper.AddressIsWhitelisted(ctx, addr)
		require.True(t, isWhitelisted)
	}
}

func TestProcessWhitelistAction(t *testing.T) {
	ctx, _, keeper, _, _ := createTestInput(t, false, 100)

	testCases := []struct {
		addr        sdk.AccAddress
		action      types.WhitelistAction
		expectedErr error
		signer      sdk.AccAddress
	}{
		{sdk.AccAddress{}, types.WhitelistActionAdd, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty"), EntSignerAddr},
		{sdk.AccAddress{}, types.WhitelistActionRemove, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty"), EntSignerAddr},
		{TestAddrs[0], types.WhitelistActionRemove, sdkerrors.Wrap(types.ErrAddressNotWhitelisted, fmt.Sprintf("%s not whitelisted", TestAddrs[0])), EntSignerAddr},
		{TestAddrs[0], types.WhitelistActionAdd, nil, EntSignerAddr},
		{TestAddrs[0], types.WhitelistActionRemove, nil, EntSignerAddr},
		{TestAddrs[0], 0x03, sdkerrors.Wrap(types.ErrInvalidWhitelistAction, "action should be add or remove"), EntSignerAddr},
		{TestAddrs[0], 0x04, sdkerrors.Wrap(types.ErrInvalidWhitelistAction, "action should be add or remove"), EntSignerAddr},
		{TestAddrs[0], types.WhitelistActionAdd, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorised signer modifying whitelist"), TestAddrs[1]},
	}

	for _, tc := range testCases {
		err := keeper.ProcessWhitelistAction(ctx, tc.addr, tc.action, tc.signer)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
	}
}
