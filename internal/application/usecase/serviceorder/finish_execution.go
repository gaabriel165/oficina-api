package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type FinishExecutionUseCase struct {
	orderRepo repository.ServiceOrderRepository
	notifier  StatusNotifier
}

func NewFinishExecutionUseCase(orderRepo repository.ServiceOrderRepository, notifier StatusNotifier) *FinishExecutionUseCase {
	return &FinishExecutionUseCase{orderRepo: orderRepo, notifier: notifier}
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

	u.notifier.Notify(order)

	return order, nil
}
