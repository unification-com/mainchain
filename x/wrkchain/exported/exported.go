package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	RouterKey      = types.RouterKey
	RegisterAction = types.RegisterAction
	RecordAction   = types.RecordAction
)

var (
	ErrIncorrectFeeDenomination = types.ErrIncorrectFeeDenomination
	ErrInsufficientWrkChainFee  = types.ErrInsufficientWrkChainFee
	ErrTooMuchWrkChainFee       = types.ErrTooMuchWrkChainFee
	ErrExceedsMaxStorage        = types.ErrExceedsMaxStorage
)

func CheckIsWrkChainTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgRegisterWrkChain:
			return true
		case *types.MsgRecordWrkChainBlock:
			return true
		case *types.MsgPurchaseWrkChainStateStorage:
			return true
		}

	}
	return false
}
