package stream

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/unification-com/mainchain/x/stream/simulation"
	"github.com/unification-com/mainchain/x/stream/types"
)

// avoid unused import issue
var (
	_ = baseapp.Paramspace
)

// GenerateGenesisState creates a randomized GenState of the feegrant module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for feegrant module's types
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns all the feegrant module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(
		simState.AppParams, simState.Cdc,
		am.keeper, am.bankKeeper, am.accountKeeper,
	)
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return simulation.ProposalMsgs()
}
