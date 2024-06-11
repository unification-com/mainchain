package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	ErrStreamDoesNotExist   = sdkerrors.Register(ModuleName, CodeStreamDoesNotExist, "stream does not exist")
	ErrStreamExists         = sdkerrors.Register(ModuleName, CodeStreamAlreadyExists, "stream exists")
	ErrInvalidData          = sdkerrors.Register(ModuleName, CodeInvalidData, "invalid data")
	ErrMissingData          = sdkerrors.Register(ModuleName, CodeMissingData, "missing data")
	ErrContentTooLarge      = sdkerrors.Register(ModuleName, CodeContentTooLarge, "msg content too large")
	ErrStreamNotCancellable = sdkerrors.Register(ModuleName, CodeNotCancellable, "stream not cancellable")

	ErrInsufficientDeposit = sdkerrors.Register(ModuleName, CodeInsufficientDeposit, "insufficient deposit")
)
