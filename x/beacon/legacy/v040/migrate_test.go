package v040_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/app"
	v038 "github.com/unification-com/mainchain/x/beacon/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/beacon/legacy/v040"
)

func TestMigrate(t *testing.T) {
	// set correct und address prefixes
	app.SetConfig()
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithCodec(encodingConfig.Marshaler)

	owner, err := sdk.AccAddressFromBech32("und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v")
	require.NoError(t, err)

	ts1 := v038.BeaconTimestamp{
		BeaconID:    1,
		TimestampID: 1,
		SubmitTime:  1629467105,
		Hash:        "f0abc33a31f3023e8caf3c2ccedb42ae02345d30cd3c08ab4d460d02c44f0f19",
		Owner:       owner,
	}

	b1 := v038.BeaconExport{
		Beacon: v038.Beacon{
			BeaconID:        1,
			Moniker:         "test_b1",
			Name:            "Test B1",
			LastTimestampID: 1,
			Owner:           owner,
		},
		BeaconTimestamps: []v038.BeaconTimestamp{ts1},
	}

	v038State := v038.GenesisState{
		Params: v038.Params{
			FeeRegister: 1000,
			FeeRecord:   100,
			Denom:       "nund",
		},
		StartingBeaconID: 2,
		Beacons:          []v038.BeaconExport{b1},
	}

	migrated := v040.Migrate(v038State)

	bz, err := clientCtx.Codec.MarshalJSON(migrated)
	require.NoError(t, err)

	// Indent the JSON bz correctly.
	var jsonObj map[string]interface{}
	err = json.Unmarshal(bz, &jsonObj)
	require.NoError(t, err)
	indentedBz, err := json.MarshalIndent(jsonObj, "", "  ")
	require.NoError(t, err)

	expected := `{
  "params": {
    "denom": "nund",
    "fee_record": "100",
    "fee_register": "1000"
  },
  "registered_beacons": [
    {
      "beacon": {
        "beacon_id": "1",
        "first_id_in_state": "1",
        "last_timestamp_id": "1",
        "moniker": "test_b1",
        "name": "Test B1",
        "num_in_state": "1",
        "owner": "und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v",
        "reg_time": "0"
      },
      "timestamps": [
        {
          "h": "f0abc33a31f3023e8caf3c2ccedb42ae02345d30cd3c08ab4d460d02c44f0f19",
          "id": "1",
          "t": "1629467105"
        }
      ]
    }
  ],
  "starting_beacon_id": "2"
}`
	require.Equal(t, expected, string(indentedBz))
}
