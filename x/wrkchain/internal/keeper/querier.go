package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the wrkchain Querier
const (
	QueryWrkChain            = "get"
	QueryWrkChainBlock       = "get-block"
	QueryWrkChainBlockHashes = "blocks"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryWrkChain:
			return queryWrkChain(ctx, path[1:], req, keeper)
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
	wrkchain := keeper.GetWrkChain(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, wrkchain)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryWrkChainBlock(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	height, err := strconv.Atoi(path[1])

	if err != nil {
		height = 0
	}

	wrkchainBlock := keeper.GetWrkChainBlock(ctx, path[0], uint64(height))

	res, err := codec.MarshalJSONIndent(keeper.cdc, wrkchainBlock)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryWrkChainBlockHashes(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	blockHashList := keeper.GetWrkChainBlockHashes(ctx, path[0])

	res, err := codec.MarshalJSONIndent(keeper.cdc, blockHashList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
