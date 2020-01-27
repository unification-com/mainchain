package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodeInvalidGenesis                = 101
	CodePurchaseOrderDoesNotExist     = 102
	CodePurchaseOrderAlreadyProcessed = 103
	CodeInvalidDecision               = 104
	CodeInvalidDenomination           = 105
	CodeInvalidStatus                 = 106
	CodePurchaseOrderNotRaised        = 107
	CodeSignerAlreadyMadeDecision     = 108
	CodeMissingData                   = 109
	CodeDataInvalid                   = 110
)

var (
	ErrInvalidGenesis                = sdkerrors.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrPurchaseOrderDoesNotExist     = sdkerrors.Register(ModuleName, CodePurchaseOrderDoesNotExist, "purchase order does not exist")
	ErrPurchaseOrderAlreadyProcessed = sdkerrors.Register(ModuleName, CodePurchaseOrderAlreadyProcessed, "purchase order already processed")
	ErrInvalidDecision               = sdkerrors.Register(ModuleName, CodeInvalidDecision, "invalid decision")
	ErrInvalidDenomination           = sdkerrors.Register(ModuleName, CodeInvalidDenomination, "invalid denomination")
	ErrInvalidStatus                 = sdkerrors.Register(ModuleName, CodeInvalidStatus, "invalid status")
	ErrPurchaseOrderNotRaised        = sdkerrors.Register(ModuleName, CodePurchaseOrderNotRaised, "purchase order not raised")
	ErrSignerAlreadyMadeDecision     = sdkerrors.Register(ModuleName, CodeSignerAlreadyMadeDecision, "signer already made decision")
	ErrMissingData                   = sdkerrors.Register(ModuleName, CodeMissingData, "missing data")
	ErrInvalidData                   = sdkerrors.Register(ModuleName, CodeDataInvalid, "invalid data")
)
