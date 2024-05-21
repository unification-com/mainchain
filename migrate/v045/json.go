package v045

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {

	if appState[ibctransfertypes.ModuleName] != nil {
		transferGenState := &ibctransfertypes.GenesisState{}
		clientCtx.Codec.MustUnmarshalJSON(appState[ibctransfertypes.ModuleName], transferGenState)

		substituteTraces := make([]ibctransfertypes.DenomTrace, len(transferGenState.DenomTraces))
		for i, dt := range transferGenState.DenomTraces {
			// replace all previous traces with the latest trace if validation passes
			// note most traces will have same value
			newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())

			if err := newTrace.Validate(); err != nil {
				substituteTraces[i] = dt
			} else {
				substituteTraces[i] = newTrace
			}
		}

		transferGenState.DenomTraces = substituteTraces

		// delete old genesis state
		delete(appState, ibctransfertypes.ModuleName)

		// set new ibc transfer genesis state
		appState[ibctransfertypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(transferGenState)
	}

	return appState
}
