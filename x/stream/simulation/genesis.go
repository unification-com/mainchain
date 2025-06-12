package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	mathmod "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/unification-com/mainchain/x/stream/types"
)

// Simulation parameter constants
const (
	ValidatorFee = "validator_fee"
)

// GenValidatorFee randomized ValidatorFee
func GenValidatorFee(r *rand.Rand) mathmod.LegacyDec {
	// 0 to 50%
	return mathmod.LegacyNewDecWithPrec(int64(r.Intn(24)), 2)
}

// RandomizedGenState generates a random GenesisState for the stream module.
func RandomizedGenState(simState *module.SimulationState) {
	var streams []types.StreamExport

	var validatorFee mathmod.LegacyDec
	simState.AppParams.GetOrGenerate(
		ValidatorFee, &validatorFee, simState.Rand,
		func(r *rand.Rand) { validatorFee = GenValidatorFee(r) },
	)

	params := types.NewParams(validatorFee)

	streamGenState := types.NewGenesisState(streams, params)
	bz, err := json.MarshalIndent(&streamGenState, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Selected randomly generated stream parameters:\n%s\n", bz)

	simState.GenState[types.ModuleName] = bz
}
