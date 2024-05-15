package keeper_test

import (
	simapp "github.com/unification-com/mainchain/app"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestAddAddressesToWhitelist(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(10)

	for _, addr := range testAddrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		require.NoError(t, err)
	}

	expectedErr := sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, sdk.AccAddress{})
	require.Equal(t, expectedErr.Error(), err.Error())
}

func TestRemoveAddressesToWhitelist(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(100)

	for _, addr := range testAddrs {
		_ = app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
	}

	for _, addr := range testAddrs {
		err := app.EnterpriseKeeper.RemoveAddressFromWhitelist(ctx, addr)
		require.NoError(t, err)
	}

	expectedErr := sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	err := app.EnterpriseKeeper.RemoveAddressFromWhitelist(ctx, sdk.AccAddress{})
	require.Equal(t, expectedErr.Error(), err.Error())
}

func TestAddressIsWhitelisted(t *testing.T) {
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	simapp.SetKeeperTestParamsAndDefaultValues(app, ctx)
	testAddrs := simapp.GenerateRandomTestAccounts(100)

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
