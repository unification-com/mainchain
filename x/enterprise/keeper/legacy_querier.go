package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"strconv"
)

const (
	QueryParameters       = "params"
	QueryPurchaseOrders   = "orders"
	QueryGetPurchaseOrder = "order"
	QueryGetLocked        = "locked"
	QueryTotalLocked      = "total-locked"
	QueryTotalUnlocked    = "total-unlocked"
	QueryEnterpriseSupply = "ent-supply"
	QueryTotalSupply      = "total-supply"
	QueryTotalSupplyOf    = "total-supply-of"
	QueryWhitelist        = "whitelist"
	QueryWhitelisted      = "whitelisted"
)

// NewLegacyQuerier is the module level router for state queries
func NewLegacyQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryParameters:
			return queryParams(ctx, keeper, legacyQuerierCdc)
		case QueryPurchaseOrders:
			return queryPurchaseOrders(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case QueryGetPurchaseOrder:
			return queryPurchaseOrderById(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryGetLocked:
			return queryLockedUndByAddress(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryTotalLocked:
			return queryTotalLocked(ctx, keeper, legacyQuerierCdc)
		case QueryTotalUnlocked:
			return queryTotalUnlocked(ctx, keeper, legacyQuerierCdc)
		case QueryEnterpriseSupply:
			return queryEnterpriseSupply(ctx, keeper, legacyQuerierCdc)
		case QueryTotalSupply:
			return queryTotalSupply(ctx, req, keeper, legacyQuerierCdc)
		case QueryTotalSupplyOf:
			return queryTotalSupplyOf(ctx, path[1:], keeper, legacyQuerierCdc)
		case QueryWhitelist:
			return queryWhitelist(ctx, keeper, legacyQuerierCdc)
		case QueryWhitelisted:
			return queryQueryWhitelisted(ctx, path[1:], keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// DONE
func queryParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryPurchaseOrders(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	var queryParams types.QueryPurchaseOrdersParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &queryParams)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	filteredPurchaseOrders := k.GetPurchaseOrdersFiltered(ctx, queryParams)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, filteredPurchaseOrders)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryPurchaseOrderById(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	purchaseOrderId, err := strconv.Atoi(path[0])

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	purchaseOrder, found := k.GetPurchaseOrder(ctx, uint64(purchaseOrderId))

	if !found {
		return nil, sdkerrors.Wrapf(types.ErrPurchaseOrderDoesNotExist, "purchase order id %d does not exist", purchaseOrderId)
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, purchaseOrder)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryLockedUndByAddress(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	lockedUnd := k.GetLockedUndForAccount(ctx, address)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, lockedUnd)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryTotalLocked(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	totalLocked := k.GetTotalLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, totalLocked)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryTotalUnlocked(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	totalUnlocked := k.GetTotalUnLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, totalUnlocked)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryEnterpriseSupply(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	entSupply := k.GetEnterpriseSupplyIncludingLockedUnd(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, entSupply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryTotalSupply(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTotalSupplyRequest

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	totalSupply, pageRes, err := k.GetTotalSupplyWithLockedNundRemoved(ctx, params.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	supplyRes := &types.QueryTotalSupplyResponse{
		Supply:     totalSupply,
		Pagination: pageRes,
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, supplyRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryTotalSupplyOf(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	denom := path[0]
	supplyOf := k.GetSupplyOfWithLockedNundRemoved(ctx, denom)
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, supplyOf)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryWhitelist(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	whitelist := k.GetAllWhitelistedAddresses(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, whitelist)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// DONE
func queryQueryWhitelisted(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	isWhiteListed := k.AddressIsWhitelisted(ctx, address)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, isWhiteListed)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
