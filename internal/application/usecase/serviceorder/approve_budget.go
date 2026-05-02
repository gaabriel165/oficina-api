package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ApproveBudgetUseCase struct {
	orderRepo repository.ServiceOrderRepository
	partRepo  repository.PartRepository
}

func NewApproveBudgetUseCase(
	orderRepo repository.ServiceOrderRepository,
	partRepo repository.PartRepository,
) *ApproveBudgetUseCase {
	return &ApproveBudgetUseCase{orderRepo: orderRepo, partRepo: partRepo}
}

func (u *ApproveBudgetUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if err := order.ApproveBudget(); err != nil {
		return nil, err
	}

	for _, orderPart := range order.Parts() {
		part, err := u.partRepo.FindByID(orderPart.PartID())
		if err != nil {
			return nil, err
		}

		if err := part.DebitStock(orderPart.Quantity()); err != nil {
			return nil, err
		}

		if err := u.partRepo.Update(part); err != nil {
			return nil, err
		}
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
