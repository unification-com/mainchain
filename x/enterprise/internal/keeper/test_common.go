package keeper

import (
	"bytes"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
	"math/rand"
	"testing"

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

const TestDenomination = "testc"

// dummy addresses used for testing
var (
	entSrcPk   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB53")
	entPk1     = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51")
	entPk2     = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	entPk3     = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52")
	entPk4     = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB54")
	entPk5     = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB55")
	EntSrcAddr = sdk.AccAddress(entSrcPk.Address())
	entAddr1   = sdk.AccAddress(entPk1.Address())
	entAddr2   = sdk.AccAddress(entPk2.Address())
	entAddr3   = sdk.AccAddress(entPk3.Address())
	entAddr4   = sdk.AccAddress(entPk4.Address())
	entAddr5   = sdk.AccAddress(entPk5.Address())

	TestAddrs = []sdk.AccAddress{
		entAddr1, entAddr2, entAddr3, entAddr4, entAddr5,
	}
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

func createTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, auth.AccountKeeper, Keeper, staking.Keeper, types.SupplyKeeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyEnt := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyEnt, sdk.StoreTypeIAVL, db)
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
	entAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter, supply.Staking)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[entAcc.GetAddress().String()] = true
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
		keyEnt, supplyKeeper, accountKeeper, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace, cdc,
	)

	keeper.SetHighestPurchaseOrderID(ctx, types.DefaultStartingPurchaseOrderID)
	entParams := types.DefaultParams()
	entParams.EntSource = EntSrcAddr
	entParams.Denom = stakingKeeper.BondDenom(ctx)
	keeper.SetParams(ctx, entParams)

	initCoins := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	_, err := bankKeeper.AddCoins(ctx, EntSrcAddr, initCoins)
	require.Nil(t, err)

	keeper.supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, entAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	return ctx, accountKeeper, keeper, stakingKeeper, supplyKeeper
}

// PurchaseOrderEqual checks if two purchase orders are equal
func PurchaseOrderEqual(poA types.EnterpriseUndPurchaseOrder, poB types.EnterpriseUndPurchaseOrder) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(poA),
		types.ModuleCdc.MustMarshalBinaryBare(poB))
}

func ParamsEqual(paramsA, paramsB types.Params) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(paramsA),
		types.ModuleCdc.MustMarshalBinaryBare(paramsB))
}

func LockedUndEqual(lA, lB types.LockedUnd) bool {
	return bytes.Equal(types.ModuleCdc.MustMarshalBinaryBare(lA),
		types.ModuleCdc.MustMarshalBinaryBare(lB))
}

func RandomDecision() types.PurchaseOrderStatus {
	rnd := rand.Intn(100)
	if rnd >= 50 {
		return types.StatusAccepted
	}
	return types.StatusRejected
}

func RandomStatus() types.PurchaseOrderStatus {
	rnd := RandInBetween(1, 4)
	switch rnd {
	case 1:
		return types.StatusRaised
	case 2:
		return types.StatusAccepted
	case 3:
		return types.StatusRejected
	case 4:
		return types.StatusCompleted
	default:
		return types.StatusRaised
	}
}

func RandInBetween(min, max int) int {
	return rand.Intn(max-min) + min
}
