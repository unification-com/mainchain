package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	RegisterAction = "register"
	RecordAction = "record"
)
// --- Register a WRKChain Msg ---

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
func (msg MsgRegisterWrkChain) Type() string { return RegisterAction }

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

// --- Record a WRKChain Block hash Msg ---

// MsgRecordWrkChainBlock defines a RecordWrkChainBlock message
type MsgRecordWrkChainBlock struct {
	WrkChainID   string         `json:"id"`
	Height       uint64         `json:"height"`
	BlockHash    string         `json:"blockhash"`
	ParentHash   string         `json:"parenthash"`
	Hash1        string         `json:"hash1"`
	Hash2        string         `json:"hash2"`
	Hash3        string         `json:"hash3"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewMsgRecordWrkChainBlock is a constructor function for MsgRecordWrkChainBlock
func NewMsgRecordWrkChainBlock(
	wrkchainId string,
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
	owner sdk.AccAddress) MsgRecordWrkChainBlock {

	return MsgRecordWrkChainBlock{
	    WrkChainID:  wrkchainId,
		Height: height,
		BlockHash: blockHash,
		ParentHash: parentHash,
		Hash1: hash1,
		Hash2: hash2,
		Hash3: hash3,
	    Owner: owner,
	}
}

// Route should return the name of the module
func (msg MsgRecordWrkChainBlock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecordWrkChainBlock) Type() string { return RecordAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecordWrkChainBlock) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.WrkChainID) == 0 || len(msg.BlockHash) == 0 {
		return sdk.ErrUnknownRequest("WrkChainID and BlockHash cannot be empty")
	}
	if msg.Height == 0 {
		return sdk.ErrUnknownRequest("Height cannot be zero")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRecordWrkChainBlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRecordWrkChainBlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}