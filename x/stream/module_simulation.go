package stream

//// avoid unused import issue
//var (
//	_ = sample.AccAddress
//	_ = streamsimulation.FindAccount
//	_ = simulation.MsgEntryKind
//	_ = baseapp.Paramspace
//	_ = rand.Rand{}
//)
//
//const (
//)
//
//// GenerateGenesisState creates a randomized GenState of the module.
//func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
//	accs := make([]string, len(simState.Accounts))
//	for i, acc := range simState.Accounts {
//		accs[i] = acc.Address.String()
//	}
//	streamGenesis := types.GenesisState{
//		Params: types.DefaultParams(),
//	}
//	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&streamGenesis)
//}
//
//// RegisterStoreDecoder registers a decoder.
//func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}
//
//// ProposalContents doesn't return any content functions for governance proposals.
//func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
//	return nil
//}
//
//// WeightedOperations returns the all the gov module operations with their respective weights.
//func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
//	operations := make([]simtypes.WeightedOperation, 0)
//
//
//	return operations
//}
//
//// ProposalMsgs returns msgs used for governance proposals for simulations.
//func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
//	return []simtypes.WeightedProposalMsg{
//	}
//}
