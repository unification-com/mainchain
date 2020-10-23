package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestEqualPurchaseOrderID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	state1.StartingWrkChainID = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	state2.StartingWrkChainID = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}

func TestNewGenesisState(t *testing.T) {
	params1 := NewParams(1000, 100, "nund")
	state1 := NewGenesisState(params1, 1)

	params2 := NewParams(1000, 100, "nund")
	state2 := NewGenesisState(params2, 1)

	require.Equal(t, state1, state2)
}

func TestDefaultGenesisState(t *testing.T) {
	state1 := DefaultGenesisState()
	state2 := DefaultGenesisState()

	require.Equal(t, state1, state2)
}

func TestValidateGenesis(t *testing.T) {
	state1 := DefaultGenesisState()
	err := ValidateGenesis(state1)
	require.NoError(t, err)

	state2 := GenesisState{}
	err = ValidateGenesis(state2)
	require.Error(t, err)

	state3 := DefaultGenesisState()
	wrkchain1 := WrkChainExport{
		WrkChain: WrkChain{
			WrkChainID: 0,
		},
	}

	state3.WrkChains = append(state3.WrkChains, wrkchain1)

	expectedErr := fmt.Errorf("invalid WrkChain: ID: %d. Error: Missing ID", 0)
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.WrkChains[0].WrkChain.WrkChainID = 1
	expectedErr = fmt.Errorf("invalid WrkChain: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	bOwnerAddr := sdk.AccAddress(pubKey.Address())
	state3.WrkChains[0].WrkChain.Owner = bOwnerAddr

	expectedErr = fmt.Errorf("invalid Beacon: Moniker: . Error: Missing Moniker")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.WrkChains[0].WrkChain.Moniker = "wrkchain"
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	expectedErr = fmt.Errorf("invalid Beacon: BaseType: . Error: Missing BaseType")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.WrkChains[0].WrkChain.BaseType = "tendermint"
	err = ValidateGenesis(state3)
	require.NoError(t, err)

	block := WrkChainBlockGenesisExport{}
	state3.WrkChains[0].WrkChainBlocks = append(state3.WrkChains[0].WrkChainBlocks, block)

	expectedErr = fmt.Errorf("invalid WrkChain block: WrkChainID: 0. Error: Missing WrkChainID")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.WrkChains[0].WrkChainBlocks[0].BlockHash = "ljbhouhgygiuyiug"
	expectedErr = fmt.Errorf("invalid WrkChain block: Height: . Error: Missing Height")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.WrkChains[0].WrkChainBlocks[0].Height = 12345
	err = ValidateGenesis(state3)
	require.NoError(t, err)
}
