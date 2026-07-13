package serviceorder

import (
	"log"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type StatusNotifier interface {
	Notify(order *entity.ServiceOrder)
}

type OrderStatusNotifier struct {
	customerRepo  repository.CustomerRepository
	notifications repository.NotificationService
}

func NewOrderStatusNotifier(
	customerRepo repository.CustomerRepository,
	notifications repository.NotificationService,
) *OrderStatusNotifier {
	return &OrderStatusNotifier{customerRepo: customerRepo, notifications: notifications}
}

func (n *OrderStatusNotifier) Notify(order *entity.ServiceOrder) {
	customer, err := n.customerRepo.FindByID(order.CustomerID())
	if err != nil {
		log.Printf("order %s notification skipped: %v", order.ID(), err)
		return
	}

	if err := n.notifications.NotifyOrderStatusChanged(customer, order); err != nil {
		log.Printf("order %s notification failed: %v", order.ID(), err)
	}
}
