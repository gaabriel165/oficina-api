package repository

import (
	"errors"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormServiceOrderRepository struct {
	db *gorm.DB
}

func NewGormServiceOrderRepository(db *gorm.DB) *GormServiceOrderRepository {
	return &GormServiceOrderRepository{db: db}
}

func (r *GormServiceOrderRepository) Create(order *entity.ServiceOrder) error {
	return r.db.Create(model.ServiceOrderModelFromDomain(order)).Error
}

func (r *GormServiceOrderRepository) Update(order *entity.ServiceOrder) error {
	m := model.ServiceOrderModelFromDomain(order)
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(m).Error; err != nil {
			return err
		}
		if err := tx.Where("service_order_id = ?", m.ID).Delete(&model.ServiceOrderItemModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("service_order_id = ?", m.ID).Delete(&model.ServiceOrderPartModel{}).Error; err != nil {
			return err
		}
		if len(m.Items) > 0 {
			if err := tx.Create(&m.Items).Error; err != nil {
				return err
			}
		}
		if len(m.Parts) > 0 {
			if err := tx.Create(&m.Parts).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormServiceOrderRepository) FindByID(id string) (*entity.ServiceOrder, error) {
	var m model.ServiceOrderModel
	err := r.db.Preload("Items").Preload("Parts").First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainrepo.ErrServiceOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormServiceOrderRepository) FindAll() ([]*entity.ServiceOrder, error) {
	var models []model.ServiceOrderModel
	if err := r.db.Preload("Items").Preload("Parts").Find(&models).Error; err != nil {
		return nil, err
	}
	orders := make([]*entity.ServiceOrder, len(models))
	for i, m := range models {
		orders[i] = m.ToDomain()
	}
	return orders, nil
}
