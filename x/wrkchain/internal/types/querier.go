package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryResWrkChainBlockHashes Queries Result Payload for a WRKChain Block Hashes query
type QueryResWrkChainBlockHashes []WrkChainBlock

// implement fmt.Stringer
func (h QueryResWrkChainBlockHashes) String() (out string) {
	for _, val := range h {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
// QueryResWrkChains Queries wrkchains
type QueryResWrkChains []WrkChain

// implement fmt.Stringer
func (wc QueryResWrkChains) String() (out string) {
	for _, val := range wc {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}


// QueryWrkChainParams Params for query 'custom/wrkchain/registered'
type QueryWrkChainParams struct {
	Page    int
	Limit   int
	Moniker string
	Owner   sdk.AccAddress
}

// NewQueryWrkChainParams creates a new instance of QueryWrkChainParams
func NewQueryWrkChainParams(page, limit int, monkker string, owner sdk.AccAddress) QueryWrkChainParams {
	return QueryWrkChainParams{
		Page:    page,
		Limit:   limit,
		Moniker: monkker,
		Owner:   owner,
	}
}
