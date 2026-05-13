package stream

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/simsx"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/unification-com/mainchain/x/stream/simulation"
	"github.com/unification-com/mainchain/x/stream/types"
)

// avoid unused import issue
var (
	_ = baseapp.Paramspace
)

// GenerateGenesisState creates a randomized GenState of the stream module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// RegisterStoreDecoder registers a decoder for stream module's types
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations is the legacy simsx back-compat path; ignored when
// WeightedOperationsX is also implemented (which it is, below).
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(
		simState.AppParams, simState.Cdc, simState.TxConfig,
		am.keeper, am.bankKeeper, am.accountKeeper,
	)
}

// WeightedOperationsX registers the stream module's sim msg factories with simsx.
func (am AppModule) WeightedOperationsX(weights simsx.WeightSource, reg simsx.Registry) {
	reg.Add(weights.Get(simulation.OpWeightMsgCreateStream, uint32(simulation.DefaultWeightMsgCreateStream)),
		simulation.MsgCreateStreamFactory(am.keeper))
	reg.Add(weights.Get(simulation.OpWeightMsgClaimStream, uint32(simulation.DefaultWeightMsgClaimStream)),
		simulation.MsgClaimStreamFactory(am.keeper))
	reg.Add(weights.Get(simulation.OpWeightMsgTopUpDeposit, uint32(simulation.DefaultWeightMsgTopUpDeposit)),
		simulation.MsgTopUpDepositFactory(am.keeper))
	reg.Add(weights.Get(simulation.OpWeightMsgUpdateFlowRate, uint32(simulation.DefaultWeightMsgUpdateFlowRate)),
		simulation.MsgUpdateFlowRateFactory(am.keeper))
	reg.Add(weights.Get(simulation.OpWeightMsgCancelStream, uint32(simulation.DefaultWeightMsgCancelStream)),
		simulation.MsgCancelStreamFactory(am.keeper))
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return simulation.ProposalMsgs()
}
