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

//// StringNice outputs the action for logs/queries etc. in a nicer readable way.
//func (action WhitelistAction) StringNice() string {
//	switch action {
//	case WhitelistActionAdd:
//		return "add"
//
//	case WhitelistActionRemove:
//		return "remove"
//
//	default:
//		return ""
//	}
//}
//
//// Format implements the fmt.Formatter interface.
//// nolint: errcheck
//func (action WhitelistAction) Format(s fmt.State, verb rune) {
//	switch verb {
//	case 's':
//		s.Write([]byte(action.String()))
//	default:
//		// TODO: Do this conversion more directly
//		s.Write([]byte(fmt.Sprintf("%v", byte(action))))
//	}
//}
