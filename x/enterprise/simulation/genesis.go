package simulation

import (
	//"fmt"
	//"math/rand"

	"encoding/json"
	"fmt"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	EnterpriseSourceAddress = "ent_signers"
)

// RandomizedGenState generates a random GenesisState for enterprise
func RandomizedGenState(simState *module.SimulationState) {
	var entAddress sdk.AccAddress
	simState.AppParams.GetOrGenerate(
		simState.Cdc, EnterpriseSourceAddress, &entAddress, simState.Rand,
		func(r *rand.Rand) { entAddress = simState.Accounts[0].Address },
	)

	// NOTE: for simulation, we're using sdk.DefaultBondDenom ("stake"), since "stake" is hard-coded
	// into the SDK's module simulation functions
	entGenesis := types.NewGenesisState(
		types.NewParams(sdk.DefaultBondDenom, 1, 3600, entAddress.String()),
		uint64(1),
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
	)

	bz, err := json.MarshalIndent(&entGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated enterprise parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(entGenesis)
}
