package model

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type ServiceOrderPartModel struct {
	ID             string  `gorm:"type:uuid;primaryKey;column:id"`
	ServiceOrderID string  `gorm:"type:uuid;not null;column:service_order_id"`
	PartID         string  `gorm:"type:uuid;not null;column:part_id"`
	PartName       string  `gorm:"type:varchar(255);not null;column:part_name"`
	Quantity       int     `gorm:"not null;column:quantity"`
	UnitPrice      float64 `gorm:"type:decimal(10,2);not null;column:unit_price"`
}

func (ServiceOrderPartModel) TableName() string { return "service_order_parts" }

func (m *ServiceOrderPartModel) ToDomain() *entity.ServiceOrderPart {
	return entity.RestoreServiceOrderPart(m.ID, m.ServiceOrderID, m.PartID, m.PartName, m.Quantity, m.UnitPrice)
}

func ServiceOrderPartModelFromDomain(p *entity.ServiceOrderPart) *ServiceOrderPartModel {
	return &ServiceOrderPartModel{
		ID:             p.ID(),
		ServiceOrderID: p.ServiceOrderID(),
		PartID:         p.PartID(),
		PartName:       p.PartName(),
		Quantity:       p.Quantity(),
		UnitPrice:      p.UnitPrice(),
	}
}
