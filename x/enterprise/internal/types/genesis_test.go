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

	state1.StartingPurchaseOrderID = 1
	require.NotEqual(t, state1, state2)
	require.False(t, state1.Equal(state2))

	state2.StartingPurchaseOrderID = 1
	require.Equal(t, state1, state2)
	require.True(t, state1.Equal(state2))
}

func TestDefaultGenesisState(t *testing.T) {
	state1 := DefaultGenesisState()
	state2 := DefaultGenesisState()

	require.Equal(t, state1, state2)
}

func TestValidateGenesis(t *testing.T) {
	privK2 := ed25519.GenPrivKey()
	pubKey2 := privK2.PubKey()
	signerAddr := sdk.AccAddress(pubKey2.Address())

	state1 := DefaultGenesisState()
	state1.Params.EntSigners = signerAddr.String()
	err := ValidateGenesis(state1)
	require.NoError(t, err)

	state2 := GenesisState{}
	expectedErr := fmt.Errorf("denom cannot be blank")
	err = ValidateGenesis(state2)
	require.Equal(t, expectedErr.Error(), err.Error())

	state2.Params.Denom = "nund"
	expectedErr = fmt.Errorf("min accepts must be positive: 0")
	err = ValidateGenesis(state2)
	require.Equal(t, expectedErr.Error(), err.Error())

	state2.Params.MinAccepts = 1
	expectedErr = fmt.Errorf("decision limit must be positive: 0")
	err = ValidateGenesis(state2)
	require.Equal(t, expectedErr.Error(), err.Error())

	state2.Params.DecisionLimit = 1
	expectedErr = fmt.Errorf("must have at least one signer")
	err = ValidateGenesis(state2)
	require.Equal(t, expectedErr.Error(), err.Error())

	state2.Params.EntSigners = signerAddr.String()
	expectedErr = fmt.Errorf("enterprise starting purchase order id should be greater than 0")
	err = ValidateGenesis(state2)
	require.Equal(t, expectedErr.Error(), err.Error())

	state2.StartingPurchaseOrderID = 1
	err = ValidateGenesis(state2)
	require.NoError(t, err)

	state3 := DefaultGenesisState()
	state3.Params.EntSigners = signerAddr.String()
	po1 := EnterpriseUndPurchaseOrder{
		PurchaseOrderID: 0,
	}

	state3.PurchaseOrders = append(state3.PurchaseOrders, po1)
	expectedErr = fmt.Errorf("invalid purchase order: PurchaseOrderID: 0. Error: Missing PurchaseOrderID")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].PurchaseOrderID = 1
	expectedErr = fmt.Errorf("invalid purchase order: Purchaser: . Error: Missing Purchaser")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	privK := ed25519.GenPrivKey()
	pubKey := privK.PubKey()
	purchaserAddr := sdk.AccAddress(pubKey.Address())

	state3.PurchaseOrders[0].Purchaser = purchaserAddr
	expectedErr = fmt.Errorf("invalid purchase order: Amount: %s. Error: Missing Amount", po1.Amount)
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Amount = sdk.NewInt64Coin("nund", 0)
	expectedErr = fmt.Errorf("invalid purchase order: Amount. Error: Amount must be greater than 0")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Amount = sdk.NewInt64Coin("nund", 100)
	expectedErr = fmt.Errorf("invalid purchase order: Status: . Error: Invalid Status")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	status, _ := PurchaseOrderStatusFromString("something")
	state3.PurchaseOrders[0].Status = status
	expectedErr = fmt.Errorf("invalid purchase order: Status: . Error: Invalid Status")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Status = StatusNil
	expectedErr = fmt.Errorf("invalid purchase order: Status: . Error: Invalid Status")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Status = StatusRaised
	err = ValidateGenesis(state3)
	require.NoError(t, err)

	decision := PurchaseOrderDecision{}
	state3.PurchaseOrders[0].Decisions = append(state3.PurchaseOrders[0].Decisions, decision)
	expectedErr = fmt.Errorf("invalid decision: Signer cannot be empty")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Decisions[0].Signer = signerAddr
	expectedErr = fmt.Errorf("invalid decision: Decision: . Error: Invalid Decision")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

	state3.PurchaseOrders[0].Decisions[0].Decision = StatusRaised
	expectedErr = fmt.Errorf("invalid decision: Decision: raised. Error: Invalid Decision")
	err = ValidateGenesis(state3)
	require.Equal(t, expectedErr.Error(), err.Error())

}
