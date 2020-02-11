package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	RegisterAction = "register_beacon"
	RecordAction   = "record_beacon_timestamp"
)

// --- Register a BEACON Msg ---

// MsgRegisterBeacon defines a RegisterBeacon message
type MsgRegisterBeacon struct {
	Moniker    string         `json:"moniker"`
	BeaconName string         `json:"name"`
	Owner      sdk.AccAddress `json:"owner"`
}

// NewMsgRegisterBeacon is a constructor function for MsgRegisterBeacon
func NewMsgRegisterBeacon(moniker string, beaconName string, owner sdk.AccAddress) MsgRegisterBeacon {
	return MsgRegisterBeacon{
		Moniker:    moniker,
		BeaconName: beaconName,
		Owner:      owner,
	}
}

// Route should return the name of the module
func (msg MsgRegisterBeacon) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRegisterBeacon) Type() string { return RegisterAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterBeacon) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Moniker) == 0 || len(msg.BeaconName) == 0 {
		return sdkerrors.Wrap(ErrMissingData, "moniker and name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterBeacon) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterBeacon) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// --- Record a BEACON timestamp hash Msg ---

// MsgRecordBeaconTimestamp defines a RecordBeaconTimestamp message
type MsgRecordBeaconTimestamp struct {
	BeaconID   uint64         `json:"beacon_id"`
	Hash       string         `json:"hash"`
	SubmitTime uint64         `json:"submit_time"`
	Owner      sdk.AccAddress `json:"owner"`
}

// NewMsgRecordBeaconTimestamp is a constructor function for MsgRecordBeaconTimestamp
func NewMsgRecordBeaconTimestamp(
	beaconId uint64,
	hash string,
	subTime uint64,
	owner sdk.AccAddress) MsgRecordBeaconTimestamp {

	return MsgRecordBeaconTimestamp{
		BeaconID:   beaconId,
		Hash:       hash,
		SubmitTime: subTime,
		Owner:      owner,
	}
}

// Route should return the name of the module
func (msg MsgRecordBeaconTimestamp) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecordBeaconTimestamp) Type() string { return RecordAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecordBeaconTimestamp) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if msg.BeaconID == 0 {
		return sdkerrors.Wrap(ErrMissingData, "id must be greater than zero")
	}
	if len(msg.Hash) == 0 {
		return sdkerrors.Wrap(ErrMissingData, "hash cannot be empty")
	}
	if msg.SubmitTime == 0 {
		return sdkerrors.Wrap(ErrMissingData, "submit time cannot be zero")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRecordBeaconTimestamp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRecordBeaconTimestamp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
