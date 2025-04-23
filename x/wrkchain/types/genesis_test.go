package types

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestEqualStartingWrkChainID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	state1.StartingWrkchainId = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.StartingWrkchainId == state2.StartingWrkchainId)

	state2.StartingWrkchainId = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.StartingWrkchainId == state2.StartingWrkchainId)
}

func TestNewGenesisState(t *testing.T) {
	params1 := NewParams(1000, 100, 100, "nund", 200, 300)
	state1 := NewGenesisState(params1, 1, nil)

	params2 := NewParams(1000, 100, 100, "nund", 200, 300)
	state2 := NewGenesisState(params2, 1, nil)

	require.Equal(t, state1, state2)
}

func TestDefaultGenesisState(t *testing.T) {
	state1 := DefaultGenesisState()
	state2 := DefaultGenesisState()

	require.Equal(t, state1, state2)
}

func TestValidateGenesis(t *testing.T) {
	state1 := DefaultGenesisState()
	err := ValidateGenesis(*state1)
	require.NoError(t, err)

	state2 := GenesisState{}
	err = ValidateGenesis(state2)
	require.Error(t, err)

	state3 := DefaultGenesisState()
	wrkchain1 := WrkChainExport{
		Wrkchain: WrkChain{
			WrkchainId: 0,
		},
	}

	state3.RegisteredWrkchains = append(state3.RegisteredWrkchains, wrkchain1)

	expectedErr := fmt.Errorf("invalid WrkChain: ID: %d. Error: Missing ID", 0)
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].Wrkchain.WrkchainId = 1
	expectedErr = fmt.Errorf("invalid WrkChain: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	bOwnerAddr := sdk.AccAddress(pubKey.Address())
	state3.RegisteredWrkchains[0].Wrkchain.Owner = bOwnerAddr.String()

	expectedErr = fmt.Errorf("invalid Beacon: Moniker: . Error: Missing Moniker")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].Wrkchain.Moniker = "wrkchain"
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	expectedErr = fmt.Errorf("invalid Beacon: BaseType: . Error: Missing BaseType")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].Wrkchain.Type = "tendermint"
	expectedErr = fmt.Errorf("invalid Beacon: InStateLimit: 0. Error: Missing InStateLimit")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].InStateLimit = 100
	err = ValidateGenesis(*state3)
	require.NoError(t, err)

	block := WrkChainBlockGenesisExport{}
	state3.RegisteredWrkchains[0].Blocks = append(state3.RegisteredWrkchains[0].Blocks, block)

	expectedErr = fmt.Errorf("invalid WrkChain block: WrkChainID: 0. Error: Missing WrkChainID")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].Blocks[0].Bh = "ljbhouhgygiuyiug"
	expectedErr = fmt.Errorf("invalid WrkChain block: Height: . Error: Missing Height")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredWrkchains[0].Blocks[0].He = 12345
	err = ValidateGenesis(*state3)
	require.NoError(t, err)
}
