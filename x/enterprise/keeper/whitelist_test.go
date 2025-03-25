package keeper_test

import (
	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestAddAddressesToWhitelist(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(10)

	for _, addr := range testAddrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		require.NoError(t, err)
	}

	expectedErr := errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, sdk.AccAddress{})
	require.Equal(t, expectedErr.Error(), err.Error())
}

func TestRemoveAddressesToWhitelist(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddrs {
		_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		err := app.EnterpriseKeeper.RemoveAddressFromWhitelist(ctx, addr)
		require.NoError(t, err)
	}

	expectedErr := errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	err := app.EnterpriseKeeper.RemoveAddressFromWhitelist(ctx, sdk.AccAddress{})
	require.Equal(t, expectedErr.Error(), err.Error())
}

func TestAddressIsWhitelisted(t *testing.T) {
	app := simapphelpers.Setup(t)
	ctx := app.BaseApp.NewContext(false)
	testAddrs := simapphelpers.GenerateRandomTestAccounts(100)

	for _, addr := range testAddrs {
		isWhitelisted := app.EnterpriseKeeper.AddressIsWhitelisted(ctx, addr)
		require.False(t, isWhitelisted)
	}

	for _, addr := range testAddrs {
		_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		isWhitelisted := app.EnterpriseKeeper.AddressIsWhitelisted(ctx, addr)
		require.True(t, isWhitelisted)
	}

	isWhitelisted := app.EnterpriseKeeper.AddressIsWhitelisted(ctx, sdk.AccAddress{})
	require.False(t, isWhitelisted)
}
