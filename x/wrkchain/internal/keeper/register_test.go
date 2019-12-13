package keeper

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

// Tests for Highest WRKChain ID

func TestSetGetHighestWRKChainID(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i <= 1000; i++ {
		keeper.SetHighestWrkChainID(ctx, i)
		wcID, err := keeper.GetHighestWrkChainID(ctx)
		require.NoError(t, err)
		require.True(t, wcID == i)
	}
}

// Tests for Get/Set WRKChains

func TestSetGetWrkChain(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 1000)

	wcID := uint64(1)
	for _, addr := range TestAddrs {

		wc := types.NewWrkChain()
		wc.Owner = addr
		wc.WrkChainID = wcID
		wc.LastBlock = 1
		wc.RegisterTime = time.Now().Unix()
		wc.Moniker = GenerateRandomString(12)
		wc.Name = GenerateRandomString(20)
		wc.GenesisHash = GenerateRandomString(32)

		err := keeper.SetWrkChain(ctx, wc)
		require.NoError(t, err)

		isRegistered := keeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)

		wcDb := keeper.GetWrkChain(ctx, wcID)
		require.True(t, WRKChainEqual(wcDb, wc))

		wcDbOwner := keeper.GetWrkChainOwner(ctx, wcID)
		require.True(t, wcDbOwner.String() == addr.String())

		wcID = wcID + 1
	}
}

func TestSetLastBlock(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	lastBlock := uint64(0)
	wcID := uint64(1)

	wc := types.NewWrkChain()
	wc.Owner = TestAddrs[0]
	wc.WrkChainID = wcID
	wc.LastBlock = lastBlock
	wc.RegisterTime = time.Now().Unix()
	wc.Moniker = GenerateRandomString(12)
	wc.Name = GenerateRandomString(20)
	wc.GenesisHash = GenerateRandomString(32)

	err := keeper.SetWrkChain(ctx, wc)
	require.NoError(t, err)

	for i := uint64(1); i <= 1000; i++ {
		err := keeper.SetLastBlock(ctx, wcID, i)
		require.NoError(t, err)

		wcDb := keeper.GetWrkChain(ctx, wcID)
		require.True(t, wcDb.LastBlock == i)
		lastBlock = i
	}

	// check can't set last block to < current last block
	oldBlock := lastBlock - 1
	err = keeper.SetLastBlock(ctx, wcID, oldBlock)
	require.NoError(t, err)
	wcDb := keeper.GetWrkChain(ctx, wcID)
	require.True(t, wcDb.LastBlock == lastBlock)

}

func TestEmptyWrkChainValuesReturnError(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	wc0 := types.NewWrkChain()

	wc1 := wc0
	wc1.Owner = TestAddrs[0]

	wc2 := wc1
	wc2.WrkChainID = uint64(1)

	wc3 := wc2
	wc3.Moniker = "moniker"

	testCases := []struct {
		wc          types.WrkChain
		expectedErr sdk.Error
	}{
		{wc0, sdk.ErrInternal("unable to set WRKChain - must have an owner")},
		{wc1, sdk.ErrInternal("unable to set WRKChain - id must be positive non-zero")},
		{wc2, sdk.ErrInternal("unable to set WRKChain - must have a moniker")},
		{wc3, nil},
	}

	for _, tc := range testCases {
		err := keeper.SetWrkChain(ctx, tc.wc)
		require.Equal(t, tc.expectedErr, err, "unexpected type of error: %s", err)
	}
}

// Tests for Registering a new WRKChain

func TestRegisterWrkChain(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 1000)

	i, _ := keeper.GetHighestWrkChainID(ctx)

	for _, addr := range TestAddrs {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		genesisHash := GenerateRandomString(32)

		expectedWc := types.NewWrkChain()
		expectedWc.Owner = addr
		expectedWc.WrkChainID = i
		expectedWc.LastBlock = 0
		expectedWc.RegisterTime = time.Now().Unix()
		expectedWc.Moniker = moniker
		expectedWc.Name = name
		expectedWc.GenesisHash = genesisHash
		expectedWc.BaseType = "geth"

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", addr)
		require.NoError(t, err)
		require.True(t, wcID == expectedWc.WrkChainID)

		isRegistered := keeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)

		wcDb := keeper.GetWrkChain(ctx, wcID)

		// hackery for reg time, otherwise following test fails
		expectedWc.RegisterTime = wcDb.RegisterTime
		require.True(t, WRKChainEqual(wcDb, expectedWc))

		wcDbOwner := keeper.GetWrkChainOwner(ctx, wcID)
		require.True(t, wcDbOwner.String() == addr.String())

		i = i + 1
	}
}

func TestFailRegisterNewWrkChain(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	testCases := []struct {
		moniker      string
		name         string
		genHash      string
		owner        sdk.AccAddress
		expectedErr  sdk.Error
		expectedWcID uint64
	}{
		{"moniker", "name", "genhash", sdk.AccAddress{}, sdk.ErrInternal("unable to set WRKChain - must have an owner"), 0},
		{"", "name", "genhash", TestAddrs[0], sdk.ErrInternal("unable to set WRKChain - must have a moniker"), 0},
		{"", "", "genhash", TestAddrs[0], sdk.ErrInternal("unable to set WRKChain - must have a moniker"), 0},
		{"", "name", "", TestAddrs[0], sdk.ErrInternal("unable to set WRKChain - must have a moniker"), 0},
		{"testmoniker", "", "", TestAddrs[0], nil, 1},
		{"testmoniker", "", "", TestAddrs[0], types.ErrWrkChainAlreadyRegistered(keeper.codespace, fmt.Sprintf("wrkchain already registered with moniker 'testmoniker' - id: 1, owner: %s", TestAddrs[0])), 0},
	}

	for _, tc := range testCases {
		wcID, err := keeper.RegisterWrkChain(ctx, tc.moniker, tc.name, tc.genHash, "geth", tc.owner)
		require.Equal(t, tc.expectedErr, err, "unexpected type of error: %s", err)
		require.True(t, wcID == tc.expectedWcID)
	}
}

func TestHighestWrkChainIdAfterRegister(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i < 1000; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		genesisHash := GenerateRandomString(32)
		owner := TestAddrs[1]

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		require.NoError(t, err)

		nextID, _ := keeper.GetHighestWrkChainID(ctx)
		expectedNextID := wcID + 1
		require.True(t, nextID == expectedNextID)
	}
}

func TestWrkChainIsRegisteredAfterRegister(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)

	for i := uint64(1); i < 1000; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		genesisHash := GenerateRandomString(32)
		owner := TestAddrs[1]

		wcID, err := keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		require.NoError(t, err)

		isRegistered := keeper.IsWrkChainRegistered(ctx, wcID)
		require.True(t, isRegistered)
	}
}

func TestGetWrkChainFilter(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 100, 0)
	numToReg := 100
	lastMoniker := ""

	for i := 0; i < numToReg; i++ {
		name := GenerateRandomString(20)
		moniker := GenerateRandomString(12)
		genesisHash := GenerateRandomString(32)
		owner := TestAddrs[1]

		_, _ = keeper.RegisterWrkChain(ctx, moniker, name, genesisHash, "geth", owner)
		lastMoniker = moniker
	}

	params := types.NewQueryWrkChainParams(1, 1000, "", TestAddrs[1])
	results := keeper.GetWrkChainsFiltered(ctx, params)
	require.True(t, len(results) == numToReg)

	params = types.NewQueryWrkChainParams(1, 1000, lastMoniker, sdk.AccAddress{})
	results = keeper.GetWrkChainsFiltered(ctx, params)
	require.True(t, len(results) == 1)
}
