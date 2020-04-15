package simapp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	undtypes "github.com/unification-com/mainchain/types"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestSimAppExport(t *testing.T) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(undtypes.Bech32PrefixAccAddr, undtypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(undtypes.Bech32PrefixValAddr, undtypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(undtypes.Bech32PrefixConsAddr, undtypes.Bech32PrefixConsPub)
	config.SetCoinType(undtypes.CoinType)
	config.SetFullFundraiserPath(undtypes.HdWalletPath)
	config.Seal()
	db := dbm.NewMemDB()
	app := NewUndSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	genesisState := NewDefaultGenesisState()

	stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
	require.NoError(t, err)

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewUndSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err = app2.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := dbm.NewMemDB()
	app := NewUndSimApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, app.BankKeeper.BlacklistedAddr(app.SupplyKeeper.GetModuleAddress(acc)))
	}
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}
