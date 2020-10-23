package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestEqualStartingBeaconID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	state1.StartingBeaconID = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	state2.StartingBeaconID = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
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
	beacon1 := BeaconExport{
		Beacon: Beacon{
			BeaconID: 0,
		},
	}

	state3.Beacons = append(state3.Beacons, beacon1)

	expectedErr := fmt.Errorf("invalid Beacon: ID: %d. Error: Missing ID", 0)
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.Beacons[0].Beacon.BeaconID = 1
	expectedErr = fmt.Errorf("invalid Beacon: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	bOwnerAddr := sdk.AccAddress(pubKey.Address())
	state3.Beacons[0].Beacon.Owner = bOwnerAddr

	expectedErr = fmt.Errorf("invalid Beacon: Moniker: . Error: Missing Moniker")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.Beacons[0].Beacon.Moniker = "beacon"
	err = ValidateGenesis(state3)
	require.NoError(t, err)

	timestamp := BeaconTimestampGenesisExport{}
	state3.Beacons[0].BeaconTimestamps = append(state3.Beacons[0].BeaconTimestamps, timestamp)

	expectedErr = fmt.Errorf("invalid Beacon timestamp: BeaconID: 0. Error: Missing BeaconID")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.Beacons[0].BeaconTimestamps[0].TimestampID = 1
	expectedErr = fmt.Errorf("invalid Beacon timestamp: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.Beacons[0].BeaconTimestamps[0].Hash = "ljbhouhgygiuyiug"
	expectedErr = fmt.Errorf("invalid Beacon timestamp: SubmitTime: . Error: Missing SubmitTime")
	err = ValidateGenesis(state3)
	require.Error(t, expectedErr, err.Error())

	state3.Beacons[0].BeaconTimestamps[0].SubmitTime = 12345
	err = ValidateGenesis(state3)
	require.NoError(t, err)
}
