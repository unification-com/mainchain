package types

var (
	EventTypeRaisePurchaseOrder           = "raise_purchase_order"
	EventTypeProcessPurchaseOrderDecision = "process_purchase_order_decision"
	EventTypeAutoRejectStalePurchaseOrder = "auto_reject_stale_purchase_order"
	EventTypeTallyPurchaseOrderDecisions  = "tally_purchase_order_decisions"
	EventTypeUndPurchaseComplete          = "und_purchase_complete"
	EventTypeUndUnlocked                  = "und_unlocked"
	EventTypeWhitelistAddress             = "whitelist_purchase_order_address"

	AttributeValueCategory = ModuleName

	AttributeKeyPurchaseOrderID = "id"
	AttributeKeyPurchaser       = "purchaser"
	AttributeKeyAmount          = "amount"
	AttributeKeyDecision        = "decision"
	AttributeKeySigner          = "signer"
	AttributeKeyNumAccepts      = "accepts"
	AttributeKeyNumRejects      = "rejects"
	AttributeKeyWhitelistAction = "action"
	AttributeWhitelistAddress   = "address"
)

