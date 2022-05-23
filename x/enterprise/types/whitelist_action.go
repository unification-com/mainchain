package types

import "fmt"

// PurchaseOrderStatusFromString turns a string into a ProposalStatus
func WhitelistActionFromString(str string) (WhitelistAction, error) {
	switch str {
	case "add":
		return WhitelistActionAdd, nil

	case "remove":
		return WhitelistActionRemove, nil

	case "":
		return WhitelistActionNil, nil

	default:
		return WhitelistActionNil, fmt.Errorf("'%s' is not a valid whitelist action", str)
	}
}

// ValidWhitelistAction returns true if the purchase order status is valid and false
// otherwise.
func ValidWhitelistAction(action WhitelistAction) bool {
	if action == WhitelistActionAdd ||
		action == WhitelistActionRemove {
		return true
	}
	return false
}
