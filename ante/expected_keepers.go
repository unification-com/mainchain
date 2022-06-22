package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}
