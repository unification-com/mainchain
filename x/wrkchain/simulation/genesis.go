package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

// RandomizedGenState generates a random GenesisState for wrkchain
func RandomizedGenState(simState *module.SimulationState) {
	startingWrkChainID := uint64(simState.Rand.Intn(100))

	// NOTE: for simulation, we're using sdk.DefaultBondDenom ("stake"), since "stake" is hard-coded
	// into the SDK's module simulation functions
	wrkchainGenesis := types.NewGenesisState(
		types.NewParams(10000000, 10000000, sdk.DefaultBondDenom),
		startingWrkChainID,
	)

	fmt.Printf("Selected randomly generated wrkchain parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, wrkchainGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(wrkchainGenesis)
}
