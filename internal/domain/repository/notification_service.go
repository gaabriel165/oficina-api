package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type NotificationService interface {
	NotifyOrderStatusChanged(customer *entity.Customer, order *entity.ServiceOrder) error
}
