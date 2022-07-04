package types

var (
	EventTypeRegisterWrkChain    = RegisterAction
	EventTypeRecordWrkChainBlock = RecordAction
	EventTypePurchaseStorage     = PurchaseStorageAction

	AttributeValueCategory = ModuleName

	AttributeKeyOwner                         = "wrkchain_owner"
	AttributeKeyWrkChainId                    = "wrkchain_id"
	AttributeKeyWrkChainMoniker               = "wrkchain_moniker"
	AttributeKeyWrkChainName                  = "wrkchain_name"
	AttributeKeyWrkChainGenesisHash           = "wrkchain_genesis_hash"
	AttributeKeyBaseType                      = "wrkchain_base_type"
	AttributeKeyBlockHash                     = "wrkchain_block_hash"
	AttributeKeyBlockHeight                   = "wrkchain_block_height"
	AttributeKeyPrunedBlockHeight             = "wrkchain_pruned_block_height"
	AttributeKeyParentHash                    = "wrkchain_parent_hash"
	AttributeKeyHash1                         = "wrkchain_hash1"
	AttributeKeyHash2                         = "wrkchain_hash2"
	AttributeKeyHash3                         = "wrkchain_hash3"
	AttributeKeyWrkChainStorageNumPurchased   = "wrkchain_storage_num_purchased"
	AttributeKeyWrkChainStorageNumCanPurchase = "wrkchain_storage_num_can_purchase"
)
