package types

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgUndPurchaseOrder(t *testing.T) {

	denom := "testc"
	tests := []struct {
		purchaser  sdk.AccAddress
		amount     sdk.Coin
		expectPass bool
	}{
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.NewInt64Coin(denom, 1), true},
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.NewInt64Coin(denom, 10), true},
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.NewInt64Coin(denom, 0), false},
		{sdk.AccAddress{}, sdk.NewInt64Coin(denom, 1), false},
		{sdk.AccAddress{}, sdk.NewInt64Coin(denom, 0), false},
	}

	for i, tc := range tests {
		msg := NewMsgUndPurchaseOrder(
			tc.purchaser,
			tc.amount,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgProcessUndPurchaseOrder(t *testing.T) {
	tests := []struct {
		poID       uint64
		decision   PurchaseOrderStatus
		signer     sdk.AccAddress
		expectPass bool
	}{
		{1, StatusAccepted, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{1, StatusRejected, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{1, StatusCompleted, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{1, StatusNil, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{1, StatusRaised, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{0, StatusAccepted, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{1, StatusAccepted, sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := NewMsgProcessUndPurchaseOrder(
			tc.poID,
			tc.decision,
			tc.signer,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
