package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	ErrInvalidGenesis               = sdkerrors.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrWrkChainDoesNotExist         = sdkerrors.Register(ModuleName, CodeWrkChainDoesNotExist, "wrkchain does not exist")
	ErrNotWrkChainOwner             = sdkerrors.Register(ModuleName, CodeNotWrkChainOwner, "not wrkchain owner")
	ErrWrkChainAlreadyRegistered    = sdkerrors.Register(ModuleName, CodeWrkChainAlreadyRegistered, "wrkchain already registered")
	ErrMissingData                  = sdkerrors.Register(ModuleName, CodeMissingData, "missing data")
	ErrInvalidData                  = sdkerrors.Register(ModuleName, CodeInvalidData, "invalid data")
	ErrWrkChainBlockAlreadyRecorded = sdkerrors.Register(ModuleName, CodeWrkChainBlockAlreadyRecorded, "wrkchain hashes already recorded")
	ErrInsufficientWrkChainFee      = sdkerrors.Register(ModuleName, CodeWrkChainInsufficientFee, "insufficient wrkchain fee")
	ErrTooMuchWrkChainFee           = sdkerrors.Register(ModuleName, CodeWrkChainTooMuchFee, "too much wrkchain fee")
	ErrFeePayerNotOwner             = sdkerrors.Register(ModuleName, CodeWrkChainFeePayerNotOwner, "fee payer not wrkchain owner")
	ErrIncorrectFeeDenomination     = sdkerrors.Register(ModuleName, CodeWrkChainIncorrectFeeDenom, "incorrect wrkchain fee doenomination")
	ErrContentTooLarge              = sdkerrors.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrExceedsMaxStorage            = sdkerrors.Register(ModuleName, CodeExceedsMaxStorage, "exceeds max storage")
	ErrNewHeightMustBeHigher        = sdkerrors.Register(ModuleName, CodeNewHeightMustBeHigher, "invalid height")
)
