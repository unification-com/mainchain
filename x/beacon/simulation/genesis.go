package simulation

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

// RandomizedGenState generates a random GenesisState for beacon module
func RandomizedGenState(simState *module.SimulationState) {
	startingBeaconID := uint64(simState.Rand.Intn(100))

	// NOTE: for simulation, we're using sdk.DefaultBondDenom ("stake"), since "stake" is hard-coded
	// into the SDK's module simulation functions
	beaconGenesis := types.NewGenesisState(
		types.NewParams(10000000, 1000000, sdk.DefaultBondDenom),
		startingBeaconID,
	)

	fmt.Printf("Selected randomly generated beacon parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, beaconGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(beaconGenesis)
}
