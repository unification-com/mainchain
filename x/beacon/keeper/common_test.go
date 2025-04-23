package keeper_test

import (
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/unification-com/mainchain/x/beacon/types"
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

func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKey
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

// BeaconEqual checks if two Beacons are equal
func BeaconEqual(wcA types.Beacon, wcB types.Beacon) bool {
	return wcA == wcB
}

// // ParamsEqual checks params are equal
func ParamsEqual(paramsA, paramsB types.Params) bool {
	return paramsA == paramsB
}

// BeaconTimestampEqual checks if two BeaconTimestamps are equal
func BeaconTimestampEqual(lA, lB types.BeaconTimestamp) bool {
	return lA == lB
}

// BeaconTimestampLegacyEqual checks if two BeaconTimestampLegacy are equal
func BeaconTimestampLegacyEqual(lA, lB types.BeaconTimestampLegacy) bool {
	return lA == lB
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
