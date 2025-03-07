package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	RegisterAction        = "register_beacon"
	RecordAction          = "record_beacon_timestamp"
	PurchaseStorageAction = "purchase_beacon_storage"
)

var (
	_ sdk.Msg = &MsgRegisterBeacon{}
	_ sdk.Msg = &MsgPurchaseBeaconStateStorage{}
	_ sdk.Msg = &MsgRecordBeaconTimestamp{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// --- Register a BEACON Msg ---

// NewMsgRegisterBeacon is a constructor function for MsgRegisterBeacon
func NewMsgRegisterBeacon(moniker string, beaconName string, owner sdk.AccAddress) *MsgRegisterBeacon {
	return &MsgRegisterBeacon{
		Moniker: moniker,
		Name:    beaconName,
		Owner:   owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgRegisterBeacon) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRegisterBeacon) Type() string { return RegisterAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterBeacon) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	if len(msg.Moniker) == 0 || len(msg.Name) == 0 {
		return errorsmod.Wrap(ErrMissingData, "moniker and name cannot be empty")
	}

	if len(msg.Name) > 128 {
		return errorsmod.Wrap(ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return errorsmod.Wrap(ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterBeacon) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterBeacon) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Record a BEACON timestamp hash Msg ---

// NewMsgRecordBeaconTimestamp is a constructor function for MsgRecordBeaconTimestamp
func NewMsgRecordBeaconTimestamp(
	beaconId uint64,
	hash string,
	subTime uint64,
	owner sdk.AccAddress) *MsgRecordBeaconTimestamp {

	return &MsgRecordBeaconTimestamp{
		BeaconId:   beaconId,
		Hash:       hash,
		SubmitTime: subTime,
		Owner:      owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgRecordBeaconTimestamp) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecordBeaconTimestamp) Type() string { return RecordAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecordBeaconTimestamp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}
	if msg.BeaconId == 0 {
		return errorsmod.Wrap(ErrMissingData, "id must be greater than zero")
	}
	if len(msg.Hash) == 0 {
		return errorsmod.Wrap(ErrMissingData, "hash cannot be empty")
	}
	if msg.SubmitTime == 0 {
		return errorsmod.Wrap(ErrMissingData, "submit time cannot be zero")
	}
	if len(msg.Hash) > 66 {
		return errorsmod.Wrap(ErrContentTooLarge, "hash too big. 66 character limit")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRecordBeaconTimestamp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgRecordBeaconTimestamp) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Purchase state storage Msg ---

// NewMsgRecordBeaconTimestamp is a constructor function for MsgRecordBeaconTimestamp
func NewMsgPurchaseBeaconStateStorage(
	beaconId uint64,
	number uint64,
	owner sdk.AccAddress) *MsgPurchaseBeaconStateStorage {

	return &MsgPurchaseBeaconStateStorage{
		BeaconId: beaconId,
		Number:   number,
		Owner:    owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgPurchaseBeaconStateStorage) Route() string { return RouterKey }

// Type should return the action
func (msg MsgPurchaseBeaconStateStorage) Type() string { return PurchaseStorageAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgPurchaseBeaconStateStorage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}
	if msg.BeaconId == 0 {
		return errorsmod.Wrap(ErrMissingData, "id must be greater than zero")
	}
	if msg.Number == 0 {
		return errorsmod.Wrap(ErrMissingData, "number cannot be zero")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPurchaseBeaconStateStorage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgPurchaseBeaconStateStorage) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Modify Params Msg Type ---

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}
