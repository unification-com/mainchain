package v038

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// ModuleName module name
	ModuleName = "wrkchain"

	RouterKey = ModuleName // defined in keys.go file

	RegisterAction = "register_beacon"
	RecordAction   = "record_beacon_timestamp"
)

type (
	WrkChains      []WrkChain
	WrkChainBlocks []WrkChainBlock

	Params struct {
		FeeRegister uint64 `json:"fee_register" yaml:"fee_register"` // Fee for registering a WRKChain
		FeeRecord   uint64 `json:"fee_record" yaml:"fee_record"`     // Fee for recording hashes for a WRKChain
		Denom       string `json:"denom" yaml:"denom"`               // Fee denomination
	}

	WrkChain struct {
		WrkChainID   uint64         `json:"wrkchain_id"`
		Moniker      string         `json:"moniker"`
		Name         string         `json:"name"`
		GenesisHash  string         `json:"genesis"`
		BaseType     string         `json:"type"`
		LastBlock    uint64         `json:"lastblock"`
		NumberBlocks uint64         `json:"num_blocks"`
		RegisterTime int64          `json:"reg_time"`
		Owner        sdk.AccAddress `json:"owner"`
	}

	WrkChainBlock struct {
		WrkChainID uint64         `json:"wrkchain_id"`
		Height     uint64         `json:"height"`
		BlockHash  string         `json:"blockhash"`
		ParentHash string         `json:"parenthash"`
		Hash1      string         `json:"hash1"`
		Hash2      string         `json:"hash2"`
		Hash3      string         `json:"hash3"`
		SubmitTime int64          `json:"sub_time"`
		Owner      sdk.AccAddress `json:"owner"`
	}

	GenesisState struct {
		Params             Params           `json:"params" yaml:"params"`                             // wrkchain params
		StartingWrkChainID uint64           `json:"starting_wrkchain_id" yaml:"starting_wrkchain_id"` // should be 1
		WrkChains          []WrkChainExport `json:"registered_wrkchains" yaml:"registered_wrkchains"`
	}

	WrkChainExport struct {
		WrkChain       WrkChain        `json:"wrkchain" yaml:"wrkchain"`
		WrkChainBlocks []WrkChainBlock `json:"blocks" yaml:"blocks"`
	}
)

// --- Register a WRKChain Msg ---

// MsgRegisterWrkChain defines a RegisterWrkChain message
type MsgRegisterWrkChain struct {
	Moniker      string         `json:"moniker"`
	WrkChainName string         `json:"name"`
	GenesisHash  string         `json:"genesis"`
	BaseType     string         `json:"type"`
	Owner        sdk.AccAddress `json:"owner"`
}

// Route should return the name of the module
func (msg MsgRegisterWrkChain) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRegisterWrkChain) Type() string { return RegisterAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterWrkChain) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Moniker) == 0 {
		return fmt.Errorf("Moniker cannot be empty")
	}

	if len(msg.WrkChainName) > 128 {
		return fmt.Errorf("name too big. 128 character limit")
	}

	if len(msg.Moniker) > 64 {
		return fmt.Errorf("moniker too big. 64 character limit")
	}

	if len(msg.GenesisHash) > 66 {
		return fmt.Errorf("genesis hash too big. 66 character limit")
	}

	return nil
}

// --- Record a WRKChain Block hash Msg ---

// MsgRecordWrkChainBlock defines a RecordWrkChainBlock message
type MsgRecordWrkChainBlock struct {
	WrkChainID uint64         `json:"wrkchain_id"`
	Height     uint64         `json:"height"`
	BlockHash  string         `json:"blockhash"`
	ParentHash string         `json:"parenthash"`
	Hash1      string         `json:"hash1"`
	Hash2      string         `json:"hash2"`
	Hash3      string         `json:"hash3"`
	Owner      sdk.AccAddress `json:"owner"`
}

// Route should return the name of the module
func (msg MsgRecordWrkChainBlock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRecordWrkChainBlock) Type() string { return RecordAction }

// ValidateBasic runs stateless checks on the message
func (msg MsgRecordWrkChainBlock) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if msg.WrkChainID == 0 {
		return fmt.Errorf("ID must be greater than zero")
	}
	if len(msg.BlockHash) == 0 {
		return fmt.Errorf("BlockHash cannot be empty")
	}
	if msg.Height == 0 {
		return fmt.Errorf("Height cannot be zero")
	}
	if len(msg.BlockHash) > 66 {
		return fmt.Errorf("block hash too big. 66 character limit")
	}
	if len(msg.ParentHash) > 66 {
		return fmt.Errorf("parent hash too big. 66 character limit")
	}
	if len(msg.Hash1) > 66 {
		return fmt.Errorf("hash1 too big. 66 character limit")
	}
	if len(msg.Hash2) > 66 {
		return fmt.Errorf("hash2 too big. 66 character limit")
	}
	if len(msg.Hash3) > 66 {
		return fmt.Errorf("hash3 too big. 66 character limit")
	}

	return nil
}

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgRegisterWrkChain{}, "wrkchain/RegisterWrkChain", nil)
	cdc.RegisterConcrete(MsgRecordWrkChainBlock{}, "wrkchain/RecordWrkChainBlock", nil)
}
