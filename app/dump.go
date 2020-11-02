package app

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/beacon"
	"github.com/unification-com/mainchain/x/wrkchain"
)

// DumpWrkchainOrBeaconData dump data for selected WRKChain or BEACON ready for import into new chain
func (app *MainchainApp) DumpWrkchainOrBeaconData(what string,
	id uint64) (dataDump json.RawMessage, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	switch what {
	case "beacon":
		beaconRegData := app.beaconKeeper.GetBeacon(ctx, id)
		beaconTimestampState := app.beaconKeeper.GetAllBeaconTimestampsForExport(ctx, id)

		dumpState := beacon.BeaconExport{
			Beacon:           beaconRegData,
			BeaconTimestamps: beaconTimestampState,
		}

		dumpStateJson, err := codec.MarshalJSONIndent(app.cdc, dumpState)
		if err != nil {
			return nil,  err
		}

		return dumpStateJson, nil
	case "wrkchain":
		wrkchainRegData := app.wrkChainKeeper.GetWrkChain(ctx, id)
		wrkchainHashState := app.wrkChainKeeper.GetAllWrkChainBlockHashesForGenesisExport(ctx, id)

		dumpState := wrkchain.WrkChainExport{
			WrkChain: wrkchainRegData,
			WrkChainBlocks: wrkchainHashState,
		}

		dumpStateJson, err := codec.MarshalJSONIndent(app.cdc, dumpState)
		if err != nil {
			return nil,  err
		}

		return dumpStateJson, nil
	default:
		return nil, nil
	}
}
