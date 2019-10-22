package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Wrkchain is a struct that contains all the metadata of a registered WRKChain
type WrkChain struct {
	WrkChainID   string         `json:"id"`
	WrkChainName string         `json:"name"`
	GenesisHash  string         `json:"genesis"`
	LastBlock    uint64         `json:"lastblock"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewWrkChain returns a new WrkChain struct
func NewWrkChain() WrkChain {
	return WrkChain{}
}

// implement fmt.Stringer
func (w WrkChain) String() string {
	return strings.TrimSpace(fmt.Sprintf(`WrkChainID: %s
WrkChainName: %s
GenesisHash: %s
LastBlock: %s
Owner: %s`, w.WrkChainID, w.WrkChainName, w.GenesisHash, w.LastBlock, w.Owner))
}
