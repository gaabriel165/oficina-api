package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type ServiceModel struct {
	ID               string    `gorm:"type:uuid;primaryKey;column:id"`
	Name             string    `gorm:"type:varchar(255);not null;column:name"`
	Description      string    `gorm:"type:text;column:description"`
	LaborPrice       float64   `gorm:"type:decimal(10,2);not null;column:labor_price"`
	EstimatedMinutes int       `gorm:"not null;column:estimated_minutes"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (ServiceModel) TableName() string { return "services" }

func (m *ServiceModel) ToDomain() *entity.Service {
	return entity.RestoreService(m.ID, m.Name, m.Description, m.LaborPrice, m.EstimatedMinutes, m.CreatedAt, m.UpdatedAt)
}

func ServiceModelFromDomain(s *entity.Service) *ServiceModel {
	return &ServiceModel{
		ID:               s.ID(),
		Name:             s.Name(),
		Description:      s.Description(),
		LaborPrice:       s.LaborPrice(),
		EstimatedMinutes: s.EstimatedMinutes(),
		CreatedAt:        s.CreatedAt(),
		UpdatedAt:        s.UpdatedAt(),
	}
}
