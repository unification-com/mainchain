package types_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	simapp "github.com/unification-com/mainchain/app"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"testing"
)

func TestMsgUndPurchaseOrder_Route(t *testing.T) {
	msg := types.MsgUndPurchaseOrder{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgUndPurchaseOrder_Type(t *testing.T) {
	msg := types.MsgUndPurchaseOrder{}
	require.Equal(t, types.PurchaseAction, msg.Type())
}

func TestMsgProcessUndPurchaseOrder_Route(t *testing.T) {
	msg := types.MsgProcessUndPurchaseOrder{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgProcessUndPurchaseOrder_Type(t *testing.T) {
	msg := types.MsgProcessUndPurchaseOrder{}
	require.Equal(t, types.ProcessAction, msg.Type())
}

func TestMsgWhitelistAddress_Route(t *testing.T) {
	msg := types.MsgWhitelistAddress{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgWhitelistAddress_Type(t *testing.T) {
	msg := types.MsgWhitelistAddress{}
	require.Equal(t, types.WhitelistAddressAction, msg.Type())
}

func TestMsgUndPurchaseOrder_Validate(t *testing.T) {
	tests := []struct {
		amount     sdk.Coin
		purchaser  string
		expectPass bool
	}{
		{
			sdk.NewInt64Coin(simapp.TestDenomination, 1),
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			true,
		},
		{
			sdk.NewInt64Coin(simapp.TestDenomination, 0),
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			false,
		},
		{
			sdk.NewInt64Coin(simapp.TestDenomination, 1),
			"rubbish",
			false,
		},
		{
			sdk.NewInt64Coin(simapp.TestDenomination, 0),
			"rubbish",
			false,
		},
	}

	for i, tc := range tests {
		msg := types.MsgUndPurchaseOrder{
			Purchaser: tc.purchaser,
			Amount:    tc.amount,
		}

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgProcessUndPurchaseOrder_Validate(t *testing.T) {
	tests := []struct {
		id         uint64
		decision   types.PurchaseOrderStatus
		signer     string
		expectPass bool
	}{
		{
			1,
			2,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			true,
		},
		{
			0,
			2,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			false,
		},
		{
			1,
			99,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			false,
		},
		{
			1,
			2,
			"rubbish",
			false,
		},
	}

	for i, tc := range tests {
		msg := types.MsgProcessUndPurchaseOrder{
			PurchaseOrderId: tc.id,
			Decision:        tc.decision,
			Signer:          tc.signer,
		}

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgWhitelistAddress_Validate(t *testing.T) {
	tests := []struct {
		action     types.WhitelistAction
		address    string
		signer     string
		expectPass bool
	}{
		{
			1,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			true,
		},
		{
			0,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			false,
		},
		{
			1,
			"rubbish",
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			false,
		},
		{
			1,
			sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
			"rubbish",
			false,
		},
	}

	for i, tc := range tests {
		msg := types.MsgWhitelistAddress{
			Address: tc.address,
			Signer:  tc.signer,
			Action:  tc.action,
		}

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
