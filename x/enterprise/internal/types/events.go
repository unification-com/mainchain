package types

var (
	EventTypeRaisePurchaseOrder   = "raise_purchase_order"
	EventTypeProcessPurchaseOrder = "process_purchase_order"
	EventTypeUndPurchaseComplete  = "und_purchase_complete"
	EventTypeUndUnlocked          = "und_unlocked"

	AttributeValueCategory = ModuleName

	AttributeKeyPurchaseOrderID = "id"
	AttributeKeyPurchaser       = "purchaser"
	AttributeKeyAmount          = "amount"
	AttributeKeyDecision        = "decision"
)
