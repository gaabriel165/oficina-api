package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type SendBudgetUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewSendBudgetUseCase(orderRepo repository.ServiceOrderRepository) *SendBudgetUseCase {
	return &SendBudgetUseCase{orderRepo: orderRepo}
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

	return order, nil
}
