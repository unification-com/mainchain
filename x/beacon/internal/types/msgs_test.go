package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestMsgRegisterBeacon_Route(t *testing.T) {
	msg := MsgRegisterBeacon{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRegisterBeacon_Type(t *testing.T) {
	msg := MsgRegisterBeacon{}
	require.Equal(t, "register_beacon", msg.Type())
}

func TestMsgRegisterBeacon_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	ownerAddr := sdk.AccAddress(pubKey2.Address())
	msg := MsgRegisterBeacon{Owner: ownerAddr}
	require.True(t, msg.GetSigners()[0].Equals(ownerAddr))
}

func TestMsgRecordBeaconTimestamp_Route(t *testing.T) {
	msg := MsgRecordBeaconTimestamp{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRecordBeaconTimestamp_Type(t *testing.T) {
	msg := MsgRecordBeaconTimestamp{}
	require.Equal(t, "record_beacon_timestamp", msg.Type())
}

func TestMsgRecordBeaconTimestamp_GetSigners(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	ownerAddr := sdk.AccAddress(pubKey2.Address())
	msg := MsgRecordBeaconTimestamp{Owner: ownerAddr}
	require.True(t, msg.GetSigners()[0].Equals(ownerAddr))
}

func TestMsgRegisterBeacon(t *testing.T) {

	tests := []struct {
		moniker    string
		name       string
		owner      sdk.AccAddress
		expectPass bool
	}{
		{"b1", "BEACON 1", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"", "BEACON 1", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"b2", "", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"b3", "BEACON 3", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"b4", "BEACON 4", sdk.AccAddress{}, false},
	}

	for i, tc := range tests {
		msg := NewMsgRegisterBeacon(
			tc.moniker,
			tc.name,
			tc.owner,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgRecordBeaconTimestamp(t *testing.T) {

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	emptyAddr := sdk.AccAddress{}

	tests := []struct {
		beaconID   uint64
		subTime    uint64
		hash       string
		Owner      sdk.AccAddress
		expectPass bool
	}{
		{1, 1234, "beaconhash", addr, true},
		{1, 1, "beaconhash", emptyAddr, false},
		{1, 1, "", addr, false},
		{1, 0, "beaconhash", addr, false},
		{0, 1, "beaconhash", addr, false},
		{0, 0, "", addr, false},
		{0, 0, "", emptyAddr, false},
	}

	for i, tc := range tests {
		msg := NewMsgRecordBeaconTimestamp(
			tc.beaconID,
			tc.hash,
			tc.subTime,
			tc.Owner,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}
