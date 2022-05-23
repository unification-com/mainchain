package types

var (
	EventTypeRegisterBeacon        = RegisterAction
	EventTypeRecordBeaconTimestamp = RecordAction

	AttributeValueCategory = ModuleName

	AttributeKeyOwner               = "beacon_owner"
	AttributeKeyBeaconId            = "beacon_id"
	AttributeKeyBeaconMoniker       = "beacon_moniker"
	AttributeKeyBeaconName          = "beacon_name"
	AttributeKeyTimestampID         = "beacon_timestamp_id"
	AttributeKeyTimestampHash       = "beacon_timestamp_hash"
	AttributeKeyTimestampSubmitTime = "beacon_timestamp_submit_time"
)
