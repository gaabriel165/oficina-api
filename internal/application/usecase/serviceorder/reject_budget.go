package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type RejectBudgetUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewRejectBudgetUseCase(orderRepo repository.ServiceOrderRepository) *RejectBudgetUseCase {
	return &RejectBudgetUseCase{orderRepo: orderRepo}
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

	return order, nil
}
