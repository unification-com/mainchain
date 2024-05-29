package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/stream/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) CalculateFlowRate(ctx context.Context, request *types.QueryCalculateFlowRateRequest) (*types.QueryCalculateFlowRateResponse, error) {

	coin, err := sdk.ParseCoinNormalized(request.Coin)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if request.Duration < 1 {
		return nil, status.Error(codes.InvalidArgument, "duration cannot be zero")
	}

	if request.Period < 1 || request.Period > types.StreamPeriodYear {
		return nil, status.Error(codes.InvalidArgument, "invalid period")
	}

	totalDuration, _, flowRateInt64 := types.CalculateFlowRateForCoin(coin, request.Period, request.Duration)

	return &types.QueryCalculateFlowRateResponse{
		Coin:     coin,
		Period:   request.Period,
		Duration: request.Duration,
		Seconds:  totalDuration,
		FlowRate: flowRateInt64,
	}, nil
}
