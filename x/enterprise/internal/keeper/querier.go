package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
	"strconv"
)

const (
	QueryParameters       = "params"
	QueryPurchaseOrders   = "orders"
	QueryGetPurchaseOrder = "order"
	QueryGetLocked        = "locked"
	QueryTotalLocked      = "total-locked"
	QueryTotalUnlocked    = "total-unlocked"
	QueryTotalSupply      = "total-supply"
	QueryWhitelist        = "whitelist"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, keeper)
		case QueryPurchaseOrders:
			return queryPurchaseOrders(ctx, path[1:], req, keeper)
		case QueryGetPurchaseOrder:
			return queryPurchaseOrderById(ctx, path[1:], keeper)
		case QueryGetLocked:
			return queryLockedUndByAddress(ctx, path[1:], keeper)
		case QueryTotalLocked:
			return queryTotalLocked(ctx, keeper)
		case QueryTotalUnlocked:
			return queryTotalUnlocked(ctx, keeper)
		case QueryTotalSupply:
			return queryTotalSupply(ctx, keeper)
		case QueryWhitelist:
			return queryWhitelist(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPurchaseOrders(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {

	var queryParams types.QueryPurchaseOrdersParams

	err := k.cdc.UnmarshalJSON(req.Data, &queryParams)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	filteredPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)

	if filteredPurchaseOrders == nil {
		filteredPurchaseOrders = types.PurchaseOrders{}
	}

	res, err := codec.MarshalJSONIndent(k.cdc, filteredPurchaseOrders)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPurchaseOrderById(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	purchaseOrderId, err := strconv.Atoi(path[0])

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	purchaseOrder := k.GetPurchaseOrder(ctx, uint64(purchaseOrderId))

	res, err := codec.MarshalJSONIndent(k.cdc, purchaseOrder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryLockedUndByAddress(ctx sdk.Context, path []string, k Keeper) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	lockedUnd := k.GetLockedUndForAccount(ctx, address)

	res, err := codec.MarshalJSONIndent(k.cdc, lockedUnd)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTotalLocked(ctx sdk.Context, k Keeper) ([]byte, error) {

	totalLocked := k.GetTotalLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, totalLocked)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTotalUnlocked(ctx sdk.Context, k Keeper) ([]byte, error) {

	totalUnlocked := k.GetTotalUnLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, totalUnlocked)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTotalSupply(ctx sdk.Context, k Keeper) ([]byte, error) {
	totalSupply := k.GetTotalSupplyIncludingLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, totalSupply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryWhitelist(ctx sdk.Context, k Keeper) ([]byte, error) {
	whitelist := k.GetAllWhitelistedAddresses(ctx)

	res, err := codec.MarshalJSONIndent(k.cdc, whitelist)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
