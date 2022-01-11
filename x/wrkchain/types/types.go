package types

import (
	undtypes "github.com/unification-com/mainchain/types"
)

const (
	// WRKChain fees, in nano FUND
	RegFee    = 1000000000000                // 1000 FUND - used in init genesis
	RecordFee = 1000000000                   // 1 FUND - used in init genesis
	FeeDenom  = undtypes.DefaultDenomination // used in init genesis

	DefaultStartingWrkChainID      uint64 = 1 // used in init genesis
	MaxBlockSubmissionsKeepInState        = 20000
)

func NewWrkchain(wrkchainId uint64, moniker, name, genesis, wcType string,
	lastBlock, numBlocks, regTime uint64, owner string) (WrkChain, error) {
	wc := WrkChain{
		WrkchainId: wrkchainId,
		Moniker:    moniker,
		Name:       name,
		Genesis:    genesis,
		Type:       wcType,
		Lastblock:  lastBlock,
		NumBlocks:  numBlocks,
		RegTime:    regTime,
		Owner:      owner,
	}

	return wc, nil
}

func NewWrkchainBlock(wrkchainId, height uint64, blockHash, parentHash, h1, h2, h3 string,
	subTime uint64, owner string) (WrkChainBlock, error) {
	b := WrkChainBlock{
		WrkchainId: wrkchainId,
		Height:     height,
		Blockhash:  blockHash,
		Parenthash: parentHash,
		Hash1:      h1,
		Hash2:      h2,
		Hash3:      h3,
		SubTime:    subTime,
		Owner:      owner,
	}
	return b, nil
}