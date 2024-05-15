package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	ModuleName     = types.ModuleName
	RouterKey      = types.RouterKey
	RegisterAction = types.RegisterAction
	RecordAction   = types.RecordAction

	PurchaseStorageFee     = types.PurchaseStorageFee
	DefaultStorageLimit    = types.DefaultStorageLimit
	DefaultMaxStorageLimit = types.DefaultMaxStorageLimit
)

var (
	ErrIncorrectFeeDenomination = types.ErrIncorrectFeeDenomination
	ErrInsufficientBeaconFee    = types.ErrInsufficientBeaconFee
	ErrTooMuchBeaconFee         = types.ErrTooMuchBeaconFee
	ErrExceedsMaxStorage        = types.ErrExceedsMaxStorage
)

func CheckIsBeaconTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgRegisterBeacon:
			return true
		case *types.MsgRecordBeaconTimestamp:
			return true
		case *types.MsgPurchaseBeaconStateStorage:
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
