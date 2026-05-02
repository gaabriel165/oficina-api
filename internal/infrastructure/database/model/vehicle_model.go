package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type VehicleModel struct {
	ID         string    `gorm:"type:uuid;primaryKey;column:id"`
	CustomerID string    `gorm:"type:uuid;not null;column:customer_id"`
	Plate      string    `gorm:"type:varchar(7);not null;uniqueIndex;column:plate"`
	Brand      string    `gorm:"type:varchar(100);not null;column:brand"`
	Model      string    `gorm:"type:varchar(100);not null;column:model"`
	Year       int       `gorm:"not null;column:year"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (VehicleModel) TableName() string { return "vehicles" }

func (m *VehicleModel) ToDomain() *entity.Vehicle {
	return entity.RestoreVehicle(m.ID, m.CustomerID, m.Plate, m.Brand, m.Model, m.Year, m.CreatedAt, m.UpdatedAt)
}

func VehicleModelFromDomain(v *entity.Vehicle) *VehicleModel {
	return &VehicleModel{
		ID:         v.ID(),
		CustomerID: v.CustomerID(),
		Plate:      v.Plate(),
		Brand:      v.Brand(),
		Model:      v.Model(),
		Year:       v.Year(),
		CreatedAt:  v.CreatedAt(),
		UpdatedAt:  v.UpdatedAt(),
	}
}
