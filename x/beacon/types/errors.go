package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodeInvalidGenesis = 101

	CodeBeaconDoesNotExist             = 201
	CodeBeaconAlreadyRegistered        = 202
	CodeBeaconTimestampAlreadyRecorded = 203
	CodeNotBeaconOwner                 = 204
	CodeMissingData                    = 205
	CodeContentTooLarge                = 206
	CodeExceedsMaxStorage              = 207

	CodeBeaconInsufficientFee   = 301
	CodeBeaconTooMuchFee        = 302
	CodeBeaconIncorrectFeeDenom = 303

	CodeBeaconFeePayerNotOwner = 401
)

var (
	// ErrInvalidGenesis error for an invalid beacon GenesisState
	ErrInvalidGenesis                 = sdkerrors.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrBeaconDoesNotExist             = sdkerrors.Register(ModuleName, CodeBeaconDoesNotExist, "beacon does not exist")
	ErrNotBeaconOwner                 = sdkerrors.Register(ModuleName, CodeNotBeaconOwner, "not beacon owner")
	ErrBeaconAlreadyRegistered        = sdkerrors.Register(ModuleName, CodeBeaconAlreadyRegistered, "beacon already registered")
	ErrBeaconTimestampAlreadyRecorded = sdkerrors.Register(ModuleName, CodeBeaconTimestampAlreadyRecorded, "beacon timestamp already recorded")
	ErrMissingData                    = sdkerrors.Register(ModuleName, CodeMissingData, "missing data")
	ErrInsufficientBeaconFee          = sdkerrors.Register(ModuleName, CodeBeaconInsufficientFee, "insufficient beacon fee")
	ErrTooMuchBeaconFee               = sdkerrors.Register(ModuleName, CodeBeaconTooMuchFee, "too much beacon fee")
	ErrFeePayerNotOwner               = sdkerrors.Register(ModuleName, CodeBeaconFeePayerNotOwner, "fee payer is not beacon owner")
	ErrIncorrectFeeDenomination       = sdkerrors.Register(ModuleName, CodeBeaconIncorrectFeeDenom, "incorrect beacon fee denomination")
	ErrContentTooLarge                = sdkerrors.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrExceedsMaxStorage              = sdkerrors.Register(ModuleName, CodeExceedsMaxStorage, "exceeds max storage")
)
