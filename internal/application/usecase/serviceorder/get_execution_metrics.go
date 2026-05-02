package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetExecutionMetricsUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewGetExecutionMetricsUseCase(orderRepo repository.ServiceOrderRepository) *GetExecutionMetricsUseCase {
	return &GetExecutionMetricsUseCase{orderRepo: orderRepo}
}

func (u *GetExecutionMetricsUseCase) Execute() (*repository.ExecutionMetrics, error) {
	return u.orderRepo.GetExecutionMetrics()
}
