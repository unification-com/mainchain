package beacon

import (
	"github.com/unification-com/mainchain-cosmos/x/beacon/exported"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = types.DefaultParamspace

	DefaultCodespace = types.DefaultCodespace

	QueryParameters = keeper.QueryParameters
)

var (
	NewKeeper       = keeper.NewKeeper
	NewQuerier      = keeper.NewQuerier
	NewGenesisState = types.NewGenesisState

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	NewQueryBeaconParams = types.NewQueryBeaconParams

	// Errors
	ErrInvalidGenesis                 = types.ErrInvalidGenesis
	ErrBeaconAlreadyRegistered        = types.ErrBeaconAlreadyRegistered
	ErrBeaconDoesNotExist             = types.ErrBeaconDoesNotExist
	ErrNotBeaconOwner                 = types.ErrNotBeaconOwner
	ErrBeaconTimestampAlreadyRecorded = types.ErrBeaconTimestampAlreadyRecorded
	ErrInsufficientBeaconFee          = types.ErrInsufficientBeaconFee
	ErrTooMuchBeaconFee               = types.ErrTooMuchBeaconFee
	ErrFeePayerNotOwner               = types.ErrFeePayerNotOwner

	// Events
	EventTypeRegisterBeacon         = types.EventTypeRegisterBeacon
	EventTypeRecordBeaconTimestamp  = types.EventTypeRecordBeaconTimestamp
	AttributeValueCategory          = types.AttributeValueCategory
	AttributeKeyOwner               = types.AttributeKeyOwner
	AttributeKeyBeaconId            = types.AttributeKeyBeaconId
	AttributeKeyBeaconMoniker       = types.AttributeKeyBeaconMoniker
	AttributeKeyBeaconName          = types.AttributeKeyBeaconName
	AttributeKeyTimestampID         = types.AttributeKeyTimestampID
	AttributeKeyTimestampHash       = types.AttributeKeyTimestampHash
	AttributeKeyTimestampSubmitTime = types.AttributeKeyTimestampSubmitTime

	GetBeaconIDBytes    = types.GetBeaconIDBytes
	GetTimestampIDBytes = types.GetTimestampIDBytes

	CheckIsBeaconTx = exported.CheckIsBeaconTx

)

type (
	Keeper = keeper.Keeper

	GenesisState = types.GenesisState

	BeaconExport = types.BeaconExport
	BeaconTimestamp = types.BeaconTimestamp

	// Msgs
	MsgRegisterBeacon        = types.MsgRegisterBeacon
	MsgRecordBeaconTimestamp = types.MsgRecordBeaconTimestamp
)
