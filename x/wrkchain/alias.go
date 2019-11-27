package wrkchain

import (
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/exported"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/ante"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultParamspace = types.DefaultParamspace

	DefaultCodespace = types.DefaultCodespace

	QueryParameters          = keeper.QueryParameters
)

var (
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	NewMsgRegisterWrkChain = types.NewMsgRegisterWrkChain
	NewWrkChain            = types.NewWrkChain
	RegisterCodec          = types.RegisterCodec
	ModuleCdc              = types.ModuleCdc

	// Event tags
	EventTypeRegisterWrkChain    = types.EventTypeRegisterWrkChain
	EventTypeRecordWrkChainBlock = types.EventTypeRecordWrkChainBlock

	AttributeKeyWrkChainId          = types.AttributeKeyWrkChainId
	AttributeKeyWrkChainMoniker     = types.AttributeKeyWrkChainMoniker
	AttributeKeyWrkChainName        = types.AttributeKeyWrkChainName
	AttributeKeyWrkChainGenesisHash = types.AttributeKeyWrkChainGenesisHash
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

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	NewQueryWrkChainParams      = types.NewQueryWrkChainParams
	NewQueryWrkChainBlockParams = types.NewQueryWrkChainBlockParams

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
	WrkChain       = types.WrkChain
	WrkChainExport = types.WrkChainExport
	WrkChainBlock  = types.WrkChainBlock

	WrkChains = types.WrkChains

	// Queries
	QueryResWrkChainBlockHashes = types.QueryResWrkChainBlockHashes
)
