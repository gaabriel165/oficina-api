package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type ServiceOrderRepository interface {
	Create(order *entity.ServiceOrder) error
	Update(order *entity.ServiceOrder) error
	FindByID(id string) (*entity.ServiceOrder, error)
	FindAll() ([]*entity.ServiceOrder, error)
}
