package v040

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	v038beacon "github.com/unification-com/mainchain/x/beacon/legacy/v038"
	v040beacon "github.com/unification-com/mainchain/x/beacon/legacy/v040"
	v038enterprise "github.com/unification-com/mainchain/x/enterprise/legacy/v038"
	v040enterprise "github.com/unification-com/mainchain/x/enterprise/legacy/v040"
	v038wrkchain "github.com/unification-com/mainchain/x/wrkchain/legacy/v038"
	v040wrkchain "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
)

// Migrate migrates exported state from v0.39 to a v0.40 genesis state.
func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	v038Codec := codec.NewLegacyAmino()
	v038beacon.RegisterLegacyAminoCodec(v038Codec)
	v038wrkchain.RegisterLegacyAminoCodec(v038Codec)
	v038enterprise.RegisterLegacyAminoCodec(v038Codec)
	v040Codec := clientCtx.JSONMarshaler

	if appState[v038beacon.ModuleName] != nil {
		var beaconGenState v038beacon.GenesisState
		v038Codec.MustUnmarshalJSON(appState[v038beacon.ModuleName], &beaconGenState)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040beacon.ModuleName] = v040Codec.MustMarshalJSON(v040beacon.Migrate(beaconGenState))
	}

	if appState[v038wrkchain.ModuleName] != nil {
		var wrkchainGenState v038wrkchain.GenesisState
		v038Codec.MustUnmarshalJSON(appState[v038wrkchain.ModuleName], &wrkchainGenState)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040wrkchain.ModuleName] = v040Codec.MustMarshalJSON(v040wrkchain.Migrate(wrkchainGenState))
	}

	if appState[v038enterprise.ModuleName] != nil {
		var enterpriseGenState v038enterprise.GenesisState
		v038Codec.MustUnmarshalJSON(appState[v038enterprise.ModuleName], &enterpriseGenState)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v040enterprise.ModuleName] = v040Codec.MustMarshalJSON(v040enterprise.Migrate(enterpriseGenState))
	}

	return appState
}
