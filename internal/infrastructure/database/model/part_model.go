package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type PartModel struct {
	ID            string    `gorm:"type:uuid;primaryKey;column:id"`
	Name          string    `gorm:"type:varchar(255);not null;column:name"`
	Description   string    `gorm:"type:text;column:description"`
	UnitPrice     float64   `gorm:"type:decimal(10,2);not null;column:unit_price"`
	StockQuantity int       `gorm:"not null;default:0;column:stock_quantity"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (PartModel) TableName() string { return "parts" }

func (m *PartModel) ToDomain() *entity.Part {
	return entity.RestorePart(m.ID, m.Name, m.Description, m.UnitPrice, m.StockQuantity, m.CreatedAt, m.UpdatedAt)
}

func PartModelFromDomain(p *entity.Part) *PartModel {
	return &PartModel{
		ID:            p.ID(),
		Name:          p.Name(),
		Description:   p.Description(),
		UnitPrice:     p.UnitPrice(),
		StockQuantity: p.StockQuantity(),
		CreatedAt:     p.CreatedAt(),
		UpdatedAt:     p.UpdatedAt(),
	}
}
