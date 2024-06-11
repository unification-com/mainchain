package simulation

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/unification-com/mainchain/x/stream/types"
	"math/rand"
)

// Simulation parameter constants
const (
	ValidatorFee = "validator_fee"
)

// GenInflationRateChange randomized InflationRateChange
func GenValidatorFee(r *rand.Rand) math.LegacyDec {
	// 0 to 50%
	return sdk.NewDecWithPrec(int64(r.Intn(24)), 2)
}

// RandomizedGenState generates a random GenesisState for the stream module.
func RandomizedGenState(simState *module.SimulationState) {
	var streams []types.StreamExport

	var validatorFee sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, ValidatorFee, &validatorFee, simState.Rand,
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
