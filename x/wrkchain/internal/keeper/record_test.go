package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

// Tests for Highest WRKChain ID

func TestSetGetWrkChainBlock(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(128)
		moniker := GenerateRandomString(64)
		genesisHash := GenerateRandomString(66)

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.NewWrkChainBlock()
			block.WrkChainID = wcID
			block.Owner = addr
			block.Height = h
			block.BlockHash = GenerateRandomString(66)
			block.ParentHash = GenerateRandomString(66)
			block.Hash1 = GenerateRandomString(66)
			block.Hash2 = GenerateRandomString(66)
			block.Hash3 = GenerateRandomString(66)
			block.SubmitTime = time.Now().Unix()

			err := keeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)

			blockDb := keeper.GetWrkChainBlock(ctx, wcID, h)
			require.True(t, WRKChainBlockEqual(blockDb, block))
		}
	}
}

func TestIsWrkChainBlockRecorded(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(128)
		moniker := GenerateRandomString(64)
		genesisHash := GenerateRandomString(66)

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.NewWrkChainBlock()
			block.WrkChainID = wcID
			block.Owner = addr
			block.Height = h
			block.BlockHash = GenerateRandomString(66)
			block.ParentHash = GenerateRandomString(66)
			block.Hash1 = GenerateRandomString(66)
			block.Hash2 = GenerateRandomString(66)
			block.Hash3 = GenerateRandomString(66)
			block.SubmitTime = time.Now().Unix()

			err := keeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)

			isRecorded := keeper.IsWrkChainBlockRecorded(ctx, wcID, h)
			require.True(t, isRecorded)
		}
	}
}

func TestGetWrkChainBlockHashes(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(1000)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(128)
		moniker := GenerateRandomString(64)
		genesisHash := GenerateRandomString(66)

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		var testBlocks []types.WrkChainBlock

		for h := uint64(1); h <= numToRecord; h++ {
			block := types.NewWrkChainBlock()
			block.WrkChainID = wcID
			block.Owner = addr
			block.Height = h
			block.BlockHash = GenerateRandomString(66)
			block.ParentHash = GenerateRandomString(66)
			block.Hash1 = GenerateRandomString(66)
			block.Hash2 = GenerateRandomString(66)
			block.Hash3 = GenerateRandomString(66)
			block.SubmitTime = time.Now().Unix()

			testBlocks = append(testBlocks, block)

			err := keeper.SetWrkChainBlock(ctx, block)
			require.NoError(t, err)
		}

		allBlocks := keeper.GetAllWrkChainBlockHashes(ctx, wcID)
		require.True(t, len(allBlocks) == int(numToRecord) && len(allBlocks) == len(testBlocks))

		for i := 0; i < int(numToRecord); i++ {
			require.True(t, WRKChainBlockEqual(allBlocks[i], testBlocks[i]))
		}
	}
}

func TestIsAuthorisedToRecord(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 100)

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	unauthorisedAddr := sdk.AccAddress(pubKey.Address())

	for _, addr := range TestAddrs {
		name := GenerateRandomString(128)
		moniker := GenerateRandomString(64)
		genesisHash := GenerateRandomString(66)

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)

		isAuthorised := keeper.IsAuthorisedToRecord(ctx, wcID, addr)
		require.True(t, isAuthorised)

		isAuthorised = keeper.IsAuthorisedToRecord(ctx, wcID, unauthorisedAddr)
		require.False(t, isAuthorised)
	}
}

func TestRecordWrkchainHashes(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToRecord := uint64(100)

	name := GenerateRandomString(128)
	moniker := GenerateRandomString(64)
	genesisHash := GenerateRandomString(66)

	wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", TestAddrs[0])
	require.NoError(t, err)

	for h := uint64(1); h <= numToRecord; h++ {
		expectedBlock := types.NewWrkChainBlock()
		expectedBlock.WrkChainID = wcID
		expectedBlock.Owner = TestAddrs[0]
		expectedBlock.Height = h
		expectedBlock.BlockHash = GenerateRandomString(66)
		expectedBlock.ParentHash = GenerateRandomString(66)
		expectedBlock.Hash1 = GenerateRandomString(66)
		expectedBlock.Hash2 = GenerateRandomString(66)
		expectedBlock.Hash3 = GenerateRandomString(66)
		expectedBlock.SubmitTime = time.Now().Unix()

		err := keeper.RecordWrkchainHashes(ctx, wcID, h, expectedBlock.BlockHash, expectedBlock.ParentHash, expectedBlock.Hash1, expectedBlock.Hash2, expectedBlock.Hash3, TestAddrs[0])
		require.NoError(t, err)

		blockDb := keeper.GetWrkChainBlock(ctx, wcID, h)
		// hackery
		expectedBlock.SubmitTime = blockDb.SubmitTime
		require.True(t, WRKChainBlockEqual(blockDb, expectedBlock))
	}

}

func TestRecordWrkchainHashesFail(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	name := GenerateRandomString(128)
	moniker := GenerateRandomString(64)
	genesisHash := GenerateRandomString(66)

	goodHash := GenerateRandomString(66)

	wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", TestAddrs[0])
	require.NoError(t, err)

	testCases := []struct {
		wrkchainId  uint64
		height      uint64
		blockHash   string
		parentHash  string
		hash1       string
		hash2       string
		hash3       string
		owner       sdk.AccAddress
		expectedErr error
	}{
		{wcID, 1, "", "", "", "", "", TestAddrs[0], sdkerrors.Wrap(types.ErrMissingData, "must include owner, id, height and hash")},
		{wcID, 0, goodHash, "", "", "", "", TestAddrs[0], sdkerrors.Wrap(types.ErrMissingData, "must include owner, id, height and hash")},
		{wcID, 1, goodHash, "", "", "", "", TestAddrs[0], nil},
	}

	for _, tc := range testCases {
		err := keeper.RecordWrkchainHashes(ctx, tc.wrkchainId, tc.height, tc.blockHash, tc.parentHash, tc.hash1, tc.hash2, tc.hash3, tc.owner)
		if tc.expectedErr != nil {
			require.Equal(t, tc.expectedErr.Error(), err.Error(), "unexpected type of error: %s", err.Error())
		} else {
			require.Nil(t, err)
		}
	}
}
