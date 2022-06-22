package v040_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/app"
	v038 "github.com/unification-com/mainchain/x/enterprise/legacy/v038"
	v040 "github.com/unification-com/mainchain/x/enterprise/legacy/v040"
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

	es1 := "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6"

	u1 := "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq"
	u2 := "und173qnkw458p646fahmd53xa45vqqvga7kyu6ryy"
	u3 := "und10kx65ezcenza0n5ex7r7pgltdnv2932rwhsmfw"
	u4 := "und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v"

	p1, err := sdk.AccAddressFromBech32(u1)
	require.NoError(t, err)
	p2, err := sdk.AccAddressFromBech32(u2)
	require.NoError(t, err)
	p3, err := sdk.AccAddressFromBech32(u3)
	require.NoError(t, err)
	p4, err := sdk.AccAddressFromBech32(u4)
	require.NoError(t, err)

	e1, err := sdk.AccAddressFromBech32(es1)
	require.NoError(t, err)

	locked := v038.LockedUnd{
		Owner: p1,
		Amount: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(998957000000000),
		},
	}

	po1d1 := v038.PurchaseOrderDecision{
		Signer:       e1,
		Decision:     v038.StatusAccepted,
		DecisionTime: 1629467086,
	}

	po1 := v038.EnterpriseUndPurchaseOrder{
		PurchaseOrderID: 1,
		Purchaser:       p1,
		Amount: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(1000000000000000),
		},
		Status:         v038.StatusCompleted,
		RaisedTime:     1629467080,
		Decisions:      []v038.PurchaseOrderDecision{po1d1},
		CompletionTime: 1629467091,
	}

	po2d1 := v038.PurchaseOrderDecision{
		Signer:       e1,
		Decision:     v038.StatusRejected,
		DecisionTime: 1629467086,
	}

	po2 := v038.EnterpriseUndPurchaseOrder{
		PurchaseOrderID: 2,
		Purchaser:       p2,
		Amount: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(1000000000000000),
		},
		Status:         v038.StatusRejected,
		RaisedTime:     1629467080,
		Decisions:      []v038.PurchaseOrderDecision{po2d1},
		CompletionTime: 0,
	}

	po3d1 := v038.PurchaseOrderDecision{
		Signer:       e1,
		Decision:     v038.StatusAccepted,
		DecisionTime: 1629467086,
	}

	po3 := v038.EnterpriseUndPurchaseOrder{
		PurchaseOrderID: 3,
		Purchaser:       p3,
		Amount: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(1000000000000000),
		},
		Status:         v038.StatusAccepted,
		RaisedTime:     1629467080,
		Decisions:      []v038.PurchaseOrderDecision{po3d1},
		CompletionTime: 0,
	}

	po4 := v038.EnterpriseUndPurchaseOrder{
		PurchaseOrderID: 4,
		Purchaser:       p4,
		Amount: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(1000000000000000),
		},
		Status:         v038.StatusRaised,
		RaisedTime:     1629467080,
		Decisions:      []v038.PurchaseOrderDecision{},
		CompletionTime: 0,
	}

	v038State := v038.GenesisState{
		Params: v038.Params{
			EntSigners:    "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6,und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed",
			Denom:         "nund",
			MinAccepts:    1,
			DecisionLimit: 3600,
		},
		StartingPurchaseOrderID: 5,
		PurchaseOrders:          v038.PurchaseOrders{po1, po2, po3, po4},
		LockedUnds:              v038.LockedUnds{locked},
		TotalLocked: sdk.Coin{
			Denom:  "nund",
			Amount: sdk.NewInt(998957000000000),
		},
		Whitelist: v038.WhitelistAddresses{p1, p2, p3, p4},
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
  "locked_und": [
    {
      "amount": {
        "amount": "998957000000000",
        "denom": "nund"
      },
      "owner": "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq"
    }
  ],
  "params": {
    "decision_time_limit": "3600",
    "denom": "nund",
    "ent_signers": "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6,und1eq239sgefyzm4crl85nfyvt7kw83vrna3f0eed",
    "min_accepts": "1"
  },
  "purchase_orders": [
    {
      "amount": {
        "amount": "1000000000000000",
        "denom": "nund"
      },
      "completion_time": "1629467091",
      "decisions": [
        {
          "decision": "STATUS_ACCEPTED",
          "decision_time": "1629467086",
          "signer": "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6"
        }
      ],
      "id": "1",
      "purchaser": "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
      "raise_time": "1629467080",
      "status": "STATUS_COMPLETED"
    },
    {
      "amount": {
        "amount": "1000000000000000",
        "denom": "nund"
      },
      "completion_time": "0",
      "decisions": [
        {
          "decision": "STATUS_REJECTED",
          "decision_time": "1629467086",
          "signer": "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6"
        }
      ],
      "id": "2",
      "purchaser": "und173qnkw458p646fahmd53xa45vqqvga7kyu6ryy",
      "raise_time": "1629467080",
      "status": "STATUS_REJECTED"
    },
    {
      "amount": {
        "amount": "1000000000000000",
        "denom": "nund"
      },
      "completion_time": "0",
      "decisions": [
        {
          "decision": "STATUS_ACCEPTED",
          "decision_time": "1629467086",
          "signer": "und1djn9sr7vtghtarp5ccvtrc0mwg9dlzjrj7alw6"
        }
      ],
      "id": "3",
      "purchaser": "und10kx65ezcenza0n5ex7r7pgltdnv2932rwhsmfw",
      "raise_time": "1629467080",
      "status": "STATUS_ACCEPTED"
    },
    {
      "amount": {
        "amount": "1000000000000000",
        "denom": "nund"
      },
      "completion_time": "0",
      "decisions": [],
      "id": "4",
      "purchaser": "und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v",
      "raise_time": "1629467080",
      "status": "STATUS_RAISED"
    }
  ],
  "starting_purchase_order_id": "5",
  "total_locked": {
    "amount": "998957000000000",
    "denom": "nund"
  },
  "whitelist": [
    "und100aex49fh53r7mpeghdq6e645epp6r9qyqk5jq",
    "und173qnkw458p646fahmd53xa45vqqvga7kyu6ryy",
    "und10kx65ezcenza0n5ex7r7pgltdnv2932rwhsmfw",
    "und10wl769hge8nhszv70uxc9zu0lgrc2lggkhst8v"
  ]
}`
	require.Equal(t, expected, string(indentedBz))
}
