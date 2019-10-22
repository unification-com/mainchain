package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeWrkChainDoesNotExist sdk.CodeType = 101
)

func ErrWrkChainDoesNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeWrkChainDoesNotExist, "WrkChain does not exist")
}
