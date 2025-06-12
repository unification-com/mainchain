package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/unification-com/mainchain/x/wrkchain/types"
)

const (
	DefaultWeightMsgUpdateParams int = 100

	OpWeightMsgUpdateParams = "op_weight_msg_update_params" //nolint:gosec
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams,
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()
	params.FeeRegister = uint64(simtypes.RandIntBetween(r, 100000000000, 1000000000000))
	params.FeeRecord = uint64(simtypes.RandIntBetween(r, 1000000000, 10000000000))
	params.FeePurchaseStorage = uint64(simtypes.RandIntBetween(r, 5000000000, 10000000000))
	params.DefaultStorageLimit = uint64(simtypes.RandIntBetween(r, 50000, 60000))
	params.MaxStorageLimit = uint64(simtypes.RandIntBetween(r, 600000, 700000))

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
