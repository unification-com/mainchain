package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // defined in keys.go file

// MsgRegisterWrkChain defines a RegisterWrkChain message
type MsgRegisterWrkChain struct {
	WrkChainID   string         `json:"id"`
	WrkChainName string         `json:"name"`
	GenesisHash  string         `json:"genesis"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewMsgRegisterWrkChain is a constructor function for MsgRegisterWrkChain
func NewMsgRegisterWrkChain(wrkchainId string, genesisHash string, wrkchainName string, owner sdk.AccAddress) MsgRegisterWrkChain {
	return MsgRegisterWrkChain{
		WrkChainID:  wrkchainId,
		WrkChainName: wrkchainName,
		GenesisHash: genesisHash,
		Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgRegisterWrkChain) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRegisterWrkChain) Type() string { return "register" }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterWrkChain) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.WrkChainID) == 0 || len(msg.GenesisHash) == 0 || len(msg.WrkChainName) == 0 {
		return sdk.ErrUnknownRequest("WrkChainID, Genesis Hash and WRKChain Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterWrkChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterWrkChain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
