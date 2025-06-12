package types

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	CodeInvalidGenesis = 101

	CodeWrkChainDoesNotExist         = 201
	CodeWrkChainAlreadyRegistered    = 202
	CodeWrkChainBlockAlreadyRecorded = 203
	CodeNotWrkChainOwner             = 204
	CodeMissingData                  = 205
	CodeInvalidData                  = 206
	CodeContentTooLarge              = 207
	CodeExceedsMaxStorage            = 208
	CodeNewHeightMustBeHigher        = 209

	CodeWrkChainInsufficientFee   = 301
	CodeWrkChainTooMuchFee        = 302
	CodeWrkChainIncorrectFeeDenom = 303

	CodeWrkChainFeePayerNotOwner = 401
)

var (
	ErrInvalidGenesis               = errorsmod.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrWrkChainDoesNotExist         = errorsmod.Register(ModuleName, CodeWrkChainDoesNotExist, "wrkchain does not exist")
	ErrNotWrkChainOwner             = errorsmod.Register(ModuleName, CodeNotWrkChainOwner, "not wrkchain owner")
	ErrWrkChainAlreadyRegistered    = errorsmod.Register(ModuleName, CodeWrkChainAlreadyRegistered, "wrkchain already registered")
	ErrMissingData                  = errorsmod.Register(ModuleName, CodeMissingData, "missing data")
	ErrInvalidData                  = errorsmod.Register(ModuleName, CodeInvalidData, "invalid data")
	ErrWrkChainBlockAlreadyRecorded = errorsmod.Register(ModuleName, CodeWrkChainBlockAlreadyRecorded, "wrkchain hashes already recorded")
	ErrInsufficientWrkChainFee      = errorsmod.Register(ModuleName, CodeWrkChainInsufficientFee, "insufficient wrkchain fee")
	ErrTooMuchWrkChainFee           = errorsmod.Register(ModuleName, CodeWrkChainTooMuchFee, "too much wrkchain fee")
	ErrFeePayerNotOwner             = errorsmod.Register(ModuleName, CodeWrkChainFeePayerNotOwner, "fee payer not wrkchain owner")
	ErrIncorrectFeeDenomination     = errorsmod.Register(ModuleName, CodeWrkChainIncorrectFeeDenom, "incorrect wrkchain fee doenomination")
	ErrContentTooLarge              = errorsmod.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrExceedsMaxStorage            = errorsmod.Register(ModuleName, CodeExceedsMaxStorage, "exceeds max storage")
	ErrNewHeightMustBeHigher        = errorsmod.Register(ModuleName, CodeNewHeightMustBeHigher, "invalid height")
)
