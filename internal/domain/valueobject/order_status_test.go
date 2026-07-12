package valueobject_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestOrderStatus_ShouldAllowValidTransitions(t *testing.T) {
	cases := []struct {
		from valueobject.OrderStatus
		to   valueobject.OrderStatus
	}{
		{valueobject.OrderStatusReceived, valueobject.OrderStatusInDiagnosis},
		{valueobject.OrderStatusInDiagnosis, valueobject.OrderStatusWaitingApproval},
		{valueobject.OrderStatusWaitingApproval, valueobject.OrderStatusInExecution},
		{valueobject.OrderStatusWaitingApproval, valueobject.OrderStatusCancelled},
		{valueobject.OrderStatusInExecution, valueobject.OrderStatusFinished},
		{valueobject.OrderStatusFinished, valueobject.OrderStatusDelivered},
	}

	for _, c := range cases {
		assert.True(t, c.from.CanTransitionTo(c.to), "expected %s -> %s to be allowed", c.from, c.to)
	}
}

func TestOrderStatus_ShouldRankListingPriorityByUrgency(t *testing.T) {
	assert.Less(t, valueobject.OrderStatusInExecution.ListingPriority(), valueobject.OrderStatusWaitingApproval.ListingPriority())
	assert.Less(t, valueobject.OrderStatusWaitingApproval.ListingPriority(), valueobject.OrderStatusInDiagnosis.ListingPriority())
	assert.Less(t, valueobject.OrderStatusInDiagnosis.ListingPriority(), valueobject.OrderStatusReceived.ListingPriority())
}

func TestOrderStatus_ShouldDenyInvalidTransitions(t *testing.T) {
	cases := []struct {
		from valueobject.OrderStatus
		to   valueobject.OrderStatus
	}{
		{valueobject.OrderStatusReceived, valueobject.OrderStatusInExecution},
		{valueobject.OrderStatusReceived, valueobject.OrderStatusDelivered},
		{valueobject.OrderStatusInDiagnosis, valueobject.OrderStatusDelivered},
		{valueobject.OrderStatusDelivered, valueobject.OrderStatusReceived},
		{valueobject.OrderStatusCancelled, valueobject.OrderStatusInExecution},
		{valueobject.OrderStatusFinished, valueobject.OrderStatusCancelled},
	}

	for _, c := range cases {
		assert.False(t, c.from.CanTransitionTo(c.to), "expected %s -> %s to be denied", c.from, c.to)
	}
}
