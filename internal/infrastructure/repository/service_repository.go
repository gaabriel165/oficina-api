package repository

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormServiceRepository struct {
	db *gorm.DB
}

func NewGormServiceRepository(db *gorm.DB) *GormServiceRepository {
	return &GormServiceRepository{db: db}
}

func (r *GormServiceRepository) Create(service *entity.Service) error {
	return r.db.Create(model.ServiceModelFromDomain(service)).Error
}

func (r *GormServiceRepository) Update(service *entity.Service) error {
	return r.db.Save(model.ServiceModelFromDomain(service)).Error
}

func (r *GormServiceRepository) Delete(id string) error {
	return r.db.Delete(&model.ServiceModel{}, "id = ?", id).Error
}

func (r *GormServiceRepository) FindByID(id string) (*entity.Service, error) {
	var m model.ServiceModel
	err := r.db.First(&m, "id = ?", id).Error
	if isNotFound(err) {
		return nil, domainrepo.ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormServiceRepository) FindAll() ([]*entity.Service, error) {
	var models []model.ServiceModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	services := make([]*entity.Service, len(models))
	for i, m := range models {
		services[i] = m.ToDomain()
	}
	return services, nil
}
