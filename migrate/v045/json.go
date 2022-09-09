package v045

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/authz"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfer "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibctypes "github.com/cosmos/ibc-go/v3/modules/core/types"

	v040beacon "github.com/unification-com/mainchain/x/beacon/legacy/v040"
	v045beacon "github.com/unification-com/mainchain/x/beacon/legacy/v045"
	v040ent "github.com/unification-com/mainchain/x/enterprise/legacy/v040"
	v045ent "github.com/unification-com/mainchain/x/enterprise/legacy/v045"
	v040wrkchain "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
	v045wrkchain "github.com/unification-com/mainchain/x/wrkchain/legacy/v045"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {

	// add authz default genesis
	if appState[authz.ModuleName] == nil {
		appState[authz.ModuleName] = clientCtx.Codec.MustMarshalJSON(authz.DefaultGenesisState())
	}

	// add capability default genesis
	if appState[capabilitytypes.ModuleName] == nil {
		appState[capabilitytypes.ModuleName] = clientCtx.Codec.MustMarshalJSON(capabilitytypes.DefaultGenesis())
	}

	// add feegrant default genesis
	if appState[feegrant.ModuleName] == nil {
		appState[feegrant.ModuleName] = clientCtx.Codec.MustMarshalJSON(feegrant.DefaultGenesisState())
	}

	// add ibc default genesis
	if appState[ibchost.ModuleName] == nil {
		appState[ibchost.ModuleName] = clientCtx.Codec.MustMarshalJSON(ibctypes.DefaultGenesisState())
	}

	// add ibc transfer default genesis
	if appState[ibctransfer.ModuleName] == nil {
		appState[ibctransfer.ModuleName] = clientCtx.Codec.MustMarshalJSON(ibctransfer.DefaultGenesisState())
	}

	// migrate BEACON
	if appState[v040beacon.ModuleName] != nil {
		var oldBeaconGenState v040beacon.GenesisState
		clientCtx.Codec.MustUnmarshalJSON(appState[v040beacon.ModuleName], &oldBeaconGenState)

		// delete deprecated x/beacon genesis state
		delete(appState, v040beacon.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v045beacon.ModuleName] = clientCtx.Codec.MustMarshalJSON(v045beacon.MigrateJSON(&oldBeaconGenState))
	}

	// migrate Wrkchain
	if appState[v040wrkchain.ModuleName] != nil {
		var oldWrkchainGenState v040wrkchain.GenesisState
		clientCtx.Codec.MustUnmarshalJSON(appState[v040wrkchain.ModuleName], &oldWrkchainGenState)

		// delete deprecated x/wrkchain genesis state
		delete(appState, v040wrkchain.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v045wrkchain.ModuleName] = clientCtx.Codec.MustMarshalJSON(v045wrkchain.MigrateJSON(&oldWrkchainGenState))
	}

	// migrate enterprise
	if appState[v040ent.ModuleName] != nil {
		var oldEnterpriseGenState v040ent.GenesisState
		clientCtx.Codec.MustUnmarshalJSON(appState[v040ent.ModuleName], &oldEnterpriseGenState)

		// delete deprecated x/enterprise genesis state
		delete(appState, v040ent.ModuleName)

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState[v045ent.ModuleName] = clientCtx.Codec.MustMarshalJSON(v045ent.MigrateJSON(&oldEnterpriseGenState))
	}

	return appState
}
