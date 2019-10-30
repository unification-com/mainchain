package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = ModuleName

	CodeWrkChainDoesNotExist      = 101
	CodeWrkChainAlreadyRegistered = 102

	CodeWrkChainInsufficientFee = 201
	CodeWrkChainTooMuchFee      = 202

	CodeWrkChainFeePayerNotOwner = 301
)

var (
	ErrWrkChainDoesNotExist      = sdkerrors.Register(DefaultCodespace, CodeWrkChainDoesNotExist, "WrkChain does not exist")
	ErrWrkChainAlreadyRegistered = sdkerrors.Register(DefaultCodespace, CodeWrkChainAlreadyRegistered, "WrkChain already registered")

	ErrInsufficientWrkChainFee = sdkerrors.Register(DefaultCodespace, CodeWrkChainInsufficientFee, "insufficient fee to pay for WrkChain tx")
	ErrTooMuchWrkChainFee      = sdkerrors.Register(DefaultCodespace, CodeWrkChainTooMuchFee, "too much fee sent to pay for WrkChain tx")

	ErrFeePayerNotOwner = sdkerrors.Register(DefaultCodespace, CodeWrkChainFeePayerNotOwner, "fee payer is not WRKChain owner")
)
