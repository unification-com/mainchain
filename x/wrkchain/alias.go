package wrkchain

import (
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/keeper"
	"github.com/unification-com/mainchain-cosmos/x/wrkchain/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	QueryWrkChain      = keeper.QueryWrkChain
	QueryWrkChainBlock = keeper.QueryWrkChainBlock
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
	FeesWrkChainRegistrationCoin = types.FeesWrkChainRegistrationCoin
	FeesWrkChainRecordHashCoin   = types.FeesWrkChainRecordHashCoin

	// WRKChain fees in sdk.Coins[]
	FeesWrkChainRegistration = types.FeesWrkChainRegistration
	FeesWrkChainRecordHash   = types.FeesChainRecordHash
)

type (
	Keeper = keeper.Keeper

	// Msgs
	MsgRegisterWrkChain    = types.MsgRegisterWrkChain
	MsgRecordWrkChainBlock = types.MsgRecordWrkChainBlock

	// Structs
	WrkChain = types.WrkChain
)
