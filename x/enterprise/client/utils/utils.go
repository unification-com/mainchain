package utils

// NormalisePurchaseOrderStatus - normalise user specified purchase order status
func NormalisePurchaseOrderStatus(status string) string {
	switch status {
	case "Accept", "accept", "Accepted", "accepted":
		return "accept"
	case "Reject", "reject", "Rejected", "rejected":
		return "reject"
	case "Raised", "raised":
		return "raised"
	case "Complete", "complete", "Completed", "completed":
		return "complete"
	}
	return ""
}
