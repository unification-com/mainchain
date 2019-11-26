package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidGenesis                = 101
	CodePurchaseOrderDoesNotExist     = 102
	CodePurchaseOrderAlreadyProcessed = 103
	CodeInvalidDecision               = 104
	CodeInvalidDenomination           = 105
	CodeInvalidStatus                 = 106
	CodePurchaseOrderNotRaised        = 107
)

// ErrInvalidGenesis error for an invalid enterprise GenesisState
func ErrInvalidGenesis(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidGenesis, msg)
}

func ErrPurchaseOrderDoesNotExist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePurchaseOrderDoesNotExist, msg)
}

func ErrPurchaseOrderAlreadyProcessed(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePurchaseOrderAlreadyProcessed, msg)
}

func ErrInvalidDecision(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDecision, msg)
}

func ErrInvalidDenomination(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDenomination, msg)
}

func ErrPurchaseOrderNotRaised(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePurchaseOrderNotRaised, msg)
}
