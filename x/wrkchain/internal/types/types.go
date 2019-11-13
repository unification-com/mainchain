package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

const (
	// WRKChain fees, in nano UND
	RegFee    = 1000000000000 // 1000 UND
	RecordFee = 1000000000    // 1 UND
	FeeDenom  = "nund"

	DefaultStartingWrkChainID uint64 = 1 // used in init genesis
)

var (
	// WRKChain fees in sdk.Coin (denom=und) format. Exported in alias.go
	FeesBaseDenomination         = sdk.NewInt64Coin(FeeDenom, 0)
	FeesWrkChainRegistrationCoin = sdk.NewInt64Coin(FeeDenom, RegFee)
	FeesWrkChainRecordHashCoin   = sdk.NewInt64Coin(FeeDenom, RecordFee)

	// WRKChain Fees in sdk.Coins[]. Exported in alias.go
	FeesWrkChainRegistration = sdk.Coins{FeesWrkChainRegistrationCoin}
	FeesWrkChainRecordHash   = sdk.Coins{FeesWrkChainRecordHashCoin}
)

// Wrkchain is a struct that contains all the metadata of a registered WRKChain
type WrkChain struct {
	WrkChainID  uint64         `json:"wrkchain_id"`
	Moniker     string         `json:"moniker"`
	Name        string         `json:"name"`
	GenesisHash string         `json:"genesis"`
	LastBlock   uint64         `json:"lastblock"`
	Owner       sdk.AccAddress `json:"owner"`
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
LastBlock: %d
Owner: %s`, w.WrkChainID, w.Moniker, w.Name, w.GenesisHash, w.LastBlock, w.Owner))
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
	SubmitTime uint64         `json:"time"`
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
