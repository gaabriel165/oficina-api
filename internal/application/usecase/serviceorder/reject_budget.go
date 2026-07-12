package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type RejectBudgetUseCase struct {
	orderRepo repository.ServiceOrderRepository
	notifier  StatusNotifier
}

func NewRejectBudgetUseCase(orderRepo repository.ServiceOrderRepository, notifier StatusNotifier) *RejectBudgetUseCase {
	return &RejectBudgetUseCase{orderRepo: orderRepo, notifier: notifier}
}

func (u *RejectBudgetUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.RejectBudget(); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	u.notifier.Notify(order)

	return order, nil
}
