package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	SuperSecretSimAppSeed = "SuperSecretSimAppEntKeySeed"
)

func GenerateEntSourceSimAccount() simulation.Account {
	var simAppAcc simulation.Account
	simAppAcc.PrivKey = secp256k1.GenPrivKeySecp256k1([]byte(SuperSecretSimAppSeed))
	simAppAcc.PubKey = simAppAcc.PrivKey.PubKey()
	simAppAcc.Address = sdk.AccAddress(simAppAcc.PubKey.Address())
	return simAppAcc
}

func GetEntSourceAddress() sdk.AccAddress {
	return GenerateEntSourceSimAccount().Address
}
