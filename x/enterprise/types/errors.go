package types

import (
	errorsmod "cosmossdk.io/errors"
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
	CodeInvalidWhitelistAction        = 111
	CodeAlreadyWhitelisted            = 112
	CodeAddressNotWhitelisted         = 113
	CodeNotAuthorisedToRaisePO        = 114
)

var (
	ErrInvalidGenesis                = errorsmod.Register(ModuleName, CodeInvalidGenesis, "invalid genesis")
	ErrPurchaseOrderDoesNotExist     = errorsmod.Register(ModuleName, CodePurchaseOrderDoesNotExist, "purchase order does not exist")
	ErrPurchaseOrderAlreadyProcessed = errorsmod.Register(ModuleName, CodePurchaseOrderAlreadyProcessed, "purchase order already processed")
	ErrInvalidDecision               = errorsmod.Register(ModuleName, CodeInvalidDecision, "invalid decision")
	ErrInvalidDenomination           = errorsmod.Register(ModuleName, CodeInvalidDenomination, "invalid denomination")
	ErrInvalidStatus                 = errorsmod.Register(ModuleName, CodeInvalidStatus, "invalid status")
	ErrPurchaseOrderNotRaised        = errorsmod.Register(ModuleName, CodePurchaseOrderNotRaised, "purchase order not raised")
	ErrSignerAlreadyMadeDecision     = errorsmod.Register(ModuleName, CodeSignerAlreadyMadeDecision, "signer already made decision")
	ErrMissingData                   = errorsmod.Register(ModuleName, CodeMissingData, "missing data")
	ErrInvalidData                   = errorsmod.Register(ModuleName, CodeDataInvalid, "invalid data")
	ErrInvalidWhitelistAction        = errorsmod.Register(ModuleName, CodeInvalidWhitelistAction, "invalid whitelist action")
	ErrAlreadyWhitelisted            = errorsmod.Register(ModuleName, CodeAlreadyWhitelisted, "address already whitelisted")
	ErrAddressNotWhitelisted         = errorsmod.Register(ModuleName, CodeAddressNotWhitelisted, "address not whitelisted")
	ErrNotAuthorisedToRaisePO        = errorsmod.Register(ModuleName, CodeNotAuthorisedToRaisePO, "address not authorised to raise purchase orders")
)
