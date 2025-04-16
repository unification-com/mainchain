package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/unification-com/mainchain/x/enterprise/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// Purchase Order queries PO details based on PurchaseOrderId
func (q Keeper) EnterpriseUndPurchaseOrder(c context.Context, req *types.QueryEnterpriseUndPurchaseOrderRequest) (*types.QueryEnterpriseUndPurchaseOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.PurchaseOrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "purchase order id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	purchaseOrder, found := q.GetPurchaseOrder(ctx, req.PurchaseOrderId)

	if !found {
		return nil, status.Errorf(codes.NotFound, "purchase order %d doesn't exist", req.PurchaseOrderId)
	}

	return &types.QueryEnterpriseUndPurchaseOrderResponse{PurchaseOrder: purchaseOrder}, nil
}

// Purchase Orders paginated
func (q Keeper) EnterpriseUndPurchaseOrders(c context.Context, req *types.QueryEnterpriseUndPurchaseOrdersRequest) (*types.QueryEnterpriseUndPurchaseOrdersResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(q.storeKey)
	var purchaseOrders []types.EnterpriseUndPurchaseOrder

	poStore := prefix.NewStore(store, types.PurchaseOrderIDKeyPrefix)

	pageRes, err := query.FilteredPaginate(poStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var info types.EnterpriseUndPurchaseOrder
		err := q.cdc.Unmarshal(value, &info)
		if err != nil {
			return false, err
		}

		if req.Status.String() != "STATUS_NIL" && !strings.EqualFold(info.GetStatus().String(), req.Status.String()) {
			return false, nil
		}

		if req.Purchaser != "" && !strings.EqualFold(info.Purchaser, req.Purchaser) {
			return false, nil
		}

		if accumulate {
			purchaseOrders = append(purchaseOrders, info)
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return &types.QueryEnterpriseUndPurchaseOrdersResponse{PurchaseOrders: purchaseOrders, Pagination: pageRes}, nil
}

func (q Keeper) LockedUndByAddress(c context.Context, req *types.QueryLockedUndByAddressRequest) (*types.QueryLockedUndByAddressResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Owner == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	lockedUnd := q.GetLockedUndAmountForAccount(ctx, addr)

	return &types.QueryLockedUndByAddressResponse{Amount: lockedUnd}, nil
}

func (q Keeper) TotalLocked(c context.Context, req *types.QueryTotalLockedRequest) (*types.QueryTotalLockedResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	amount := q.GetTotalLockedUnd(ctx)

	return &types.QueryTotalLockedResponse{Amount: amount}, nil
}

func (q Keeper) Whitelist(c context.Context, req *types.QueryWhitelistRequest) (*types.QueryWhitelistResponse, error) {

	ctx := sdk.UnwrapSDKContext(c)
	whitelist := q.GetAllWhitelistedAddresses(ctx)

	return &types.QueryWhitelistResponse{Addresses: whitelist}, nil
}

func (q Keeper) Whitelisted(c context.Context, req *types.QueryWhitelistedRequest) (*types.QueryWhitelistedResponse, error) {

	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	isWhilelisted := q.AddressIsWhitelisted(ctx, addr)

	return &types.QueryWhitelistedResponse{Address: req.Address, Whitelisted: isWhilelisted}, nil
}

func (q Keeper) EnterpriseAccount(c context.Context, req *types.QueryEnterpriseAccountRequest) (*types.QueryEnterpriseAccountResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	userAcc := q.GetEnterpriseUserAccount(ctx, addr)

	return &types.QueryEnterpriseAccountResponse{
		Account: userAcc,
	}, nil
}

func (q Keeper) TotalSpentEFUND(c context.Context, req *types.QueryTotalSpentEFUNDRequest) (*types.QueryTotalSpentEFUNDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	totalSpent := q.GetTotalSpentEFUND(ctx)

	return &types.QueryTotalSpentEFUNDResponse{
		Amount: totalSpent,
	}, nil
}

func (q Keeper) SpentEFUNDByAddress(c context.Context, req *types.QuerySpentEFUNDByAddressRequest) (*types.QuerySpentEFUNDByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	spent := q.GetSpentEFUNDAmountForAccount(ctx, addr)

	return &types.QuerySpentEFUNDByAddressResponse{
		Amount: spent,
	}, nil
}
