package v040

import (
	v038 "github.com/unification-com/mainchain/x/wrkchain/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/wrkchain/types"
)

func Migrate(oldWrkchainState v038.GenesisState) *v040.GenesisState {
	newWrkchains := make(v040.WrkChainExports, len(oldWrkchainState.WrkChains))

	for i, oldWrkchain := range oldWrkchainState.WrkChains {
		newWrkchainBlocks := make(v040.WrkChainBlockGenesisExports, len(oldWrkchain.WrkChainBlocks))
		for j, oldWrkchainBlock := range oldWrkchain.WrkChainBlocks {
			newWrkchainBlocks[j] = v040.WrkChainBlockGenesisExport{
				He: oldWrkchainBlock.Height,
				Bh: oldWrkchainBlock.BlockHash,
				Ph: oldWrkchainBlock.ParentHash,
				H1: oldWrkchainBlock.Hash1,
				H2: oldWrkchainBlock.Hash2,
				H3: oldWrkchainBlock.Hash3,
				St: uint64(oldWrkchainBlock.SubmitTime),
			}
		}

		newWrkchains[i] = v040.WrkChainExport{
			Wrkchain: v040.WrkChain{
				WrkchainId: oldWrkchain.WrkChain.WrkChainID,
				Moniker:    oldWrkchain.WrkChain.Moniker,
				Name:       oldWrkchain.WrkChain.Name,
				Genesis:    oldWrkchain.WrkChain.GenesisHash,
				Type:       oldWrkchain.WrkChain.BaseType,
				Lastblock:  oldWrkchain.WrkChain.LastBlock,
				NumBlocks:  oldWrkchain.WrkChain.NumberBlocks,
				RegTime:    uint64(oldWrkchain.WrkChain.RegisterTime),
				Owner:      oldWrkchain.WrkChain.Owner.String(),
			},
			Blocks: newWrkchainBlocks,
		}
	}
	
	return &v040.GenesisState{
		Params: v040.Params{
			FeeRegister: oldWrkchainState.Params.FeeRegister,
			FeeRecord:   oldWrkchainState.Params.FeeRecord,
			Denom:       oldWrkchainState.Params.Denom,
		},
		StartingWrkchainId:  oldWrkchainState.StartingWrkChainID,
		RegisteredWrkchains: newWrkchains,
	}
}
