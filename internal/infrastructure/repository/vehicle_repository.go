package repository

import (
	"errors"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormVehicleRepository struct {
	db *gorm.DB
}

func NewGormVehicleRepository(db *gorm.DB) *GormVehicleRepository {
	return &GormVehicleRepository{db: db}
}

func (r *GormVehicleRepository) Create(vehicle *entity.Vehicle) error {
	return r.db.Create(model.VehicleModelFromDomain(vehicle)).Error
}

func (r *GormVehicleRepository) Update(vehicle *entity.Vehicle) error {
	return r.db.Save(model.VehicleModelFromDomain(vehicle)).Error
}

func (r *GormVehicleRepository) Delete(id string) error {
	return r.db.Delete(&model.VehicleModel{}, "id = ?", id).Error
}

func (r *GormVehicleRepository) FindByID(id string) (*entity.Vehicle, error) {
	var m model.VehicleModel
	err := r.db.First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainrepo.ErrVehicleNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormVehicleRepository) FindByPlate(plate string) (*entity.Vehicle, error) {
	var m model.VehicleModel
	err := r.db.First(&m, "plate = ?", plate).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainrepo.ErrVehicleNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormVehicleRepository) FindByCustomerID(customerID string) ([]*entity.Vehicle, error) {
	var models []model.VehicleModel
	if err := r.db.Find(&models, "customer_id = ?", customerID).Error; err != nil {
		return nil, err
	}
	vehicles := make([]*entity.Vehicle, len(models))
	for i, m := range models {
		vehicles[i] = m.ToDomain()
	}
	return vehicles, nil
}

func (r *GormVehicleRepository) FindAll() ([]*entity.Vehicle, error) {
	var models []model.VehicleModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	vehicles := make([]*entity.Vehicle, len(models))
	for i, m := range models {
		vehicles[i] = m.ToDomain()
	}
	return vehicles, nil
}
