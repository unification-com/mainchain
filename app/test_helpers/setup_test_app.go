package test_helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/unification-com/mainchain/app/params"

	authsign "github.com/cosmos/cosmos-sdk/x/auth/signing"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	beacontypes "github.com/unification-com/mainchain/x/beacon/types"
	enttypes "github.com/unification-com/mainchain/x/enterprise/types"
	wrkchaintypes "github.com/unification-com/mainchain/x/wrkchain/types"

	"github.com/unification-com/mainchain/app"
)

// SimAppChainID hardcoded chainID for simulation
const (
	TestDenomination   = sdk.DefaultBondDenom
	SimAppChainID      = "simulation-app"
	charset            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	TestDefaultStorage = 200
	TestMaxStorage     = 300
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// app.App testing.
var (
	DefaultConsensusParams = &abci.ConsensusParams{
		Block: &abci.BlockParams{
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
	EncodingConfig = params.EncodingConfig{}
	seededRand     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	nodeHome := filepath.Join(userHomeDir, ".und_testapp")
	config := sdk.GetConfig()
	if config.GetBech32AccountAddrPrefix() == "cosmos" {
		app.SetConfig()
	}
	encCdc := app.MakeEncodingConfig()
	EncodingConfig = encCdc

	testApp := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, nodeHome, invCheckPeriod, encCdc, EmptyAppOptions{})
	if withGenesis {
		return testApp, NewDefaultGenesisState(encCdc.Marshaler)
	}
	return testApp, map[string]json.RawMessage{}
}

// Setup initializes a new app.App. A Nop logger is set in app.App.
func Setup(isCheckTx bool) *app.App {
	testApp, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		testApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return testApp
}

func SetKeeperTestParamsAndDefaultValues(app *app.App, ctx sdk.Context) {
	app.BeaconKeeper.SetParams(ctx, beacontypes.NewParams(24, 2, 2, TestDenomination, TestDefaultStorage, TestMaxStorage))
	app.WrkchainKeeper.SetParams(ctx, wrkchaintypes.NewParams(24, 2, 2, TestDenomination, TestDefaultStorage, TestMaxStorage))
	app.EnterpriseKeeper.SetParams(ctx, enttypes.Params{
		EntSigners:        sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
		Denom:             TestDenomination,
		MinAccepts:        1,
		DecisionTimeLimit: 1000,
	})
	_ = app.EnterpriseKeeper.SetTotalLockedUnd(ctx, sdk.NewInt64Coin(TestDenomination, 0))
	_ = app.EnterpriseKeeper.SetTotalSpentEFUND(ctx, sdk.NewInt64Coin(TestDenomination, 0))
	app.EnterpriseKeeper.SetHighestPurchaseOrderID(ctx, 1)
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func GenerateRandomTestAccounts(accNum int) []sdk.AccAddress {
	return createRandomAccounts(accNum)
}

func GenerateRandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomString generates a random string given a length, based on a set character set
func GenerateRandomString(length int) string {
	return GenerateRandomStringWithCharset(length, charset)
}

func RandInBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

func fundAccount(app *app.App, ctx sdk.Context, addr sdk.AccAddress, amtCoins sdk.Coins) error {
	err := app.BankKeeper.MintCoins(ctx, enttypes.ModuleName, amtCoins)
	if err != nil {
		return err
	}
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, enttypes.ModuleName, addr, amtCoins)
	if err != nil {
		return err
	}
	return nil
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *app.App, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *app.App, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

func addTestAddrs(app *app.App, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// saveAccount saves the provided account into the app.App with balance based on initCoins.
func saveAccount(app *app.App, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := fundAccount(app, ctx, addr, initCoins)
	if err != nil {
		panic(err)
	}
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) app.GenesisState {
	return app.ModuleBasics.DefaultGenesis(cdc)
}

// GenTx generates a signed mock transaction.
func GenTx(gen client.TxConfig, msgs []sdk.Msg, feeAmt sdk.Coins, gas uint64, chainID string, accNums, accSeqs []uint64, priv ...cryptotypes.PrivKey) (sdk.Tx, error) {
	sigs := make([]signing.SignatureV2, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))

	signMode := gen.SignModeHandler().DefaultMode()

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range priv {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
			},
			Sequence: accSeqs[i],
		}
	}

	tx := gen.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = tx.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range priv {
		signerData := authsign.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
		if err != nil {
			panic(err)
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = tx.SetSignatures(sigs...)
		if err != nil {
			panic(err)
		}
	}

	return tx.GetTx(), nil
}
