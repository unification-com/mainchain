package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEqualStartingBeaconID(t *testing.T) {
	state1 := GenesisState{}
	state2 := GenesisState{}
	require.Equal(t, state1, state2)

	state1.StartingBeaconId = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.StartingBeaconId == state2.StartingBeaconId)

	state2.StartingBeaconId = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.StartingBeaconId == state2.StartingBeaconId)
}

func TestDefaultGenesisState(t *testing.T) {
	state1 := DefaultGenesisState()
	state2 := DefaultGenesisState()

	require.Equal(t, state1, state2)
}

func TestValidateGenesis(t *testing.T) {
	testAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	state1 := DefaultGenesisState()
	err := ValidateGenesis(*state1)
	require.NoError(t, err)

	state2 := GenesisState{}
	err = ValidateGenesis(state2)
	require.Error(t, err)

	state3 := DefaultGenesisState()
	beacon1 := BeaconExport{
		Beacon: &Beacon{
			BeaconId: 0,
		},
	}

	state3.RegisteredBeacons = append(state3.RegisteredBeacons, beacon1)

	expectedErr := fmt.Errorf("invalid Beacon: ID: %d. Error: Missing ID", 0)
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Beacon.BeaconId = 1
	expectedErr = fmt.Errorf("invalid Beacon: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Beacon.Owner = testAddr.String()

	expectedErr = fmt.Errorf("invalid Beacon: Moniker: . Error: Missing Moniker")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Beacon.Moniker = "beacon"
	err = ValidateGenesis(*state3)
	require.NoError(t, err)

	timestamp := BeaconTimestampGenesisExport{}
	state3.RegisteredBeacons[0].Timestamps = append(state3.RegisteredBeacons[0].Timestamps, timestamp)

	expectedErr = fmt.Errorf("invalid Beacon timestamp: BeaconID: 0. Error: Missing BeaconID")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Timestamps[0].Id = 1
	expectedErr = fmt.Errorf("invalid Beacon timestamp: Owner: %s. Error: Missing Owner", sdk.AccAddress{})
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Timestamps[0].H = "ljbhouhgygiuyiug"
	expectedErr = fmt.Errorf("invalid Beacon timestamp: SubmitTime: . Error: Missing SubmitTime")
	err = ValidateGenesis(*state3)
	require.Error(t, expectedErr, err.Error())

	state3.RegisteredBeacons[0].Timestamps[0].T = 12345
	err = ValidateGenesis(*state3)
	require.NoError(t, err)
}
