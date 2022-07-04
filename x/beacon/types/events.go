package types

var (
	EventTypeRegisterBeacon        = RegisterAction
	EventTypeRecordBeaconTimestamp = RecordAction
	EventTypePurchaseStorage       = PurchaseStorageAction

	AttributeValueCategory = ModuleName

	AttributeKeyOwner                       = "beacon_owner"
	AttributeKeyBeaconId                    = "beacon_id"
	AttributeKeyBeaconMoniker               = "beacon_moniker"
	AttributeKeyBeaconName                  = "beacon_name"
	AttributeKeyTimestampID                 = "beacon_timestamp_id"
	AttributeKeyTimestampHash               = "beacon_timestamp_hash"
	AttributeKeyTimestampIdPruned           = "beacon_timestamp_id_pruned"
	AttributeKeyTimestampSubmitTime         = "beacon_timestamp_submit_time"
	AttributeKeyBeaconStorageNumPurchased   = "beacon_storage_num_purchased"
	AttributeKeyBeaconStorageNumCanPurchase = "beacon_storage_num_can_purchase"
)
