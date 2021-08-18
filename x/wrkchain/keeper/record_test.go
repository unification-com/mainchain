package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/wrkchain/types"
	"testing"
	"time"
)

// Tests for Highest WRKChain ID

func TestSetGetWrkChainBlock(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(128)
		moniker := test_helpers.GenerateRandomString(64)
		genesisHash := test_helpers.GenerateRandomString(66)

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.WrkChainBlock{}
			block.WrkchainId = wcID
			block.Owner = addr.String()
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)

			blockDb := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, h)
			require.True(t, WRKChainBlockEqual(blockDb, block))
		}
	}
}

func TestIsWrkChainBlockRecorded(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(128)
		moniker := test_helpers.GenerateRandomString(64)
		genesisHash := test_helpers.GenerateRandomString(66)

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.WrkChainBlock{}
			block.WrkchainId = wcID
			block.Owner = addr.String()
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)

			isRecorded := app.WrkchainKeeper.IsWrkChainBlockRecorded(ctx, wcID, h)
			require.True(t, isRecorded)
		}
	}
}

func TestGetWrkChainBlockHashes(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	numToRecord := uint64(1000)

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(128)
		moniker := test_helpers.GenerateRandomString(64)
		genesisHash := test_helpers.GenerateRandomString(66)

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		var testBlocks []types.WrkChainBlock

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.WrkChainBlock{}
			block.WrkchainId = wcID
			block.Owner = addr.String()
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			testBlocks = append(testBlocks, block)

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)
		}

		allBlocks := app.WrkchainKeeper.GetAllWrkChainBlockHashes(ctx, wcID)
		require.True(t, len(allBlocks) == int(numToRecord) && len(allBlocks) == len(testBlocks))

		for i := 0; i < int(numToRecord); i++ {
			require.True(t, WRKChainBlockEqual(allBlocks[i], testBlocks[i]))
		}
	}
}

func TestIsAuthorisedToRecord(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	unauthorisedAddr := sdk.AccAddress(pubKey.Address())

	for _, addr := range testAddrs {
		name := test_helpers.GenerateRandomString(128)
		moniker := test_helpers.GenerateRandomString(64)
		genesisHash := test_helpers.GenerateRandomString(66)

		wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		isAuthorised := app.WrkchainKeeper.IsAuthorisedToRecord(ctx, wcID, addr)
		require.True(t, isAuthorised)

		isAuthorised = app.WrkchainKeeper.IsAuthorisedToRecord(ctx, wcID, unauthorisedAddr)
		require.False(t, isAuthorised)
	}
}

func TestRecordWrkchainHashes(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(10)
	numToRecord := uint64(100)

	name := test_helpers.GenerateRandomString(128)
	moniker := test_helpers.GenerateRandomString(64)
	genesisHash := test_helpers.GenerateRandomString(66)

	wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", testAddrs[0])
	require.NoError(t, err)

	for h := uint64(1); h <= numToRecord; h++ {
		expectedBlock := types.WrkChainBlock{}
		expectedBlock.WrkchainId = wcID
		expectedBlock.Owner = testAddrs[0].String()
		expectedBlock.Height = h
		expectedBlock.Blockhash = test_helpers.GenerateRandomString(66)
		expectedBlock.Parenthash = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash1 = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash2 = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash3 = test_helpers.GenerateRandomString(66)
		expectedBlock.SubTime = uint64(time.Now().Unix())

		err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcID, h, expectedBlock.Blockhash, expectedBlock.Parenthash, expectedBlock.Hash1, expectedBlock.Hash2, expectedBlock.Hash3, testAddrs[0])
		require.NoError(t, err)

		blockDb := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, h)
		// hackery
		expectedBlock.SubTime = blockDb.SubTime
		require.True(t, WRKChainBlockEqual(blockDb, expectedBlock))

		wrkChainDb, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcID)

		require.True(t, wrkChainDb.Lastblock == h)
	}

}
