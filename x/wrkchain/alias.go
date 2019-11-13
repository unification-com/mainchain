package wrkchain

import (
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/exported"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	DefaultCodespace = types.DefaultCodespace

	QueryWrkChain            = keeper.QueryWrkChain
	QueryWrkChainBlock       = keeper.QueryWrkChainBlock
	QueryWrkChainBlockHashes = keeper.QueryWrkChainBlockHashes

	FeeDenom  = types.FeeDenom
	RegFee    = types.RegFee
	RecordFee = types.RecordFee
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

	// WRKChain fees in sdk.Coin(denom=und)
	FeesBaseDenomination         = types.FeesBaseDenomination
	FeesWrkChainRegistrationCoin = types.FeesWrkChainRegistrationCoin
	FeesWrkChainRecordHashCoin   = types.FeesWrkChainRecordHashCoin

	// WRKChain fees in sdk.Coins[]
	FeesWrkChainRegistration = types.FeesWrkChainRegistration
	FeesWrkChainRecordHash   = types.FeesWrkChainRecordHash

	RegisteredWrkChainPrefix        = types.RegisteredWrkChainPrefix
	RecordedWrkChainBlockHashPrefix = types.RecordedWrkChainBlockHashPrefix

	GetWrkChainIDBytes = types.GetWrkChainIDBytes

	CheckIsWrkChainTx = exported.CheckIsWrkChainTx

	// Error messages
	ErrWrkChainDoesNotExist      = types.ErrWrkChainDoesNotExist
	ErrWrkChainAlreadyRegistered = types.ErrWrkChainAlreadyRegistered
	ErrInsufficientWrkChainFee   = types.ErrInsufficientWrkChainFee
	ErrTooMuchWrkChainFee        = types.ErrTooMuchWrkChainFee
	ErrFeePayerNotOwner          = types.ErrFeePayerNotOwner

	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)

type (
	Keeper = keeper.Keeper

	GenesisState = types.GenesisState

	// Msgs
	MsgRegisterWrkChain    = types.MsgRegisterWrkChain
	MsgRecordWrkChainBlock = types.MsgRecordWrkChainBlock

	// Structs
	WrkChain = types.WrkChain

	// Queries
	QueryResWrkChainBlockHashes = types.QueryResWrkChainBlockHashes
)
