package types

import (
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgRegisterBeacon_Route(t *testing.T) {
	msg := MsgRegisterBeacon{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRegisterBeacon_Type(t *testing.T) {
	msg := MsgRegisterBeacon{}
	require.Equal(t, RegisterAction, msg.Type())
}

func TestMsgRecordBeaconTimestamp_Route(t *testing.T) {
	msg := MsgRecordBeaconTimestamp{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgRecordBeaconTimestamp_Type(t *testing.T) {
	msg := MsgRecordBeaconTimestamp{}
	require.Equal(t, RecordAction, msg.Type())
}

func TestMsgPurchaseBeaconStateStorage_Route(t *testing.T) {
	msg := MsgPurchaseBeaconStateStorage{}
	require.Equal(t, ModuleName, msg.Route())
}

func TestMsgPurchaseBeaconStateStorage_Type(t *testing.T) {
	msg := MsgPurchaseBeaconStateStorage{}
	require.Equal(t, PurchaseStorageAction, msg.Type())
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
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), true},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a51", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
		{"c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", "c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5c14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a51", sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()), false},
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
		{1, 12345, "0xc14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a5", addr, true},
		{1, 12346, "0xc14cb7f5c98846be8668e95e99312df0c74391dd328ef07daf66de05920c44a51", addr, false},
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

func TestMsgPurchaseBeaconStateStorage(t *testing.T) {

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	emptyAddr := sdk.AccAddress{}

	tests := []struct {
		beaconID   uint64
		number     uint64
		Owner      sdk.AccAddress
		expectPass bool
	}{
		{1, 10, addr, true},
		{1, 10, emptyAddr, false},
		{1, 0, addr, false},
		{0, 1, addr, false},
		{0, 0, addr, false},
		{0, 0, emptyAddr, false},
	}

	for i, tc := range tests {
		msg := NewMsgPurchaseBeaconStateStorage(
			tc.beaconID,
			tc.number,
			tc.Owner,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgRegisterBeaconGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("addr1")
	msg := NewMsgRegisterBeacon("testbeacon", "testbeaconname", addr)
	pc := codec.NewProtoCodec(types.NewInterfaceRegistry())
	res, err := pc.MarshalAminoJSON(msg)
	require.NoError(t, err)
	expected := `{"type":"beacon/RegisterBeacon","value":{"moniker":"testbeacon","name":"testbeaconname","owner":"cosmos1v9jxgu33kfsgr5"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgPurchaseBeaconStateStorageGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("addr1")
	msg := NewMsgPurchaseBeaconStateStorage(1, 1000, addr)
	pc := codec.NewProtoCodec(types.NewInterfaceRegistry())
	res, err := pc.MarshalAminoJSON(msg)
	require.NoError(t, err)
	expected := `{"type":"beacon/PurchaseBeaconStateStorage","value":{"beacon_id":"1","number":"1000","owner":"cosmos1v9jxgu33kfsgr5"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgRecordBeaconTimestampGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("addr1")
	msg := NewMsgRecordBeaconTimestamp(1, "abc123", 1000, addr)
	pc := codec.NewProtoCodec(types.NewInterfaceRegistry())
	res, err := pc.MarshalAminoJSON(msg)
	require.NoError(t, err)
	expected := `{"type":"beacon/RecordBeaconTimestamp","value":{"beacon_id":"1","hash":"abc123","owner":"cosmos1v9jxgu33kfsgr5","submit_time":"1000"}}`
	require.Equal(t, expected, string(res))
}
