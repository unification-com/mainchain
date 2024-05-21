package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

type (
	ParamSet = paramtypes.ParamSet

	// Subspace defines an interface that implements the legacy x/params Subspace
	// type.
	//
	// NOTE: This is used solely for migration of x/params managed parameters.
	Subspace interface {
		GetParamSet(ctx sdk.Context, ps ParamSet)
	}
)
