package repository

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

type ServiceMetricRow struct {
	ServiceName     string
	AvgMinutes      float64
	CompletedOrders int
}

type ExecutionMetrics struct {
	OverallAvgMinutes float64
	ByService         []ServiceMetricRow
}

type ServiceOrderRepository interface {
	Create(order *entity.ServiceOrder) error
	Update(order *entity.ServiceOrder) error
	FindByID(id string) (*entity.ServiceOrder, error)
	FindByStatuses(statuses []valueobject.OrderStatus) ([]*entity.ServiceOrder, error)
	GetExecutionMetrics() (*ExecutionMetrics, error)
}
