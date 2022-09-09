package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	BeaconStartingId          = "beacon_start_id"
	BeaconFeeRegister         = "beacon_fee_register"
	BeaconFeeRecord           = "beacon_fee_record"
	BeaconFeePurchaseStorage  = "beacon_fee_purchase_storage"
	BeaconDefaultStorageLimit = "beacon_default_storage_limit"
	BeaconMaxStorageLimit     = "beacon_max_storage_limit"
)

// RandomizedGenState generates a random GenesisState for beacon module
func RandomizedGenState(simState *module.SimulationState) {

	var startId uint64
	var feeRegister uint64
	var feeRecord uint64
	var feePurchaseStorage uint64
	var defaultStorageLimit uint64
	var maxStorageLimit uint64

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconStartingId, &startId, simState.Rand,
		func(r *rand.Rand) {
			startId = uint64(simtypes.RandIntBetween(r, 1, 100))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconFeeRegister, &feeRegister, simState.Rand,
		func(r *rand.Rand) {
			feeRegister = uint64(simtypes.RandIntBetween(r, 1, 10))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconFeeRecord, &feeRecord, simState.Rand,
		func(r *rand.Rand) {
			feeRecord = uint64(simtypes.RandIntBetween(r, 1, 10))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconFeePurchaseStorage, &feePurchaseStorage, simState.Rand,
		func(r *rand.Rand) {
			feePurchaseStorage = uint64(simtypes.RandIntBetween(r, 1, 10))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconDefaultStorageLimit, &defaultStorageLimit, simState.Rand,
		func(r *rand.Rand) {
			defaultStorageLimit = uint64(simtypes.RandIntBetween(r, 5, 10))
		},
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, BeaconMaxStorageLimit, &maxStorageLimit, simState.Rand,
		func(r *rand.Rand) {
			maxStorageLimit = uint64(simtypes.RandIntBetween(r, 10, 20))
		},
	)

	// NOTE: for simulation, we're using sdk.DefaultBondDenom ("stake"), since "stake" is hard-coded
	// into the SDK's module simulation functions
	beaconGenesis := types.NewGenesisState(
		types.NewParams(feeRegister, feeRecord, feePurchaseStorage, sdk.DefaultBondDenom, defaultStorageLimit, maxStorageLimit),
		startId,
		nil,
	)

	bz, err := json.MarshalIndent(&beaconGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated beacon parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(beaconGenesis)
}
