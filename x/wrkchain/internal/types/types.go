package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain-cosmos/types"
)

const (
	// WRKChain fees, in nano UND
	RegFee    = 1000000000000                // 1000 UND - used in init genesis
	RecordFee = 1000000000                   // 1 UND - used in init genesis
	FeeDenom  = undtypes.DefaultDenomination // used in init genesis

	DefaultStartingWrkChainID uint64 = 1 // used in init genesis
)

// WrkChains is an array of WrkChain
type WrkChains []WrkChain

// String implements stringer interface
func (w WrkChains) String() string {
	out := "ID - [Moniker] 'Name' (Genesis) {LastBlock} Owner\n"
	for _, wc := range w {
		out += fmt.Sprintf("%d - [%s] '%s' (%s) {%d} %s\n",
			wc.WrkChainID, wc.Moniker,
			wc.Name, wc.GenesisHash, wc.LastBlock, wc.Owner)
	}
	return strings.TrimSpace(out)
}

// Wrkchain is a struct that contains all the metadata of a registered WRKChain
type WrkChain struct {
	WrkChainID   uint64         `json:"wrkchain_id"`
	Moniker      string         `json:"moniker"`
	Name         string         `json:"name"`
	GenesisHash  string         `json:"genesis"`
	BaseType     string         `json:"type"`
	LastBlock    uint64         `json:"lastblock"`
	RegisterTime int64          `json:"reg_time"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewWrkChain returns a new WrkChain struct
func NewWrkChain() WrkChain {
	return WrkChain{}
}

// implement fmt.Stringer
func (w WrkChain) String() string {
	return strings.TrimSpace(fmt.Sprintf(`WRKChainID: %d
Moniker: %s
Name: %s
GenesisHash: %s
BaseType: %s
LastBlock: %d
RegisterTime: %d
Owner: %s`, w.WrkChainID, w.Moniker, w.Name, w.GenesisHash, w.BaseType, w.LastBlock, w.RegisterTime, w.Owner))
}

// WrkChainBlocks is an array of WrkChainBlock
type WrkChainBlocks []WrkChainBlock

// String implements stringer interface
func (wcb WrkChainBlocks) String() string {
	out := "ID - [Height] 'BlockHash' (ParentHash) {Hash1} <Hash2> `Hash3` Owner\n"
	for _, b := range wcb {
		out += fmt.Sprintf("%d - [%d] '%s' (%s) {%s} <%s> `%s` %s\n",
			b.WrkChainID, b.Height, b.BlockHash, b.ParentHash,
			b.Hash1, b.Hash2, b.Hash3, b.Owner)
	}
	return strings.TrimSpace(out)
}

// WrkChainBlock is a struct that contains a wrkchain's recorded block
type WrkChainBlock struct {
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

// NewWrkChainBlock returns a new WrkChainBlock struct
func NewWrkChainBlock() WrkChainBlock {
	return WrkChainBlock{}
}

// implement fmt.Stringer
func (w WrkChainBlock) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Moniker: %d
Height: %d
BlockHash: %s
ParentHash: %s
Hash1: %s
Hash2: %s
Hash3: %s
SubmitTime: %d
Owner: %s`, w.WrkChainID, w.Height, w.BlockHash, w.ParentHash, w.Hash1, w.Hash2, w.Hash3, w.SubmitTime, w.Owner))
}
