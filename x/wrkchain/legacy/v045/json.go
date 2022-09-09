package v045

import (
	v040 "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

// MigrateJSON accepts exported v0.40 x/wrkchain genesis state and migrates it to
// v0.45 (0.43) x/wrkchain genesis state.
func MigrateJSON(oldWrkchainState *v040.GenesisState) *types.GenesisState {
	newWrkchains := make(types.WrkChainExports, len(oldWrkchainState.RegisteredWrkchains))

	for i, oldWrkchain := range oldWrkchainState.RegisteredWrkchains {
		lowestHeight := uint64(0)
		newWrkchainBlocks := make(types.WrkChainBlockGenesisExports, len(oldWrkchain.Blocks))
		for j, oldWrkchainBlock := range oldWrkchain.Blocks {
			if lowestHeight == 0 || oldWrkchainBlock.He < lowestHeight {
				lowestHeight = oldWrkchainBlock.He
			}
			newWrkchainBlocks[j] = types.WrkChainBlockGenesisExport{
				He: oldWrkchainBlock.He,
				Bh: oldWrkchainBlock.Bh,
				Ph: oldWrkchainBlock.Ph,
				H1: oldWrkchainBlock.H1,
				H2: oldWrkchainBlock.H2,
				H3: oldWrkchainBlock.H3,
				St: oldWrkchainBlock.St,
			}
		}

		newWrkchains[i] = types.WrkChainExport{
			Wrkchain: types.WrkChain{
				WrkchainId:   oldWrkchain.Wrkchain.WrkchainId,
				Moniker:      oldWrkchain.Wrkchain.Moniker,
				Name:         oldWrkchain.Wrkchain.Name,
				Genesis:      oldWrkchain.Wrkchain.Genesis,
				Type:         oldWrkchain.Wrkchain.Type,
				Lastblock:    oldWrkchain.Wrkchain.Lastblock,
				NumBlocks:    uint64(len(oldWrkchain.Blocks)),
				LowestHeight: lowestHeight,
				RegTime:      oldWrkchain.Wrkchain.RegTime,
				Owner:        oldWrkchain.Wrkchain.Owner,
			},
			Blocks:       newWrkchainBlocks,
			InStateLimit: types.DefaultStorageLimit,
		}
	}

	return &types.GenesisState{
		Params: types.Params{
			FeeRegister:         oldWrkchainState.Params.FeeRegister,
			FeeRecord:           oldWrkchainState.Params.FeeRecord,
			FeePurchaseStorage:  types.PurchaseStorageFee,
			Denom:               oldWrkchainState.Params.Denom,
			DefaultStorageLimit: types.DefaultStorageLimit,
			MaxStorageLimit:     types.DefaultMaxStorageLimit,
		},
		StartingWrkchainId:  oldWrkchainState.StartingWrkchainId,
		RegisteredWrkchains: newWrkchains,
	}
}
