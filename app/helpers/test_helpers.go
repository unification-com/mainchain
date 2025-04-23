package helpers

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	undapp "github.com/unification-com/mainchain/app"
	appparams "github.com/unification-com/mainchain/app/params"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	streamtypes "github.com/unification-com/mainchain/x/stream/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"
)

// SimAppChainID hardcoded chainID for simulation
const (
	SimAppChainID             = "fund-app"
	SimTestRegFee             = 10000000000
	SimTestRecordFee          = 1000000000
	SimTestPurchaseStorageFee = 1000000000

	SimTestDefaultStartingId uint64 = 1

	SimTestDefaultStorageLimit    uint64 = 1000
	SimTestDefaultMaxStorageLimit uint64 = 10000
)

var SimTestDefaultStreamValFee = sdkmath.LegacyNewDecWithPrec(1, 2)

// DefaultConsensusParams defines the default Tendermint consensus params used
// in Und App testing.
var DefaultConsensusParams = &tmproto.ConsensusParams{
	Block: &tmproto.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

type PV struct {
	PrivKey cryptotypes.PrivKey
}

type EmptyAppOptions struct{}

func (EmptyAppOptions) Get(_ string) interface{} { return nil }

func Setup(t *testing.T) *undapp.App {
	appparams.SetAddressPrefixes()
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := mock.NewPV()
	senderPubKey := senderPrivKey.PrivKey.PubKey()

	acc := authtypes.NewBaseAccount(senderPubKey.Address().Bytes(), senderPubKey, 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000))),
	}
	genesisAccounts := []authtypes.GenesisAccount{acc}
	app := SetupWithGenesisValSet(t, valSet, genesisAccounts, balance)

	return app
}

// SetupWithGenesisValSet initializes a new App with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the App from first genesis
// account. A Nop logger is set in App.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *undapp.App {
	t.Helper()

	fundApp, genesisState := setup()
	genesisState = genesisStateWithValSet(t, fundApp, genesisState, valSet, genAccs, balances...)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = fundApp.InitChain(
		&abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	_, err = fundApp.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             fundApp.LastBlockHeight() + 1,
		Hash:               fundApp.LastCommitID().Hash,
		NextValidatorsHash: valSet.Hash(),
	})
	require.NoError(t, err)

	return fundApp
}

func setup() (*undapp.App, undapp.GenesisState) {
	db := dbm.NewMemDB()

	dir, err := os.MkdirTemp("", "und-test-app")
	if err != nil {
		panic(err)
	}

	defer func() {
		os.RemoveAll(dir)
	}()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[server.FlagInvCheckPeriod] = 5
	appOptions[server.FlagMinGasPrices] = "0.25nund"
	appOptions[flags.FlagHome] = dir

	fundApp := undapp.NewApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appOptions,
	)
	return fundApp, fundApp.BasicModuleManager.DefaultGenesis(fundApp.AppCodec())
}

func genesisStateWithValSet(t *testing.T,
	app *undapp.App, genesisState undapp.GenesisState,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) undapp.GenesisState {
	t.Helper()
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdkmath.LegacyOneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			Commission:      stakingtypes.NewCommission(sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec()),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), sdkmath.LegacyOneDec()))

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	// Beacon params
	beaconGenesis := beacontypes.NewGenesisState(
		beacontypes.NewParams(SimTestRegFee, SimTestRecordFee, SimTestPurchaseStorageFee, sdk.DefaultBondDenom, SimTestDefaultStorageLimit, SimTestDefaultMaxStorageLimit),
		SimTestDefaultStartingId, []beacontypes.BeaconExport{},
	)
	genesisState[beacontypes.ModuleName] = app.AppCodec().MustMarshalJSON(beaconGenesis)

	// WrkChain
	wrkChainGenesis := wrkchaintypes.NewGenesisState(
		wrkchaintypes.NewParams(SimTestRegFee, SimTestRecordFee, SimTestPurchaseStorageFee, sdk.DefaultBondDenom, SimTestDefaultStorageLimit, SimTestDefaultMaxStorageLimit),
		SimTestDefaultStartingId, []wrkchaintypes.WrkChainExport{},
	)
	genesisState[wrkchaintypes.ModuleName] = app.AppCodec().MustMarshalJSON(wrkChainGenesis)

	// Enterprise
	enterpriseGenesis := enttypes.NewGenesisState(
		enttypes.NewParams(sdk.DefaultBondDenom, 1, 1000, genAccs[0].GetAddress().String()),
		1, sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewIntFromUint64(0)),
		enttypes.EnterpriseUndPurchaseOrders{}, enttypes.LockedUnds{}, enttypes.Whitelists{},
		sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewIntFromUint64(0)), enttypes.SpentEFUNDs{},
	)
	genesisState[enttypes.ModuleName] = app.AppCodec().MustMarshalJSON(enterpriseGenesis)

	// Streams
	streamGenesis := streamtypes.NewGenesisState(
		[]streamtypes.StreamExport{},
		streamtypes.NewParams(SimTestDefaultStreamValFee),
	)
	genesisState[streamtypes.ModuleName] = app.AppCodec().MustMarshalJSON(streamGenesis)

	return genesisState
}
