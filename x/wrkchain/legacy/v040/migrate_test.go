package v040_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/app"
	v038 "github.com/unification-com/mainchain/x/wrkchain/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/wrkchain/legacy/v040"
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

	blk1 := v038.WrkChainBlock{
		WrkChainID: 1,
		Height:     1,
		BlockHash:  "f0abc33a31f3023e8caf3c2ccedb42ae02345d30cd3c08ab4d460d02c44f0f19",
		ParentHash: "0x8f59efee3ab5ba2b0f81951c3a3f63e15f816e3ae6710e4f2a43b9355c551a25",
		Hash1:      "0xe231454f914940e6cacadbd52487b85e07309e1f1506481cd44ff92472ac44cb",
		Hash2:      "0x60a8237d158e369ce0ea1535fbc814d9825afa93f48e7a6638814f530a6aa53b",
		Hash3:      "0x4b7a849085bd7ec5008ec663aba113d1ab26776fbe0a65e4ee9c89559191b1c4",
		SubmitTime: 1629467105,
		Owner:      owner,
	}

	wc1 := v038.WrkChainExport{
		WrkChain: v038.WrkChain{
			WrkChainID:   1,
			Moniker:      "test_wc1",
			Name:         "Test WC1",
			GenesisHash:  "0x217b485b9eab141983512d8ec37b848cc025bc86b43bb24a157dd29f7456ae96",
			BaseType:     "geth",
			LastBlock:    1,
			NumberBlocks: 1,
			RegisterTime: 1629467000,
			Owner:        owner,
		},
		WrkChainBlocks: []v038.WrkChainBlock{blk1},
	}

	v038State := v038.GenesisState{
		Params: v038.Params{
			FeeRegister: 1000,
			FeeRecord:   100,
			Denom:       "nund",
		},
		StartingWrkChainID: 2,
		WrkChains:          []v038.WrkChainExport{wc1},
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
  "registered_wrkchains": [
    {
      "blocks": [
        {
          "bh": "f0abc33a31f3023e8caf3c2ccedb42ae02345d30cd3c08ab4d460d02c44f0f19",
          "h1": "0xe231454f914940e6cacadbd52487b85e07309e1f1506481cd44ff92472ac44cb",
          "h2": "0x60a8237d158e369ce0ea1535fbc814d9825afa93f48e7a6638814f530a6aa53b",
          "h3": "0x4b7a849085bd7ec5008ec663aba113d1ab26776fbe0a65e4ee9c89559191b1c4",
          "he": "1",
          "ph": "0x8f59efee3ab5ba2b0f81951c3a3f63e15f816e3ae6710e4f2a43b9355c551a25",
          "st": "1629467105"
        }
      ],
      "wrkchain": {
        "genesis": "0x217b485b9eab141983512d8ec37b848cc025bc86b43bb24a157dd29f7456ae96",
        "lastblock": "1",
        "lowest_height": "1",
        "moniker": "test_wc1",
        "name": "Test WC1",
        "num_blocks": "1",
        "owner": "und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v",
        "reg_time": "1629467000",
        "type": "geth",
        "wrkchain_id": "1"
      }
    }
  ],
  "starting_wrkchain_id": "2"
}`
	require.Equal(t, expected, string(indentedBz))
}
