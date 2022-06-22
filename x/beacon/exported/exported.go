package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/beacon/types"
)

const (
	RouterKey      = types.RouterKey
	RegisterAction = types.RegisterAction
	RecordAction   = types.RecordAction
)

var (
	ErrIncorrectFeeDenomination = types.ErrIncorrectFeeDenomination
	ErrInsufficientBeaconFee    = types.ErrInsufficientBeaconFee
	ErrTooMuchBeaconFee         = types.ErrTooMuchBeaconFee
)

func CheckIsBeaconTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.(type) {
		case *types.MsgRegisterBeacon:
			return true
		case *types.MsgRecordBeaconTimestamp:
			return true
		}
	}
	return false
}
