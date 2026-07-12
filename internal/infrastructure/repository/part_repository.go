package repository

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormPartRepository struct {
	db *gorm.DB
}

func NewGormPartRepository(db *gorm.DB) *GormPartRepository {
	return &GormPartRepository{db: db}
}

func (r *GormPartRepository) Create(part *entity.Part) error {
	return r.db.Create(model.PartModelFromDomain(part)).Error
}

func (r *GormPartRepository) Update(part *entity.Part) error {
	return r.db.Save(model.PartModelFromDomain(part)).Error
}

func (r *GormPartRepository) Delete(id string) error {
	return r.db.Delete(&model.PartModel{}, "id = ?", id).Error
}

func (r *GormPartRepository) FindByID(id string) (*entity.Part, error) {
	var m model.PartModel
	err := r.db.First(&m, "id = ?", id).Error
	if isNotFound(err) {
		return nil, domainrepo.ErrPartNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormPartRepository) FindAll() ([]*entity.Part, error) {
	var models []model.PartModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	parts := make([]*entity.Part, len(models))
	for i, m := range models {
		parts[i] = m.ToDomain()
	}
	return parts, nil
}
