package entity

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/google/uuid"
)

type StatusTransition struct {
	From     valueobject.OrderStatus
	To       valueobject.OrderStatus
	Duration time.Duration
}

type ServiceOrder struct {
	id             string
	customerID     string
	vehicleID      string
	status         valueobject.OrderStatus
	notes          string
	totalAmount    float64
	items          []*ServiceOrderItem
	parts          []*ServiceOrderPart
	startedAt      *time.Time
	finishedAt     *time.Time
	createdAt      time.Time
	updatedAt      time.Time
	lastTransition *StatusTransition
}

func NewServiceOrder(customerID, vehicleID, notes string) (*ServiceOrder, error) {
	if customerID == "" {
		return nil, ErrServiceOrderCustomerRequired
	}
	if vehicleID == "" {
		return nil, ErrServiceOrderVehicleRequired
	}

	now := time.Now()
	return &ServiceOrder{
		id:         uuid.NewString(),
		customerID: customerID,
		vehicleID:  vehicleID,
		status:     valueobject.OrderStatusReceived,
		notes:      notes,
		items:      []*ServiceOrderItem{},
		parts:      []*ServiceOrderPart{},
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func RestoreServiceOrder(
	id, customerID, vehicleID string,
	status valueobject.OrderStatus,
	notes string,
	totalAmount float64,
	items []*ServiceOrderItem,
	parts []*ServiceOrderPart,
	startedAt, finishedAt *time.Time,
	createdAt, updatedAt time.Time,
) *ServiceOrder {
	return &ServiceOrder{
		id:          id,
		customerID:  customerID,
		vehicleID:   vehicleID,
		status:      status,
		notes:       notes,
		totalAmount: totalAmount,
		items:       items,
		parts:       parts,
		startedAt:   startedAt,
		finishedAt:  finishedAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (o *ServiceOrder) ID() string                        { return o.id }
func (o *ServiceOrder) CustomerID() string                { return o.customerID }
func (o *ServiceOrder) VehicleID() string                 { return o.vehicleID }
func (o *ServiceOrder) Status() valueobject.OrderStatus   { return o.status }
func (o *ServiceOrder) Notes() string                     { return o.notes }
func (o *ServiceOrder) TotalAmount() float64              { return o.totalAmount }
func (o *ServiceOrder) Items() []*ServiceOrderItem        { return o.items }
func (o *ServiceOrder) Parts() []*ServiceOrderPart        { return o.parts }
func (o *ServiceOrder) StartedAt() *time.Time             { return o.startedAt }
func (o *ServiceOrder) FinishedAt() *time.Time            { return o.finishedAt }
func (o *ServiceOrder) CreatedAt() time.Time              { return o.createdAt }
func (o *ServiceOrder) UpdatedAt() time.Time              { return o.updatedAt }
func (o *ServiceOrder) LastTransition() *StatusTransition { return o.lastTransition }

func (o *ServiceOrder) StartDiagnosis() error {
	return o.transitionTo(valueobject.OrderStatusInDiagnosis)
}

func (o *ServiceOrder) AddItem(service *Service) {
	item := newServiceOrderItem(o.id, service)
	o.items = append(o.items, item)
	o.updatedAt = time.Now()
}

func (o *ServiceOrder) AddPart(part *Part, quantity int) error {
	if quantity <= 0 {
		return ErrServiceOrderItemInvalidQty
	}
	orderPart := newServiceOrderPart(o.id, part, quantity)
	o.parts = append(o.parts, orderPart)
	o.updatedAt = time.Now()
	return nil
}

func (o *ServiceOrder) SendBudget() error {
	if len(o.items) == 0 && len(o.parts) == 0 {
		return ErrServiceOrderNoItems
	}
	o.calculateTotalAmount()
	return o.transitionTo(valueobject.OrderStatusWaitingApproval)
}

func (o *ServiceOrder) ApproveBudget() error {
	if err := o.transitionTo(valueobject.OrderStatusInExecution); err != nil {
		return err
	}
	now := time.Now()
	o.startedAt = &now
	return nil
}

func (o *ServiceOrder) RejectBudget() error {
	return o.transitionTo(valueobject.OrderStatusCancelled)
}

func (o *ServiceOrder) FinishExecution() error {
	if err := o.transitionTo(valueobject.OrderStatusFinished); err != nil {
		return err
	}
	now := time.Now()
	o.finishedAt = &now
	return nil
}

func (o *ServiceOrder) Deliver() error {
	return o.transitionTo(valueobject.OrderStatusDelivered)
}

func (o *ServiceOrder) transitionTo(next valueobject.OrderStatus) error {
	if !o.status.CanTransitionTo(next) {
		return ErrServiceOrderInvalidTransition
	}
	now := time.Now()
	o.lastTransition = &StatusTransition{From: o.status, To: next, Duration: now.Sub(o.updatedAt)}
	o.status = next
	o.updatedAt = now
	return nil
}

func (o *ServiceOrder) calculateTotalAmount() {
	total := 0.0
	for _, item := range o.items {
		total += item.LaborPrice()
	}
	for _, part := range o.parts {
		total += part.TotalPrice()
	}
	o.totalAmount = total
}
