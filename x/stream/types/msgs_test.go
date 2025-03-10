package types_test

import (
	"testing"

	mathmod "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/unification-com/mainchain/x/stream/types"
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
	s := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	r := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	tests := []struct {
		deposit    sdk.Coin
		flowRate   int64
		receiver   sdk.AccAddress
		sender     sdk.AccAddress
		expectPass bool
	}{
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(10000)), 100, r, s, true},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(0)), 100, r, s, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(10000)), 0, r, s, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(10000)), 100, sdk.AccAddress{}, s, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(10000)), 100, r, sdk.AccAddress{}, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(100)), 100, r, s, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(10000)), 100, r, r, false},
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

func TestMsgClaimStream_Route(t *testing.T) {
	msg := types.MsgClaimStream{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgClaimStream_Type(t *testing.T) {
	msg := types.MsgClaimStream{}
	require.Equal(t, types.ClaimStreamAction, msg.Type())
}

func TestMsgClaimStream_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	receiverAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgClaimStream{Receiver: receiverAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(receiverAddr))
}

func TestMsgClaimStream_ValidateBasic(t *testing.T) {
	tests := []struct {
		sender     sdk.AccAddress
		receiver   sdk.AccAddress
		expectPass bool
	}{
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{sdk.AccAddress{}, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := types.NewMsgClaimStream(
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

//	MsgTopUpDeposit{}

func TestMsgTopUpDeposit_Route(t *testing.T) {
	msg := types.MsgTopUpDeposit{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgTopUpDeposit_Type(t *testing.T) {
	msg := types.MsgTopUpDeposit{}
	require.Equal(t, types.TopUpDepositAction, msg.Type())
}

func TestMsgTopUpDeposit_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	senderAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgCreateStream{Sender: senderAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(senderAddr))
}

func TestMsgTopUpDeposit_ValidateBasic(t *testing.T) {
	tests := []struct {
		deposit    sdk.Coin
		sender     sdk.AccAddress
		receiver   sdk.AccAddress
		expectPass bool
	}{
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(100)), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(100)), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress{}, false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(100)), sdk.AccAddress{}, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.NewCoin(sdk.DefaultBondDenom, mathmod.NewIntFromUint64(0)), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
	}

	for i, tc := range tests {
		msg := types.NewMsgTopUpDeposit(
			tc.receiver,
			tc.sender,
			tc.deposit,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

//	MsgUpdateFlowRate{}

func TestMsgUpdateFlowRate_Route(t *testing.T) {
	msg := types.MsgUpdateFlowRate{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgUpdateFlowRate_Type(t *testing.T) {
	msg := types.MsgUpdateFlowRate{}
	require.Equal(t, types.UpdateFlowRateAction, msg.Type())
}

func TestMsgUpdateFlowRate_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	senderAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgUpdateFlowRate{Sender: senderAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(senderAddr))
}

func TestMsgUpdateFlowRate_ValidateBasic(t *testing.T) {
	tests := []struct {
		flowRate   int64
		sender     sdk.AccAddress
		receiver   sdk.AccAddress
		expectPass bool
	}{
		{1, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{1, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress{}, false},
		{1, sdk.AccAddress{}, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{0, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
	}

	for i, tc := range tests {
		msg := types.NewMsgUpdateFlowRate(
			tc.receiver,
			tc.sender,
			tc.flowRate,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

//	MsgCancelStream{}

func TestMsgCancelStream_Route(t *testing.T) {
	msg := types.MsgCancelStream{}
	require.Equal(t, types.ModuleName, msg.Route())
}

func TestMsgCancelStream_Type(t *testing.T) {
	msg := types.MsgCancelStream{}
	require.Equal(t, types.CancelStreamAction, msg.Type())
}

func TestMsgMsgCancelStream_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	senderAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgCancelStream{Sender: senderAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(senderAddr))
}

func TestMsgCancelStream_ValidateBasic(t *testing.T) {
	tests := []struct {
		receiver   sdk.AccAddress
		sender     sdk.AccAddress
		expectPass bool
	}{
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{sdk.AccAddress{}, sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := types.NewMsgCancelStream(
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

// MsgUpdateParams{}
func TestMsgUpdateParams_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	senderAddr := sdk.AccAddress(pubKey2.Address())
	msg := types.MsgUpdateParams{Authority: senderAddr.String()}
	require.True(t, msg.GetSigners()[0].Equals(senderAddr))
}

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name            string
		msgUpdateParams types.MsgUpdateParams
		expFail         bool
		expError        string
	}{
		{
			"valid msg",
			types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params:    types.DefaultParams(),
			},
			false,
			"",
		},
		{
			"negative validator fee",
			types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: types.Params{
					ValidatorFee: mathmod.LegacyNewDecWithPrec(-1, 2),
				},
			},
			true,
			"validator fee cannot be negative:",
		},
		{
			"validator fee > 100%",
			types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: types.Params{
					ValidatorFee: mathmod.LegacyNewDecWithPrec(101, 2),
				},
			},
			true,
			"validator fee cannot be greater than 100% (1.00). Sent",
		},
		{
			"nil validator fee",
			types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: types.Params{
					ValidatorFee: mathmod.LegacyDec{},
				},
			},
			true,
			"validator fee cannot be nil",
		},
		{
			"Invalid authority",
			types.MsgUpdateParams{
				Authority: "invalid",
				Params:    types.DefaultParams(),
			},
			true,
			"invalid authority address",
		},
	}

	for _, tc := range tests {
		err := tc.msgUpdateParams.ValidateBasic()
		if tc.expFail {
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expError)
		} else {
			require.NoError(t, err)
		}
	}
}
