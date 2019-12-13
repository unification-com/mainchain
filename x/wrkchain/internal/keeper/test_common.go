package keeper

import (
	"bytes"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	TestDenomination = "testc"
)

// dummy addresses used for testing
var (
	wcPk1   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51")
	wcPk2   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	wcPk3   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52")
	wcPk4   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB54")
	wcPk5   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB55")
	wcAddr1 = sdk.AccAddress(wcPk1.Address())
	wcAddr2 = sdk.AccAddress(wcPk2.Address())
	wcAddr3 = sdk.AccAddress(wcPk3.Address())
	wcAddr4 = sdk.AccAddress(wcPk4.Address())
	wcAddr5 = sdk.AccAddress(wcPk5.Address())

	TestAddrs = []sdk.AccAddress{
		wcAddr1, wcAddr2, wcAddr3, wcAddr4, wcAddr5,
	}

	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateRandomAddresses(num int) []sdk.AccAddress {
	var testAddrs []sdk.AccAddress
	for i := 0; i < num; i++ {
		privK := ed25519.GenPrivKey()
		pubKey := privK.PubKey()
		testAddrs = append(testAddrs, sdk.AccAddress(pubKey.Address()))
	}
	return testAddrs
}

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func makeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func createTestInput(t *testing.T, isCheckTx bool, initPower int64, genAccs int) (sdk.Context, auth.AccountKeeper, Keeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyWRKChain := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyWRKChain, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	require.Nil(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "und-unit-test-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
	}

	// create module accounts
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	wrkchainAcc := supply.NewEmptyModuleAccount(types.ModuleName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[wrkchainAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)

	stakingKeeper := staking.NewKeeper(cdc, keyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	skParams := staking.DefaultParams()
	skParams.BondDenom = TestDenomination
	stakingKeeper.SetParams(ctx, skParams)

	keeper := NewKeeper(
		keyWRKChain, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace, cdc,
	)

	keeper.SetHighestWrkChainID(ctx, types.DefaultStartingWrkChainID)
	wrkchainParams := types.DefaultParams()
	wrkchainParams.FeeRegister = 10
	wrkchainParams.FeeRecord = 1
	wrkchainParams.Denom = stakingKeeper.BondDenom(ctx)
	keeper.SetParams(ctx, wrkchainParams)

	if genAccs > 0 {
		TestAddrs = GenerateRandomAddresses(genAccs)
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	return ctx, accountKeeper, keeper
}

// WRKChainEqual checks if two WRKChains are equal
func WRKChainEqual(wcA types.WrkChain, wcB types.WrkChain) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(wcA),
		types.ModuleCdc.MustMarshalBinaryBare(wcB))
}

// ParamsEqual checks params are equal
func ParamsEqual(paramsA, paramsB types.Params) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(paramsA),
		types.ModuleCdc.MustMarshalBinaryBare(paramsB))
}

// WRKChainBlockEqual checks if two WRKChainBlocks are equal
func WRKChainBlockEqual(lA, lB types.WrkChainBlock) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(lA),
		types.ModuleCdc.MustMarshalBinaryBare(lB))
}

// RandInBetween generates a random number between two given values
func RandInBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

// GenerateRandomStringWithCharset generates a random string given a length and character set
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
