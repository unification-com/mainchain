package enterprise

import (
	"github.com/unification-com/mainchain/x/enterprise/internal/ante"
	"github.com/unification-com/mainchain/x/enterprise/internal/keeper"
	"github.com/unification-com/mainchain/x/enterprise/internal/types"
)

const (
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace

	QuerierRoute          = types.QuerierRoute
	QueryParameters       = keeper.QueryParameters
	QueryPurchaseOrders   = keeper.QueryPurchaseOrders
	QueryGetPurchaseOrder = keeper.QueryGetPurchaseOrder

	StatusNil       = types.StatusNil
	StatusRaised    = types.StatusRaised
	StatusAccepted  = types.StatusAccepted
	StatusRejected  = types.StatusRejected
	StatusCompleted = types.StatusCompleted
)

var (
	NewKeeper           = keeper.NewKeeper
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewQuerier          = keeper.NewQuerier
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc

	// Events
	EventTypeRaisePurchaseOrder           = types.EventTypeRaisePurchaseOrder
	EventTypeProcessPurchaseOrderDecision = types.EventTypeProcessPurchaseOrderDecision
	EventTypeTallyPurchaseOrderDecisions  = types.EventTypeTallyPurchaseOrderDecisions
	EventTypeAutoRejectStalePurchaseOrder = types.EventTypeAutoRejectStalePurchaseOrder
	EventTypeUndPurchaseComplete          = types.EventTypeUndPurchaseComplete
	EventTypeWhitelistAddress             = types.EventTypeWhitelistAddress
	AttributeValueCategory                = types.AttributeValueCategory
	AttributeKeyPurchaseOrderID           = types.AttributeKeyPurchaseOrderID
	AttributeKeyPurchaser                 = types.AttributeKeyPurchaser
	AttributeKeyAmount                    = types.AttributeKeyAmount
	AttributeKeyDecision                  = types.AttributeKeyDecision
	AttributeKeySigner                    = types.AttributeKeySigner
	AttributeKeyNumAccepts                = types.AttributeKeyNumAccepts
	AttributeKeyNumRejects                = types.AttributeKeyNumRejects
	AttributeKeyWhitelistAction           = types.AttributeKeyWhitelistAction
	AttributeWhitelistAddress             = types.AttributeWhitelistAddress

	// Key functions
	GetPurchaseOrderIDBytes = types.GetPurchaseOrderIDBytes

	ValidPurchaseOrderAcceptRejectStatus = types.ValidPurchaseOrderAcceptRejectStatus
	ValidWhitelistAction                 = types.ValidWhitelistAction

	NewLockedUndRetriever   = keeper.NewLockedUndRetriever
	NewTotalSupplyRetriever = keeper.NewTotalSupplyRetriever

	RegisterInvariants = keeper.RegisterInvariants
	AllInvariants      = keeper.AllInvariants

	NewQueryPurchaseOrdersParams = types.NewQueryPurchaseOrdersParams

	// Msg functions
	NewMsgUndPurchaseOrder        = types.NewMsgUndPurchaseOrder
	NewMsgProcessUndPurchaseOrder = types.NewMsgProcessUndPurchaseOrder

	// Errors
	ErrInvalidDecision        = types.ErrInvalidDecision
	ErrInvalidDenomination    = types.ErrInvalidDenomination
	ErrNotAuthorisedToRaisePO = types.ErrNotAuthorisedToRaisePO

	NewCheckLockedUndDecorator = ante.NewCheckLockedUndDecorator

	DefaultParams = types.DefaultParams
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params

	QueryPurchaseOrdersParams = types.QueryPurchaseOrdersParams

	// Msgs
	MsgPurchaseUnd             = types.MsgUndPurchaseOrder
	MsgProcessUndPurchaseOrder = types.MsgProcessUndPurchaseOrder

	// Structs
	EnterpriseUndPurchaseOrder = types.EnterpriseUndPurchaseOrder
	UndSupply                  = types.UndSupply
	PurchaseOrderStatus        = types.PurchaseOrderStatus
	LockedUnd                  = types.LockedUnd
)
