package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	RouterKey = ModuleName // defined in keys.go file

	RegisterAction        = "register_wrkchain"
	RecordAction          = "record_wrkchain_hash"
	PurchaseStorageAction = "purchase_wrkchain_storage"
)

var (
	_ sdk.Msg = &MsgRegisterWrkChain{}
	_ sdk.Msg = &MsgPurchaseWrkChainStateStorage{}
	_ sdk.Msg = &MsgRecordWrkChainBlock{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// --- Register a WRKChain Msg ---

// NewMsgRegisterWrkChain is a constructor function for MsgRegisterWrkChain
func NewMsgRegisterWrkChain(moniker string, genesisHash string, wrkchainName string, baseType string, owner sdk.AccAddress) *MsgRegisterWrkChain {
	return &MsgRegisterWrkChain{
		Moniker:     moniker,
		Name:        wrkchainName,
		GenesisHash: genesisHash,
		BaseType:    baseType,
		Owner:       owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgRegisterWrkChain) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRegisterWrkChain) Type() string { return RegisterAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterWrkChain) ValidateBasic() error {
	ownerAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	if ownerAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner)
	}
	if len(msg.Moniker) == 0 {
		return sdkerrors.Wrap(ErrMissingData, "Moniker cannot be empty")
	}

	if len(msg.Name) > 128 {
		return sdkerrors.Wrap(ErrContentTooLarge, "name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return sdkerrors.Wrap(ErrContentTooLarge, "moniker too big. 64 character limit")
	}

	if len(msg.GenesisHash) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "genesis hash too big. 66 character limit")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRegisterWrkChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterWrkChain) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Record a WRKChain Block hash Msg ---

// NewMsgRecordWrkChainBlock is a constructor function for MsgRecordWrkChainBlock
func NewMsgRecordWrkChainBlock(
	wrkchainId uint64,
	height uint64,
	blockHash string,
	parentHash string,
	hash1 string,
	hash2 string,
	hash3 string,
	owner sdk.AccAddress) *MsgRecordWrkChainBlock {

	return &MsgRecordWrkChainBlock{
		WrkchainId: wrkchainId,
		Height:     height,
		BlockHash:  blockHash,
		ParentHash: parentHash,
		Hash1:      hash1,
		Hash2:      hash2,
		Hash3:      hash3,
		Owner:      owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgRecordWrkChainBlock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecordWrkChainBlock) Type() string { return RecordAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecordWrkChainBlock) ValidateBasic() error {
	ownerAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}

	if ownerAddr.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner)
	}
	if msg.WrkchainId == 0 {
		return sdkerrors.Wrap(ErrInvalidData, "ID must be greater than zero")
	}
	if len(msg.BlockHash) == 0 {
		return sdkerrors.Wrap(ErrMissingData, "BlockHash cannot be empty")
	}
	if msg.Height == 0 {
		return sdkerrors.Wrap(ErrMissingData, "Height cannot be zero")
	}
	if len(msg.BlockHash) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "block hash too big. 66 character limit")
	}
	if len(msg.ParentHash) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "parent hash too big. 66 character limit")
	}
	if len(msg.Hash1) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "hash1 too big. 66 character limit")
	}
	if len(msg.Hash2) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "hash2 too big. 66 character limit")
	}
	if len(msg.Hash3) > 66 {
		return sdkerrors.Wrap(ErrContentTooLarge, "hash3 too big. 66 character limit")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRecordWrkChainBlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgRecordWrkChainBlock) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Purchase state storage Msg ---

// NewMsgRecordBeaconTimestamp is a constructor function for MsgRecordBeaconTimestamp
func NewMsgPurchaseWrkChainStateStorage(
	wrkchainId uint64,
	number uint64,
	owner sdk.AccAddress) *MsgPurchaseWrkChainStateStorage {

	return &MsgPurchaseWrkChainStateStorage{
		WrkchainId: wrkchainId,
		Number:     number,
		Owner:      owner.String(),
	}
}

// Route should return the name of the module
func (msg MsgPurchaseWrkChainStateStorage) Route() string { return RouterKey }

// Type should return the action
func (msg MsgPurchaseWrkChainStateStorage) Type() string { return PurchaseStorageAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgPurchaseWrkChainStateStorage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid owner address (%s)", err)
	}
	if msg.WrkchainId == 0 {
		return sdkerrors.Wrap(ErrMissingData, "id must be greater than zero")
	}
	if msg.Number == 0 {
		return sdkerrors.Wrap(ErrMissingData, "number cannot be zero")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPurchaseWrkChainStateStorage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners defines whose signature is required
func (msg MsgPurchaseWrkChainStateStorage) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// --- Modify Params Msg Type ---

// GetSignBytes returns the raw bytes for a MsgUpdateParams message that
// the expected signer needs to sign.
func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check on the provided data.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}
