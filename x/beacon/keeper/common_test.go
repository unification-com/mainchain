package keeper_test

import (
	"math/rand"

	simapphelpers "github.com/unification-com/mainchain/app/helpers"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	charsetForRand = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

var (
	TestAddrs  = simapphelpers.GenerateRandomTestAccounts(10)
	seededRand = rand.New(rand.NewSource(1))
)

// BeaconEqual checks if two Beacons are equal
func BeaconEqual(wcA types.Beacon, wcB types.Beacon) bool {
	return wcA == wcB
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
