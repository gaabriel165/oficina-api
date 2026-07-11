package valueobject

import "errors"

var ErrInvalidStatusTransition = errors.New("invalid order status transition")

type OrderStatus string

const (
	OrderStatusReceived        OrderStatus = "received"
	OrderStatusInDiagnosis     OrderStatus = "in_diagnosis"
	OrderStatusWaitingApproval OrderStatus = "waiting_approval"
	OrderStatusInExecution     OrderStatus = "in_execution"
	OrderStatusFinished        OrderStatus = "finished"
	OrderStatusDelivered       OrderStatus = "delivered"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

var allowedTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusReceived:        {OrderStatusInDiagnosis},
	OrderStatusInDiagnosis:     {OrderStatusWaitingApproval},
	OrderStatusWaitingApproval: {OrderStatusInExecution, OrderStatusCancelled},
	OrderStatusInExecution:     {OrderStatusFinished},
	OrderStatusFinished:        {OrderStatusDelivered},
	OrderStatusDelivered:       {},
	OrderStatusCancelled:       {},
}

var listingPriority = map[OrderStatus]int{
	OrderStatusInExecution:     1,
	OrderStatusWaitingApproval: 2,
	OrderStatusInDiagnosis:     3,
	OrderStatusReceived:        4,
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	allowed := allowedTransitions[s]
	for _, status := range allowed {
		if status == next {
			return true
		}
	}
	return false
}

func (s OrderStatus) ListingPriority() int {
	return listingPriority[s]
}

func (s OrderStatus) String() string {
	return string(s)
}
