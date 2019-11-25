package beacon

import (
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/beacon/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = types.DefaultParamspace

	DefaultCodespace = types.DefaultCodespace

	QueryParameters = keeper.QueryParameters
)

var (
	NewKeeper       = keeper.NewKeeper
	NewQuerier      = keeper.NewQuerier
	NewGenesisState = types.NewGenesisState

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)

type (
	Keeper = keeper.Keeper

	GenesisState = types.GenesisState
)
