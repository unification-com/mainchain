package types

const MaxBlockSubmissionsKeepInState = 20000

type WrkChainExports []WrkChainExport
type WrkChainBlockGenesisExports []WrkChainBlockGenesisExport

// WrkChainBlockLegacy is only used to support old style hash output for the legacy REST endpoint
type WrkChainBlockLegacy struct {
	WrkChainID uint64 `json:"wrkchain_id"`
	Height     uint64 `json:"height"`
	BlockHash  string `json:"blockhash"`
	ParentHash string `json:"parenthash"`
	Hash1      string `json:"hash1"`
	Hash2      string `json:"hash2"`
	Hash3      string `json:"hash3"`
	SubmitTime uint64 `json:"sub_time"`
	Owner      string `json:"owner"`
}

func NewWrkchain(wrkchainId uint64, moniker, name, genesis, wcType string,
	lastBlock, numBlocksInState, regTime uint64, owner string) (WrkChain, error) {
	wc := WrkChain{
		WrkchainId: wrkchainId,
		Moniker:    moniker,
		Name:       name,
		Genesis:    genesis,
		BaseType:   wcType,
		Lastblock:  lastBlock,
		NumBlocks:  numBlocksInState,
		RegTime:    regTime,
		Owner:      owner,
	}

	return wc, nil
}

func NewWrkchainBlock(height uint64, blockHash, parentHash, h1, h2, h3 string,
	subTime uint64) (WrkChainBlock, error) {
	b := WrkChainBlock{
		Height:     height,
		Blockhash:  blockHash,
		Parenthash: parentHash,
		Hash1:      h1,
		Hash2:      h2,
		Hash3:      h3,
		SubTime:    subTime,
	}
	return b, nil
}
