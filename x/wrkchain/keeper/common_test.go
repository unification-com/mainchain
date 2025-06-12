package keeper_test

import (
	"math/rand"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	charsetForRand = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

var (
	TestAddrs  = createRandomAccounts(10)
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// WRKChainEqual checks if two WRKChains are equal
func WRKChainEqual(wcA types.WrkChain, wcB types.WrkChain) bool {
	return wcA == wcB
}

// ParamsEqual checks params are equal
func ParamsEqual(paramsA, paramsB types.Params) bool {
	return paramsA == paramsB
}

// WRKChainBlockEqual checks if two WRKChainBlocks are equal
func WRKChainBlockEqual(bA, bB types.WrkChainBlock) bool {
	return bA == bB
}

// WRKChainBlockLegacyEqual checks if two WRKChainBlocks are equal
func WRKChainBlockLegacyEqual(bA, bB types.WrkChainBlockLegacy) bool {
	return bA == bB
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
	return GenerateRandomStringWithCharset(length, charsetForRand)
}
