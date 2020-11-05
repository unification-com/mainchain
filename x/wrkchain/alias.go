package wrkchain

import (
	"github.com/unification-com/mainchain/x/wrkchain/exported"
	"github.com/unification-com/mainchain/x/wrkchain/internal/ante"
	"github.com/unification-com/mainchain/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain/x/wrkchain/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = types.DefaultParamspace

	QueryParameters = keeper.QueryParameters
)

var (
	NewKeeper                 = keeper.NewKeeper
	NewQuerier                = keeper.NewQuerier
	NewMsgRegisterWrkChain    = types.NewMsgRegisterWrkChain
	NewMsgRecordWrkChainBlock = types.NewMsgRecordWrkChainBlock
	NewWrkChain               = types.NewWrkChain
	RegisterCodec             = types.RegisterCodec
	ModuleCdc                 = types.ModuleCdc

	// Event tags
	EventTypeRegisterWrkChain    = types.EventTypeRegisterWrkChain
	EventTypeRecordWrkChainBlock = types.EventTypeRecordWrkChainBlock

	AttributeKeyWrkChainId          = types.AttributeKeyWrkChainId
	AttributeKeyWrkChainMoniker     = types.AttributeKeyWrkChainMoniker
	AttributeKeyWrkChainName        = types.AttributeKeyWrkChainName
	AttributeKeyWrkChainGenesisHash = types.AttributeKeyWrkChainGenesisHash
	AttributeKeyBaseType            = types.AttributeKeyBaseType
	AttributeKeyOwner               = types.AttributeKeyOwner
	AttributeKeyBlockHash           = types.AttributeKeyBlockHash
	AttributeKeyBlockHeight         = types.AttributeKeyBlockHeight
	AttributeKeyParentHash          = types.AttributeKeyParentHash
	AttributeKeyHash1               = types.AttributeKeyHash1
	AttributeKeyHash2               = types.AttributeKeyHash2
	AttributeKeyHash3               = types.AttributeKeyHash3

	GetWrkChainIDBytes = types.GetWrkChainIDBytes

	CheckIsWrkChainTx = exported.CheckIsWrkChainTx

	// Error messages
	ErrWrkChainDoesNotExist         = types.ErrWrkChainDoesNotExist
	ErrWrkChainAlreadyRegistered    = types.ErrWrkChainAlreadyRegistered
	ErrInsufficientWrkChainFee      = types.ErrInsufficientWrkChainFee
	ErrTooMuchWrkChainFee           = types.ErrTooMuchWrkChainFee
	ErrFeePayerNotOwner             = types.ErrFeePayerNotOwner
	ErrNotWrkChainOwner             = types.ErrNotWrkChainOwner
	ErrWrkChainBlockAlreadyRecorded = types.ErrWrkChainBlockAlreadyRecorded
	ErrMissingData                  = types.ErrMissingData
	ErrContentTooLarge              = types.ErrContentTooLarge

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	NewQueryWrkChainParams      = types.NewQueryWrkChainParams
	NewQueryWrkChainBlockParams = types.NewQueryWrkChainBlockParams
	NewParamsRetriever          = keeper.NewParamsRetriever

	DefaultParams = types.DefaultParams

	NewCorrectWrkChainFeeDecorator = ante.NewCorrectWrkChainFeeDecorator
)

type (
	Keeper = keeper.Keeper

	GenesisState = types.GenesisState

	QueryWrkChainParams      = types.QueryWrkChainParams
	QueryWrkChainBlockParams = types.QueryWrkChainBlockParams

	// Msgs
	MsgRegisterWrkChain    = types.MsgRegisterWrkChain
	MsgRecordWrkChainBlock = types.MsgRecordWrkChainBlock

	// Structs
	WrkChain                    = types.WrkChain
	WrkChainExport              = types.WrkChainExport
	WrkChainBlock               = types.WrkChainBlock
	WrkChainBlocksGenesisExport = types.WrkChainBlocksGenesisExport

	WrkChains = types.WrkChains

	// Queries
	QueryResWrkChainBlockHashes = types.QueryResWrkChainBlockHashes
)
