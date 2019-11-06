package types

var (
	EventTypeRaisePurchaseOrder   = "raise-purchase-order"
	EventTypeProcessPurchaseOrder = "process-purchase-order"
	EventTypeUndPurchaseComplete  = "und-purchase-complete"

	AttributeValueCategory = ModuleName

	AttributeKeyPurchaseOrderID = "id"
	AttributeKeyPurchaser       = "purchaser"
	AttributeKeyAmount          = "amount"
	AttributeKeyDecision        = "decision"
)
