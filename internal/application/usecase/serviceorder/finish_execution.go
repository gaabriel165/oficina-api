package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type FinishExecutionUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewFinishExecutionUseCase(orderRepo repository.ServiceOrderRepository) *FinishExecutionUseCase {
	return &FinishExecutionUseCase{orderRepo: orderRepo}
}

func (u *FinishExecutionUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.FinishExecution(); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
