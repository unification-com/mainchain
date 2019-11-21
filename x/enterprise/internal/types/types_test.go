package types

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPurchaseOrderStatusFormat(t *testing.T) {
	statusRaised, _ := PurchaseOrderStatusFromString("raised")
	statusAccept, _ := PurchaseOrderStatusFromString("accept")
	statusReject, _ := PurchaseOrderStatusFromString("reject")
	statusComplete, _ := PurchaseOrderStatusFromString("complete")
	tests := []struct {
		pt                   PurchaseOrderStatus
		sprintFArgs          string
		expectedStringOutput string
	}{
		{statusRaised, "%s", "raised"},
		{statusRaised, "%v", "1"},
		{statusAccept, "%s", "accept"},
		{statusAccept, "%v", "2"},
		{statusReject, "%s", "reject"},
		{statusReject, "%v", "3"},
		{statusComplete, "%s", "complete"},
		{statusComplete, "%v", "4"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.sprintFArgs, tt.pt)
		require.Equal(t, tt.expectedStringOutput, got)
	}
}
