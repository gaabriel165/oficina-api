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

func (r *GormServiceOrderRepository) GetExecutionMetrics() (*domainrepo.ExecutionMetrics, error) {
	var overallAvg float64
	r.db.Raw(`
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (finished_at - started_at)) / 60), 0)
		FROM service_orders
		WHERE finished_at IS NOT NULL AND started_at IS NOT NULL AND deleted_at IS NULL
	`).Scan(&overallAvg)

	type row struct {
		ServiceName     string  `gorm:"column:service_name"`
		AvgMinutes      float64 `gorm:"column:avg_minutes"`
		CompletedOrders int     `gorm:"column:completed_orders"`
	}

	var rows []row
	if err := r.db.Raw(`
		SELECT
			s.name AS service_name,
			COALESCE(AVG(EXTRACT(EPOCH FROM (so.finished_at - so.started_at)) / 60), 0) AS avg_minutes,
			COUNT(*) AS completed_orders
		FROM service_orders so
		JOIN service_order_items soi ON soi.service_order_id = so.id
		JOIN services s ON s.id = soi.service_id
		WHERE so.finished_at IS NOT NULL AND so.started_at IS NOT NULL
		GROUP BY s.id, s.name
		ORDER BY s.name
	`).Scan(&rows).Error; err != nil {
		return nil, err
	}

	metrics := &domainrepo.ExecutionMetrics{
		OverallAvgMinutes: overallAvg,
		ByService:         make([]domainrepo.ServiceMetricRow, 0, len(rows)),
	}
	for _, r := range rows {
		metrics.ByService = append(metrics.ByService, domainrepo.ServiceMetricRow{
			ServiceName:     r.ServiceName,
			AvgMinutes:      r.AvgMinutes,
			CompletedOrders: r.CompletedOrders,
		})
	}
	return metrics, nil
}
