package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type VehicleRepository interface {
	Create(vehicle *entity.Vehicle) error
	Update(vehicle *entity.Vehicle) error
	Delete(id string) error
	FindByID(id string) (*entity.Vehicle, error)
	FindByPlate(plate string) (*entity.Vehicle, error)
	FindByCustomerID(customerID string) ([]*entity.Vehicle, error)
	FindAll() ([]*entity.Vehicle, error)
}
