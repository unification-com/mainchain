package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

type AccountWithLocked struct {
	Account   exported.Account `json:"account"`
	Locked    sdk.Coin         `json:"locked"`
	Available sdk.Coins        `json:"available"`
}

func NewAccountWithLocked() AccountWithLocked {
	return AccountWithLocked{}
}

func (a AccountWithLocked) String() string {
	return ""
}
