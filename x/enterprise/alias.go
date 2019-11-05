package enterprise

import (
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/enterprise/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace

	DefaultCodespace = types.DefaultCodespace

	QuerierRoute          = types.QuerierRoute
	QueryParameters       = keeper.QueryParameters
	QueryPurchaseOrders   = keeper.QueryPurchaseOrders
	QueryGetPurchaseOrder = keeper.QueryGetPurchaseOrder
)

var (
	NewKeeper           = keeper.NewKeeper
	NewMsgPurchaseUnd   = types.NewMsgRaiseUndPurchaseOrder
	NewEnterpriseUnd    = types.NewEnterpriseUnd
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc

	// Events
	EventTypeRaisePurchaseOrder = types.EventTypeRaisePurchaseOrder
	AttributeValueCategory      = types.AttributeValueCategory
	AttributeKeyPurchaseOrderID = types.AttributeKeyPurchaseOrderID
	AttributeKeyPurchaser       = types.AttributeKeyPurchaser
	AttributeKeyAmount          = types.AttributeKeyAmount

	// Key functions
	GetPurchaseOrderIDBytes = types.GetPurchaseOrderIDBytes
)

type (
	Keeper = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	// Msgs
	MsgPurchaseUnd    = types.MsgRaiseUndPurchaseOrder

	// Structs
	EnterpriseUnd = types.EnterpriseUndPurchaseOrder

)
