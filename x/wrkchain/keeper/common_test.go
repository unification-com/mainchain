package keeper_test

import (
	"math/rand"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	charsetForRand = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

var (
	TestAddrs  = simapphelpers.GenerateRandomTestAccounts(10)
	seededRand = rand.New(rand.NewSource(1))
)

// WRKChainEqual checks if two WRKChains are equal
func WRKChainEqual(wcA types.WrkChain, wcB types.WrkChain) bool {
	return wcA == wcB
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
