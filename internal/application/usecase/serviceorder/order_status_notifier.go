package serviceorder

import (
	"log/slog"

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
		slog.Warn("notification.skipped",
			slog.String("event", "notification.skipped"),
			slog.String("order_id", order.ID()),
			slog.String("error", err.Error()),
		)
		return
	}

	if err := n.notifications.NotifyOrderStatusChanged(customer, order); err != nil {
		slog.Error("integration.error",
			slog.String("event", "integration.error"),
			slog.String("integration", "resend"),
			slog.String("operation", "notify_order_status_changed"),
			slog.String("order_id", order.ID()),
			slog.String("status", order.Status().String()),
			slog.String("error", err.Error()),
		)
		return
	}

	slog.Info("notification.sent",
		slog.String("event", "notification.sent"),
		slog.String("integration", "resend"),
		slog.String("order_id", order.ID()),
		slog.String("status", order.Status().String()),
	)
}
