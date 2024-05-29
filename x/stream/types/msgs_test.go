package types_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/stream/types"
	"testing"
)

//	MsgCreateStream{}

func TestMsgCreateStream_Route(t *testing.T) {
	msg := types.MsgCreateStream{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgCreateStream_Type(t *testing.T) {
	msg := types.MsgCreateStream{}
	require.Equal(t, types.CreateStreamAction, msg.Type())
}

func TestMsgCreateStream_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	senderAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgCreateStream{Sender: senderAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(senderAddr))
}

func TestMsgCreateStream_ValidateBasic(t *testing.T) {
	tests := []struct {
		deposit    sdk.Coin
		flowRate   int64
		receiver   sdk.AccAddress
		sender     sdk.AccAddress
		expectPass bool
	}{
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(10000)), 100, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(0)), 100, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(10000)), 0, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(10000)), 100, sdk.AccAddress{}, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(10000)), 100, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress{}, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromUint64(100)), 100, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
	}

	for i, tc := range tests {
		msg := types.NewMsgCreateStream(
			tc.deposit,
			tc.flowRate,
			tc.receiver,
			tc.sender,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

//	MsgClaimStream{}
//	MsgTopUpDeposit{}
//	MsgUpdateFlowRate{}
//	MsgCancelStream{}
//	MsgUpdateParams{}
