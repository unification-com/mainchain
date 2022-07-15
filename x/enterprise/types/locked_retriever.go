package types

//import (
//	"context"
//	"fmt"
//	"strconv"
//
//	grpc "google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"
//
//	"github.com/cosmos/cosmos-sdk/client"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
//)
//
//type LockedRetriever struct{}
//
//func (ar LockedRetriever) GetLockedWithHeight(clientCtx client.Context, addr sdk.AccAddress) (*LockedUnd, int64, error) {
//	var header metadata.MD
//
//	queryClient := NewQueryClient(clientCtx)
//	//LockedUndByAddress
//	res, err := queryClient.LockedUndByAddress(context.Background(), &QueryLockedUndByAddressRequest{Owner: addr.String()}, grpc.Header(&header))
//	if err != nil {
//		return &LockedUnd{}, 0, err
//	}
//
//	blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
//	if l := len(blockHeight); l != 1 {
//		return &LockedUnd{}, 0, fmt.Errorf("unexpected '%s' header length; got %d, expected: %d", grpctypes.GRPCBlockHeightHeader, l, 1)
//	}
//
//	nBlockHeight, err := strconv.Atoi(blockHeight[0])
//	if err != nil {
//		return &LockedUnd{}, 0, fmt.Errorf("failed to parse block height: %w", err)
//	}
//
//	return res.LockedUnd, int64(nBlockHeight), nil
//}
