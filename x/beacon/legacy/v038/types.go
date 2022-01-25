package v038

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// ModuleName module name
	ModuleName = "beacon"

	RouterKey = ModuleName // defined in keys.go file

	RegisterAction = "register_beacon"
	RecordAction   = "record_beacon_timestamp"
)

type (
	Beacons          []Beacon
	BeaconTimestamps []BeaconTimestamp

	Params struct {
		FeeRegister uint64 `json:"fee_register" yaml:"fee_register"` // Fee for registering a BEACON
		FeeRecord   uint64 `json:"fee_record" yaml:"fee_record"`     // Fee for recording timestamps for a BEACON
		Denom       string `json:"denom" yaml:"denom"`               // Fee denomination
	}

	Beacon struct {
		BeaconID        uint64         `json:"beacon_id"`
		Moniker         string         `json:"moniker"`
		Name            string         `json:"name"`
		LastTimestampID uint64         `json:"last_timestamp_id"`
		Owner           sdk.AccAddress `json:"owner"`
	}

	BeaconTimestamp struct {
		BeaconID    uint64         `json:"beacon_id"`
		TimestampID uint64         `json:"timestamp_id"`
		SubmitTime  uint64         `json:"submit_time"`
		Hash        string         `json:"hash"`
		Owner       sdk.AccAddress `json:"owner"`
	}

	GenesisState struct {
		Params           Params         `json:"params" yaml:"params"`                         // beacon params
		StartingBeaconID uint64         `json:"starting_beacon_id" yaml:"starting_beacon_id"` // should be 1
		Beacons          []BeaconExport `json:"registered_beacons" yaml:"registered_beacons"`
	}

	BeaconExport struct {
		Beacon           Beacon            `json:"beacon" yaml:"beacon"`
		BeaconTimestamps []BeaconTimestamp `json:"timestamps" yaml:"timestamps"`
	}
)

// MsgRegisterBeacon defines a RegisterBeacon message
type MsgRegisterBeacon struct {
	Moniker    string         `json:"moniker"`
	BeaconName string         `json:"name"`
	Owner      sdk.AccAddress `json:"owner"`
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
		return fmt.Errorf("moniker and name cannot be empty")
	}

	if len(msg.BeaconName) > 128 {
		return fmt.Errorf("name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return fmt.Errorf("moniker too big. 64 character limit")
	}

	return nil
}

// --- Record a BEACON timestamp hash Msg ---

// MsgRecordBeaconTimestamp defines a RecordBeaconTimestamp message
type MsgRecordBeaconTimestamp struct {
	BeaconID   uint64         `json:"beacon_id"`
	Hash       string         `json:"hash"`
	SubmitTime uint64         `json:"submit_time"`
	Owner      sdk.AccAddress `json:"owner"`
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
		return fmt.Errorf("id must be greater than zero")
	}
	if len(msg.Hash) == 0 {
		return fmt.Errorf("hash cannot be empty")
	}
	if msg.SubmitTime == 0 {
		return fmt.Errorf("submit time cannot be zero")
	}
	if len(msg.Hash) > 66 {
		return fmt.Errorf("hash too big. 66 character limit")
	}

	return nil
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgRegisterBeacon{}, "beacon/RegisterBeacon", nil)
	cdc.RegisterConcrete(MsgRecordBeaconTimestamp{}, "beacon/RecordBeaconTimestamp", nil)
}
