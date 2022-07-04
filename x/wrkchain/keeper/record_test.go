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
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, wcID, block)
			require.NoError(t, err)

			blockDb, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, h)
			require.True(t, found)
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
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, wcID, block)
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
			block.Height = h
			block.Blockhash = test_helpers.GenerateRandomString(66)
			block.Parenthash = test_helpers.GenerateRandomString(66)
			block.Hash1 = test_helpers.GenerateRandomString(66)
			block.Hash2 = test_helpers.GenerateRandomString(66)
			block.Hash3 = test_helpers.GenerateRandomString(66)
			block.SubTime = uint64(time.Now().Unix())

			testBlocks = append(testBlocks, block)

			err := app.WrkchainKeeper.SetWrkChainBlock(ctx, wcID, block)
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
	numToRecord := uint64(1000)
	recordLimit := uint64(200)
	startHeight := uint64(24)
	endHeight := startHeight + numToRecord

	name := test_helpers.GenerateRandomString(128)
	moniker := test_helpers.GenerateRandomString(64)
	genesisHash := test_helpers.GenerateRandomString(66)

	wcID, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, moniker, name, genesisHash, "geth", testAddrs[0])
	require.NoError(t, err)

	// set the record limit
	wc, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcID)
	wc.InStateLimit = recordLimit
	err = app.WrkchainKeeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for h := startHeight; h < endHeight; h++ {
		expectedBlock := types.WrkChainBlock{}
		expectedBlock.Height = h
		expectedBlock.Blockhash = test_helpers.GenerateRandomString(66)
		expectedBlock.Parenthash = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash1 = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash2 = test_helpers.GenerateRandomString(66)
		expectedBlock.Hash3 = test_helpers.GenerateRandomString(66)
		expectedBlock.SubTime = uint64(time.Now().Unix())

		_, err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcID, h, expectedBlock.Blockhash, expectedBlock.Parenthash, expectedBlock.Hash1, expectedBlock.Hash2, expectedBlock.Hash3)
		require.NoError(t, err)

		blockDb, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, h)
		require.True(t, found)
		// hackery
		expectedBlock.SubTime = blockDb.SubTime
		require.True(t, WRKChainBlockEqual(blockDb, expectedBlock))

		wrkChainDb, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcID)

		require.True(t, wrkChainDb.Lastblock == h)
	}

	wrkChainDb, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcID)
	require.True(t, wrkChainDb.NumBlocks == recordLimit)
	require.True(t, wrkChainDb.Lastblock == endHeight-1)
	require.True(t, wrkChainDb.LowestHeight == numToRecord-recordLimit+startHeight)

	// should still be in state
	blockCount := uint64(0)
	for height := wrkChainDb.LowestHeight; height < endHeight; height++ {
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, height)
		require.True(t, found)
		blockCount++
	}
	require.True(t, blockCount == recordLimit)

	// should no longer be in state
	for height := startHeight; height < wrkChainDb.LowestHeight; height++ {
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcID, height)
		require.False(t, found)
	}

}

func TestIncreaseInStateStorage(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	recordLimitIncrease := uint64(200)

	wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", testAddrs[0])
	require.NoError(t, err)

	wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.InStateLimit == types.DefaultStorageLimit)

	err = app.WrkchainKeeper.IncreaseInStateStorage(ctx, wcId, recordLimitIncrease)
	require.NoError(t, err)

	wrkchain, found = app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.InStateLimit == types.DefaultStorageLimit+recordLimitIncrease)
}

func TestIncreaseInStateStorageWithBlockHashRecording(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	numToRecord := uint64(500)
	recordLimit := uint64(100)
	increaseAmount := uint64(50)

	wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", testAddrs[0])
	require.NoError(t, err)

	// set the record limit
	wc, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	wc.InStateLimit = recordLimit
	err = app.WrkchainKeeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	// record initial hashes
	for i := uint64(1); i <= numToRecord; i++ {
		hash := test_helpers.GenerateRandomString(32)
		deletedHeight, err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcId, i, hash, "", "", "", "")
		require.NoError(t, err)

		if i < recordLimit {
			require.True(t, deletedHeight == 0)
		}

		if deletedHeight > 0 {
			_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, deletedHeight)
			require.False(t, found)
		}
	}

	// sanity check
	wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.NumBlocks == recordLimit)
	require.True(t, wrkchain.LowestHeight == numToRecord-recordLimit+1)

	// increase storage capacity
	err = app.WrkchainKeeper.IncreaseInStateStorage(ctx, wcId, increaseAmount)
	require.NoError(t, err)
	wrkchain, found = app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.InStateLimit == recordLimit+increaseAmount)

	// record new timestamps
	for i := numToRecord + 1; i <= numToRecord+numToRecord; i++ {
		hash := test_helpers.GenerateRandomString(32)
		deletedHeight, err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcId, i, hash, "", "", "", "")
		require.NoError(t, err)

		if deletedHeight > 0 {
			_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, deletedHeight)
			require.False(t, found)
		}
	}

	// check final result
	wrkchain, found = app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.NumBlocks == recordLimit+increaseAmount)
	require.True(t, wrkchain.LowestHeight == (numToRecord*2)-recordLimit-increaseAmount+1)
}

func TestAsymmetricRecordNewWrkchainHashesAndDeleteOld(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	heightsToRecord := [10]uint64{1, 3, 4, 8, 11, 24, 33, 34, 40, 50}
	expectedDeleteds := [10]uint64{0, 0, 0, 0, 0, 1, 3, 4, 8, 11}
	expectedLowests := [10]uint64{1, 1, 1, 1, 1, 3, 4, 8, 11, 24}
	expectedHighests := [10]uint64{1, 3, 4, 8, 11, 24, 33, 34, 40, 50}
	expectedNumBlocks := [10]uint64{1, 2, 3, 4, 5, 5, 5, 5, 5, 5}

	finalInStates := [5]uint64{24, 33, 34, 40, 50}
	finalNotInStates := [5]uint64{1, 3, 4, 8, 11}

	wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", testAddrs[0])
	require.NoError(t, err)

	// set the record limit
	wc, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	wc.InStateLimit = 5
	err = app.WrkchainKeeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		heightToRecord := heightsToRecord[i]
		expectedDeleted := expectedDeleteds[i]
		expectedLowest := expectedLowests[i]
		expectedHighest := expectedHighests[i]
		expectedNumBlock := expectedNumBlocks[i]

		hash := test_helpers.GenerateRandomString(32)
		deletedHeight, err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcId, heightToRecord, hash, "", "", "", "")
		require.NoError(t, err)

		// should be found
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, heightToRecord)
		require.True(t, found)

		wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
		require.True(t, found)
		require.True(t, wrkchain.NumBlocks == expectedNumBlock)
		require.True(t, wrkchain.LowestHeight == expectedLowest)
		require.True(t, wrkchain.Lastblock == expectedHighest)

		_, found = app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, expectedDeleted)
		require.False(t, found)

		_, found = app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, deletedHeight)
		require.False(t, found)

		require.True(t, deletedHeight == expectedDeleted)
	}

	for i := 0; i < 5; i++ {
		finalInState := finalInStates[i]
		finalNotInState := finalNotInStates[i]

		// should be found
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, finalInState)
		require.True(t, found)

		// should not be found
		_, found = app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, finalNotInState)
		require.False(t, found)
	}

	// final check
	wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.NumBlocks == 5)
	require.True(t, wrkchain.LowestHeight == 24)
	require.True(t, wrkchain.Lastblock == 50)

}

func TestAsymmetricRecordNewWrkchainHashesAndDeleteOldWithIncrease(t *testing.T) {
	app := test_helpers.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	testAddrs := test_helpers.GenerateRandomTestAccounts(1)

	// add 15 wrkchain blocks with asymmetric heights
	heightsToRecord := [15]uint64{1, 3, 4, 8, 11, 24, 33, 34, 40, 50, 55, 69, 78, 99, 108}
	// these heights should no longer exists at nth submission
	// nothing should be deleted from 10th to 12th, since storage is increased by 3
	expectedDeleteds := [15]uint64{0, 0, 0, 0, 0, 1, 3, 4, 8, 0, 0, 0, 11, 24, 33}
	// these heights should be lowest in state at nth submission
	expectedLowests := [15]uint64{1, 1, 1, 1, 1, 3, 4, 8, 11, 11, 11, 11, 24, 33, 34}
	// these heights should be highest in state at nth submission
	expectedHighests := [15]uint64{1, 3, 4, 8, 11, 24, 33, 34, 40, 50, 55, 69, 78, 99, 108}
	// expected number of blocks held in state at nth submission. 5th to 9th should be
	// 5, the current limit. Storage increases at 10th, so should be 8 onwards
	expectedNumBlocks := [15]uint64{1, 2, 3, 4, 5, 5, 5, 5, 5, 6, 7, 8, 8, 8, 8}
	// at 10th, increase storage by 3 (before 10th submission)
	increaseStorages := [15]uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0}

	// at each of the 15 submission points, these heights should exists in state
	// this is a matrix representation of the above config
	expectedStillInStateAtThisPoints := [15][]uint64{
		{1},
		{1, 3},
		{1, 3, 4},
		{1, 3, 4, 8},
		{1, 3, 4, 8, 11},
		{3, 4, 8, 11, 24},
		{4, 8, 11, 24, 33},
		{8, 11, 24, 33, 34},
		{11, 24, 33, 34, 40},
		{11, 24, 33, 34, 40, 50},
		{11, 24, 33, 34, 40, 50, 55},
		{11, 24, 33, 34, 40, 50, 55, 69},
		{24, 33, 34, 40, 50, 55, 69, 78},
		{33, 34, 40, 50, 55, 69, 78, 99},
		{34, 40, 50, 55, 69, 78, 99, 108},
	}

	// at each of the 15 submission points, these heights should no longer exist in state
	// this is a matrix representation of the above config
	expectedDeletedAtThisPoints := [15][]uint64{
		{0},
		{0},
		{0},
		{0},
		{0},
		{1},
		{1, 3},
		{1, 3, 4},
		{1, 3, 4, 8},
		{1, 3, 4, 8},
		{1, 3, 4, 8},
		{1, 3, 4, 8},
		{1, 3, 4, 8, 11},
		{1, 3, 4, 8, 11, 24},
		{1, 3, 4, 8, 11, 24, 33},
	}

	// only these heights should exist in state at the end
	finalInStates := [8]uint64{34, 40, 50, 55, 69, 78, 99, 108}
	// these heights should not exist in state at the end
	finalNotInStates := [7]uint64{1, 3, 4, 8, 11, 24, 33}

	wcId, err := app.WrkchainKeeper.RegisterNewWrkChain(ctx, "moniker", "name", "ghash", "tm", testAddrs[0])
	require.NoError(t, err)

	// set the record limit
	wc, _ := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	wc.InStateLimit = 5
	err = app.WrkchainKeeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for i := 0; i < 15; i++ {
		heightToRecord := heightsToRecord[i]
		expectedDeleted := expectedDeleteds[i]
		expectedLowest := expectedLowests[i]
		expectedHighest := expectedHighests[i]
		expectedNumBlock := expectedNumBlocks[i]
		increaseStorage := increaseStorages[i]
		expectedStillInStateAtThisPoint := expectedStillInStateAtThisPoints[i]
		expectedDeletedAtThisPoint := expectedDeletedAtThisPoints[i]

		if increaseStorage > 0 {
			err := app.WrkchainKeeper.IncreaseInStateStorage(ctx, wcId, increaseStorage)
			require.NoError(t, err)
		}

		hash := test_helpers.GenerateRandomString(32)
		deletedHeight, err := app.WrkchainKeeper.RecordNewWrkchainHashes(ctx, wcId, heightToRecord, hash, "", "", "", "")
		require.NoError(t, err)

		// this added height should be found!
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, heightToRecord)
		require.True(t, found)

		// deletedHeight should not be found!
		_, found = app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, deletedHeight)
		require.False(t, found)

		// should still be in state, based on matrix
		for _, v := range expectedStillInStateAtThisPoint {
			_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, v)
			require.True(t, found)
		}
		// shoud all be deleted, based on matrix
		for _, v := range expectedDeletedAtThisPoint {
			_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, v)
			require.False(t, found)
		}

		wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
		require.True(t, found)
		require.True(t, wrkchain.NumBlocks == expectedNumBlock)
		require.True(t, wrkchain.LowestHeight == expectedLowest)
		require.True(t, wrkchain.Lastblock == expectedHighest)
		_, found = app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, expectedDeleted)
		require.False(t, found)

		require.True(t, deletedHeight == expectedDeleted)
	}

	// final sanity checks
	for _, v := range finalInStates {
		// should be found
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, v)
		require.True(t, found)
	}
	for _, v := range finalNotInStates {
		// should not be found
		_, found := app.WrkchainKeeper.GetWrkChainBlock(ctx, wcId, v)
		require.False(t, found)
	}

	wrkchain, found := app.WrkchainKeeper.GetWrkChain(ctx, wcId)
	require.True(t, found)
	require.True(t, wrkchain.NumBlocks == 8)
	require.True(t, wrkchain.LowestHeight == 34)
	require.True(t, wrkchain.Lastblock == 108)

}
