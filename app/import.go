package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon"
	"github.com/unification-com/mainchain/x/wrkchain"
)

func (app *MainchainApp) ImportWrkchainAndBeaconData(ctx sdk.Context, nodeHome string) {

	genesisDataDirBeacon := filepath.Join(nodeHome, "config/genesis_data/beacon")
	genesisDataDirWrkChain := filepath.Join(nodeHome, "config/genesis_data/wrkchain")

	ctx.Logger().Info("import additional beacon and wrkchain data")

	_, err := os.Stat(genesisDataDirBeacon)
	if err == nil {
		c, _ := ioutil.ReadDir(genesisDataDirBeacon)
		for _, entry := range c {
			dataFile := filepath.Join(genesisDataDirBeacon, entry.Name())
			ctx.Logger().Info("attempt to import beacon data", "dataFile", dataFile)
			app.importBeaconData(ctx, dataFile)
			ctx.Logger().Info("finished importing beacon data", "dataFile", dataFile)
		}
	}

	_, err = os.Stat(genesisDataDirWrkChain)
	if err == nil {
		c, _ := ioutil.ReadDir(genesisDataDirWrkChain)
		for _, entry := range c {
			dataFile := filepath.Join(genesisDataDirWrkChain, entry.Name())
			ctx.Logger().Info("attempt to import wrkchain data", "dataFile", dataFile)
			app.importWrkchainData(ctx, dataFile)
			ctx.Logger().Info("finished importing wrkchain data", "dataFile", dataFile)
		}
	}

	ctx.Logger().Info("genesis and data import complete!")
}

func getBytes(p string) ([]byte, error) {
	jsonFile, err := os.Open(p)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return byteValue, nil
}

// importBeaconData will import BEACON timestamps from a given file.
// An error at any point will induce panic
func (app *MainchainApp) importBeaconData(ctx sdk.Context, p string) {

	byteValue, err := getBytes(p)
	if err != nil {
		panic(err)
	}

	var beaconImport beacon.BeaconExport

	err = app.cdc.UnmarshalJSON(byteValue, &beaconImport)

	if err != nil {
		panic(err)
	}

	beaconData := beaconImport.Beacon
	timestamps := beaconImport.BeaconTimestamps

	// check BEACON exists in the current state - i.e. has been imported from genesis.json
	// ID and Moniker in the data file must match those imported form genesis.json
	b := app.beaconKeeper.GetBeacon(ctx, beaconData.BeaconID)
	if b.Moniker != beaconData.Moniker || b.BeaconID != beaconData.BeaconID {
		panic(fmt.Sprintf("cannot find BEACON %d (%s) in state from genesis!",
			beaconData.BeaconID, beaconData.Moniker))
	}

	ctx.Logger().Info("found state for beacon", "id", beaconData.BeaconID, "moniker", beaconData.Moniker)

	if b.LastTimestampID > 0 {
		panic("beacon LastTimestampID > 0 - was this already imported in genesis.json?")
	}

	// registration data has already been imported during InitGenesis.
	// just need to set LastTimestampID and import the timestamps
	if beaconData.LastTimestampID > 0 {
		err = app.beaconKeeper.SetLastTimestampID(ctx, beaconData.BeaconID, beaconData.LastTimestampID)
		if err != nil {
			panic(err)
		}
	}

	for _, timestamp := range timestamps {

		bts := beacon.BeaconTimestamp{
			BeaconID:    beaconData.BeaconID,
			TimestampID: timestamp.TimestampID,
			SubmitTime:  timestamp.SubmitTime,
			Hash:        timestamp.Hash,
			Owner:       beaconData.Owner,
		}

		err = app.beaconKeeper.SetBeaconTimestamp(ctx, bts)
		if err != nil {
			panic(err)
		}
	}
}

// importWrkchainData will import WRKChain block hashes from a given file.
// An error at any point will induce panic
func (app *MainchainApp) importWrkchainData(ctx sdk.Context, p string) {
	byteValue, err := getBytes(p)

	if err != nil {
		panic(err)
	}

	var wrkchainImport wrkchain.WrkChainExport
	err = app.cdc.UnmarshalJSON(byteValue, &wrkchainImport)

	if err != nil {
		panic(err)
	}

	wrkchainData := wrkchainImport.WrkChain
	wrkchainBlocks := wrkchainImport.WrkChainBlocks

	// check WRKChain exists in the current state - i.e. has been imported from genesis.json
	// ID and Moniker in the data file must match those imported form genesis.json
	w := app.wrkChainKeeper.GetWrkChain(ctx, wrkchainData.WrkChainID)
	if w.Moniker != wrkchainData.Moniker || w.WrkChainID != wrkchainData.WrkChainID {
		panic(fmt.Sprintf("cannot find WRKChain %d (%s) in state from genesis!",
			wrkchainData.WrkChainID, wrkchainData.Moniker))
	}

	ctx.Logger().Info("found state for wrkchain", "id", wrkchainData.WrkChainID,
		"moniker", wrkchainData.Moniker)

	if w.LastBlock > 0 || w.NumberBlocks > 0 {
		panic("wrkchain LastBlock || NumberBlocks > 0 - was this already imported in genesis.json?")
	}

	// registration data has already been imported during InitGenesis.
	// just need to set LastBlock and import the timestamps
	err = app.wrkChainKeeper.SetLastBlock(ctx, wrkchainData.WrkChainID, wrkchainData.LastBlock)
	if err != nil {
		panic(err)
	}

	for _, block := range wrkchainBlocks {
		blk := wrkchain.WrkChainBlock{
			WrkChainID: wrkchainData.WrkChainID,
			Height:     block.Height,
			BlockHash:  block.BlockHash,
			ParentHash: block.ParentHash,
			Hash1:      block.Hash1,
			Hash2:      block.Hash2,
			Hash3:      block.Hash3,
			SubmitTime: block.SubmitTime,
			Owner:      wrkchainData.Owner,
		}

		err = app.wrkChainKeeper.SetWrkChainBlock(ctx, blk)
		if err != nil {
			panic(err)
		}

		// also update NumBlocks!
		err = app.wrkChainKeeper.SetNumBlocks(ctx, wrkchainData.WrkChainID)
		if err != nil {
			panic(err)
		}
	}
}
