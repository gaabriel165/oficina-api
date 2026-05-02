package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

type ServiceOrderModel struct {
	ID          string                  `gorm:"type:uuid;primaryKey;column:id"`
	CustomerID  string                  `gorm:"type:uuid;not null;column:customer_id"`
	VehicleID   string                  `gorm:"type:uuid;not null;column:vehicle_id"`
	Status      string                  `gorm:"type:varchar(50);not null;column:status"`
	Notes       string                  `gorm:"type:text;column:notes"`
	TotalAmount float64                 `gorm:"type:decimal(10,2);not null;default:0;column:total_amount"`
	StartedAt   *time.Time              `gorm:"column:started_at"`
	FinishedAt  *time.Time              `gorm:"column:finished_at"`
	CreatedAt   time.Time               `gorm:"column:created_at"`
	UpdatedAt   time.Time               `gorm:"column:updated_at"`
	Items       []ServiceOrderItemModel `gorm:"foreignKey:ServiceOrderID"`
	Parts       []ServiceOrderPartModel `gorm:"foreignKey:ServiceOrderID"`
}

func (ServiceOrderModel) TableName() string { return "service_orders" }

func (m *ServiceOrderModel) ToDomain() *entity.ServiceOrder {
	items := make([]*entity.ServiceOrderItem, len(m.Items))
	for i, item := range m.Items {
		items[i] = item.ToDomain()
	}

	parts := make([]*entity.ServiceOrderPart, len(m.Parts))
	for i, part := range m.Parts {
		parts[i] = part.ToDomain()
	}

	return entity.RestoreServiceOrder(
		m.ID, m.CustomerID, m.VehicleID,
		valueobject.OrderStatus(m.Status),
		m.Notes, m.TotalAmount,
		items, parts,
		m.StartedAt, m.FinishedAt,
		m.CreatedAt, m.UpdatedAt,
	)
}

func ServiceOrderModelFromDomain(o *entity.ServiceOrder) *ServiceOrderModel {
	items := make([]ServiceOrderItemModel, len(o.Items()))
	for i, item := range o.Items() {
		items[i] = *ServiceOrderItemModelFromDomain(item)
	}

	parts := make([]ServiceOrderPartModel, len(o.Parts()))
	for i, part := range o.Parts() {
		parts[i] = *ServiceOrderPartModelFromDomain(part)
	}

	return &ServiceOrderModel{
		ID:          o.ID(),
		CustomerID:  o.CustomerID(),
		VehicleID:   o.VehicleID(),
		Status:      o.Status().String(),
		Notes:       o.Notes(),
		TotalAmount: o.TotalAmount(),
		StartedAt:   o.StartedAt(),
		FinishedAt:  o.FinishedAt(),
		CreatedAt:   o.CreatedAt(),
		UpdatedAt:   o.UpdatedAt(),
		Items:       items,
		Parts:       parts,
	}
}
