package simulation

import (
	"math/rand"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/unification-com/mainchain/x/enterprise/types"
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
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, accs []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	minAccepts := uint64(simtypes.RandIntBetween(r, 1, 3))
	entAddresses := make([]string, minAccepts)

	i := uint64(0)
	for i < minAccepts {
		randAcc, _ := simtypes.RandomAcc(r, accs)
		for contains(entAddresses, randAcc.Address.String()) {
			randAcc, _ = simtypes.RandomAcc(r, accs)
		}
		entAddresses[i] = randAcc.Address.String()
		i += 1
	}

	params := types.DefaultParams()
	params.Denom = sdk.DefaultBondDenom
	params.DecisionTimeLimit = uint64(simtypes.RandIntBetween(r, 360, 3600))
	params.MinAccepts = uint64(simtypes.RandIntBetween(r, 1, 2))
	params.EntSigners = strings.Join(entAddresses, ",")

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
