package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper is a modified version of the default found in Cosmos SDK x/auth/types/expected_keepers.go
// and adds the GetAllBalances and SpendableCoins methods required by the beacon and wrkchain modules' ante handlers
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	GetAllBalances(ctx context.Context, address sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx context.Context, address sdk.AccAddress) sdk.Coins
}
