package types

import "fmt"

// PurchaseOrderStatusFromString turns a string into a ProposalStatus
func PurchaseOrderStatusFromString(str string) (PurchaseOrderStatus, error) {
	switch str {
	case "accept":
		return StatusAccepted, nil

	case "reject":
		return StatusRejected, nil

	case "raised":
		return StatusRaised, nil

	case "complete":
		return StatusCompleted, nil

	case "":
		return StatusNil, nil

	default:
		return StatusNil, fmt.Errorf("'%s' is not a valid purchase order status", str)
	}
}

// ValidPurchaseOrderStatus returns true if the purchase order status is valid and false
// otherwise.
func ValidPurchaseOrderStatus(status PurchaseOrderStatus) bool {
	if status == StatusRaised ||
		status == StatusAccepted ||
		status == StatusRejected ||
		status == StatusCompleted {
		return true
	}
	return false
}

// ValidPurchaseOrderAcceptRejectStatus checks the decision - returns true if accept/reject.
func ValidPurchaseOrderAcceptRejectStatus(status PurchaseOrderStatus) bool {
	if status == StatusAccepted || status == StatusRejected {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status PurchaseOrderStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *PurchaseOrderStatus) Unmarshal(data []byte) error {
	*status = PurchaseOrderStatus(data[0])
	return nil
}

// StringNice outputs the status for logs/queries etc. in a nicer readable way.
func (status PurchaseOrderStatus) StringNice() string {
	switch status {
	case StatusAccepted:
		return "accept"

	case StatusRejected:
		return "reject"

	case StatusRaised:
		return "raised"

	case StatusCompleted:
		return "complete"

	default:
		return ""
	}
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (status PurchaseOrderStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}
