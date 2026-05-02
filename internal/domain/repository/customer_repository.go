package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type CustomerRepository interface {
	Create(customer *entity.Customer) error
	Update(customer *entity.Customer) error
	Delete(id string) error
	FindByID(id string) (*entity.Customer, error)
	FindByDocument(document string) (*entity.Customer, error)
	FindAll() ([]*entity.Customer, error)
}
