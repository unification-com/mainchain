package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis = 101

	CodeWrkChainDoesNotExist         = 201
	CodeWrkChainAlreadyRegistered    = 202
	CodeWrkChainBlockAlreadyRecorded = 203
	CodeNotWrkChainOwner             = 204

	CodeWrkChainInsufficientFee = 301
	CodeWrkChainTooMuchFee      = 302

	CodeWrkChainFeePayerNotOwner = 401
)

// ErrInvalidGenesis error for an invalid wrkchain GenesisState
func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrWrkChainDoesNotExist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainDoesNotExist, msg)
}

func ErrNotWrkChainOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNotWrkChainOwner, msg)
}

func ErrWrkChainAlreadyRegistered(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainAlreadyRegistered, msg)
}

func ErrWrkChainBlockAlreadyRecorded(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainBlockAlreadyRecorded, msg)
}

func ErrInsufficientWrkChainFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainInsufficientFee, msg)
}

func ErrTooMuchWrkChainFee(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainTooMuchFee, msg)
}

func ErrFeePayerNotOwner(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainFeePayerNotOwner, msg)
}
