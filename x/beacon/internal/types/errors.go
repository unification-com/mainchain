package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis = 101

	CodeBeaconDoesNotExist         = 201
	CodeBeaconAlreadyRegistered    = 202
	CodeBeaconBlockAlreadyRecorded = 203
	CodeNotBeaconOwner             = 204

	CodeBeaconInsufficientFee   = 301
	CodeBeaconTooMuchFee        = 302
	CodeBeaconIncorrectFeeDenom = 303

	CodeBeaconFeePayerNotOwner = 401
)

// ErrInvalidGenesis error for an invalid beacon GenesisState
func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrBeaconDoesNotExist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconDoesNotExist, msg)
}

func ErrNotBeaconOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNotBeaconOwner, msg)
}

func ErrBeaconAlreadyRegistered(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconAlreadyRegistered, msg)
}

func ErrBeaconTimestampAlreadyRecorded(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconBlockAlreadyRecorded, msg)
}

func ErrInsufficientBeaconFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconInsufficientFee, msg)
}

func ErrTooMuchBeaconFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconTooMuchFee, msg)
}

func ErrFeePayerNotOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconFeePayerNotOwner, msg)
}

func ErrIncorrectFeeDenomination(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBeaconIncorrectFeeDenom, msg)
}
