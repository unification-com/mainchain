package simulation

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"
	"strings"

	"github.com/unification-com/mainchain/x/enterprise/types"
)

// Simulation parameter constants
const (
	EnterpriseSignerAddress     = "ent_signers"
	EnterpriseDecisionTimeLimit = "decision_time_limit"
	EnterpriseMinAccepts        = "min_accepts"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// RandomizedGenState generates a random GenesisState for enterprise
func RandomizedGenState(simState *module.SimulationState) {
	var entAddress string
	var decisionLimit uint64
	var minAccepts uint64

	simState.AppParams.GetOrGenerate(
		simState.Cdc, EnterpriseMinAccepts, &minAccepts, simState.Rand,
		func(r *rand.Rand) {
			minAccepts = uint64(simtypes.RandIntBetween(r, 1, 3))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, EnterpriseSignerAddress, &entAddress, simState.Rand,
		func(r *rand.Rand) {
			entAddresses := make([]string, minAccepts)
			i := uint64(0)
			for i < minAccepts {
				randAcc, _ := simtypes.RandomAcc(r, simState.Accounts)
				for contains(entAddresses, randAcc.Address.String()) {
					randAcc, _ = simtypes.RandomAcc(r, simState.Accounts)
				}
				entAddresses[i] = randAcc.Address.String()
				i += 1
			}
			entAddress = strings.Join(entAddresses, ",")
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, EnterpriseDecisionTimeLimit, &decisionLimit, simState.Rand,
		func(r *rand.Rand) {
			decisionLimit = uint64(simtypes.RandIntBetween(r, 9000, 15000))
		},
	)

	// NOTE: for simulation, we're using sdk.DefaultBondDenom ("stake"), since "stake" is hard-coded
	// into the SDK's module simulation functions
	entGenesis := types.NewGenesisState(
		types.NewParams(sdk.DefaultBondDenom, minAccepts, decisionLimit, entAddress),
		uint64(1),
		sdk.NewInt64Coin(sdk.DefaultBondDenom, 0), nil, nil, nil,
	)

	bz, err := json.MarshalIndent(&entGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated enterprise parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(entGenesis)
}
