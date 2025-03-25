package types

import (
	errorsmod "cosmossdk.io/errors"
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
	ErrInvalidGenesis                 = errorsmod.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrBeaconDoesNotExist             = errorsmod.Register(ModuleName, CodeBeaconDoesNotExist, "beacon does not exist")
	ErrNotBeaconOwner                 = errorsmod.Register(ModuleName, CodeNotBeaconOwner, "not beacon owner")
	ErrBeaconAlreadyRegistered        = errorsmod.Register(ModuleName, CodeBeaconAlreadyRegistered, "beacon already registered")
	ErrBeaconTimestampAlreadyRecorded = errorsmod.Register(ModuleName, CodeBeaconTimestampAlreadyRecorded, "beacon timestamp already recorded")
	ErrMissingData                    = errorsmod.Register(ModuleName, CodeMissingData, "missing data")
	ErrInsufficientBeaconFee          = errorsmod.Register(ModuleName, CodeBeaconInsufficientFee, "insufficient beacon fee")
	ErrTooMuchBeaconFee               = errorsmod.Register(ModuleName, CodeBeaconTooMuchFee, "too much beacon fee")
	ErrFeePayerNotOwner               = errorsmod.Register(ModuleName, CodeBeaconFeePayerNotOwner, "fee payer is not beacon owner")
	ErrIncorrectFeeDenomination       = errorsmod.Register(ModuleName, CodeBeaconIncorrectFeeDenom, "incorrect beacon fee denomination")
	ErrContentTooLarge                = errorsmod.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrExceedsMaxStorage              = errorsmod.Register(ModuleName, CodeExceedsMaxStorage, "exceeds max storage")
)
