package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

//// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
//type FeeTx interface {
//	sdk.Tx
//	GetGas() uint64
//	GetFee() sdk.Coins
//	FeePayer() sdk.AccAddress
//}

const (
	RouterKey      = types.RouterKey
	RegisterAction = types.RegisterAction
	RecordAction   = types.RecordAction
)

var (
	ErrIncorrectFeeDenomination = types.ErrIncorrectFeeDenomination
	ErrInsufficientWrkChainFee  = types.ErrInsufficientWrkChainFee
	ErrTooMuchWrkChainFee       = types.ErrTooMuchWrkChainFee
)

func CheckIsWrkChainTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		if msg.Route() == types.RouterKey {
			switch msg.Type() {
			case types.RecordAction:
				return true
			case types.RegisterAction:
				return true
			}
		}
	}
	return false
}
