package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the wrkchain Querier
const (
	QueryWrkChain            = "get"
	QueryWrkChains           = "wrkchains"
	QueryWrkChainBlock       = "get-block"
	QueryWrkChainBlockHashes = "blocks"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryWrkChain:
			return queryWrkChain(ctx, path[1:], req, keeper)
		case QueryWrkChains:
			return queryWrkChainsFiltered(ctx, path[1:], req, keeper)
		case QueryWrkChainBlock:
			return queryWrkChainBlock(ctx, path[1:], req, keeper)
		case QueryWrkChainBlockHashes:
			return queryWrkChainBlockHashes(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown wrkchain query endpoint")
		}
	}
}

// nolint: unparam
func queryWrkChain(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	wrkchainID, err := strconv.Atoi(path[0])

	if err != nil {
		wrkchainID = 0
	}

	wrkchain := keeper.GetWrkChain(ctx, uint64(wrkchainID))

	res, err := codec.MarshalJSONIndent(keeper.cdc, wrkchain)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryWrkChainBlock(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	wrkchainID, err := strconv.Atoi(path[0])

	if err != nil {
		wrkchainID = 0
	}

	height, err := strconv.Atoi(path[1])

	if err != nil {
		height = 0
	}

	wrkchainBlock := keeper.GetWrkChainBlock(ctx, uint64(wrkchainID), uint64(height))

	res, err := codec.MarshalJSONIndent(keeper.cdc, wrkchainBlock)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryWrkChainBlockHashes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	wrkchainID, err := strconv.Atoi(path[0])

	if err != nil {
		wrkchainID = 0
	}

	blockHashList := keeper.GetWrkChainBlockHashes(ctx, uint64(wrkchainID))

	res, err := codec.MarshalJSONIndent(keeper.cdc, blockHashList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryWrkChainsFiltered(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {

	var queryParams types.QueryWrkChainParams

	err := k.cdc.UnmarshalJSON(req.Data, &queryParams)

	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("failed to parse params", err.Error()))
	}

	filteredWrkChains := k.GetWrkChainsFiltered(ctx, queryParams)

	if filteredWrkChains == nil {
		filteredWrkChains = types.WrkChains{}
	}

	res, err := codec.MarshalJSONIndent(k.cdc, filteredWrkChains)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}
