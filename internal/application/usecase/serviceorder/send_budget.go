package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type SendBudgetUseCase struct {
	orderRepo repository.ServiceOrderRepository
	notifier  StatusNotifier
}

func NewSendBudgetUseCase(orderRepo repository.ServiceOrderRepository, notifier StatusNotifier) *SendBudgetUseCase {
	return &SendBudgetUseCase{orderRepo: orderRepo, notifier: notifier}
}

func (u *SendBudgetUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.SendBudget(); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	u.notifier.Notify(order)

	return order, nil
}
