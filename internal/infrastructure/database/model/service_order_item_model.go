package model

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type ServiceOrderItemModel struct {
	ID             string  `gorm:"type:uuid;primaryKey;column:id"`
	ServiceOrderID string  `gorm:"type:uuid;not null;column:service_order_id"`
	ServiceID      string  `gorm:"type:uuid;not null;column:service_id"`
	ServiceName    string  `gorm:"type:varchar(255);not null;column:service_name"`
	LaborPrice     float64 `gorm:"type:decimal(10,2);not null;column:labor_price"`
}

func (ServiceOrderItemModel) TableName() string { return "service_order_items" }

func (m *ServiceOrderItemModel) ToDomain() *entity.ServiceOrderItem {
	return entity.RestoreServiceOrderItem(m.ID, m.ServiceOrderID, m.ServiceID, m.ServiceName, m.LaborPrice)
}

func ServiceOrderItemModelFromDomain(i *entity.ServiceOrderItem) *ServiceOrderItemModel {
	return &ServiceOrderItemModel{
		ID:             i.ID(),
		ServiceOrderID: i.ServiceOrderID(),
		ServiceID:      i.ServiceID(),
		ServiceName:    i.ServiceName(),
		LaborPrice:     i.LaborPrice(),
	}
}
