package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type ServiceRepository interface {
	Create(service *entity.Service) error
	Update(service *entity.Service) error
	Delete(id string) error
	FindByID(id string) (*entity.Service, error)
	FindAll() ([]*entity.Service, error)
}
