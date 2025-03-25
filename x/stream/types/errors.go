package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	CodeInvalidGenesis = 101

	CodeStreamDoesNotExist  = 201
	CodeStreamAlreadyExists = 202
	CodeMissingData         = 203
	CodeContentTooLarge     = 204
	CodeInvalidData         = 205
	CodeNotCancellable      = 206

	CodeInsufficientDeposit = 301
)

// x/stream module sentinel errors
var (
	ErrStreamDoesNotExist   = errorsmod.Register(ModuleName, CodeStreamDoesNotExist, "stream does not exist")
	ErrStreamExists         = errorsmod.Register(ModuleName, CodeStreamAlreadyExists, "stream exists")
	ErrInvalidData          = errorsmod.Register(ModuleName, CodeInvalidData, "invalid data")
	ErrMissingData          = errorsmod.Register(ModuleName, CodeMissingData, "missing data")
	ErrContentTooLarge      = errorsmod.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrStreamNotCancellable = errorsmod.Register(ModuleName, CodeNotCancellable, "stream not cancellable")

	ErrInsufficientDeposit = errorsmod.Register(ModuleName, CodeInsufficientDeposit, "insufficient deposit")
)
