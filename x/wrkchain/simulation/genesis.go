package simulation

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

// RandomizedGenState generates a random GenesisState for wrkchain
func RandomizedGenState(simState *module.SimulationState) {
	startingWrkChainID := uint64(simState.Rand.Intn(100))

	wrkchainGenesis := types.NewGenesisState(
		startingWrkChainID,
	)

	fmt.Printf("Selected randomly generated wrkchain parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, wrkchainGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(wrkchainGenesis)
}
