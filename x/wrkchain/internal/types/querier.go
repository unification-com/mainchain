package types

import "strings"

// QueryResNames Queries Result Payload for a WRKChain Block Hashes query
type QueryResWrkChainBlockHashes []WrkChainBlock

// implement fmt.Stringer
func (h QueryResWrkChainBlockHashes) String() (out string) {
	for _, val := range h {
		out += val.String() + "\n"
	}
	return strings.TrimSpace(out)
}
