package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type AddPartToOrderInput struct {
	OrderID  string
	PartID   string
	Quantity int
}

type AddPartToOrderUseCase struct {
	orderRepo repository.ServiceOrderRepository
	partRepo  repository.PartRepository
}

func NewAddPartToOrderUseCase(
	orderRepo repository.ServiceOrderRepository,
	partRepo repository.PartRepository,
) *AddPartToOrderUseCase {
	return &AddPartToOrderUseCase{orderRepo: orderRepo, partRepo: partRepo}
}

func (u *AddPartToOrderUseCase) Execute(input AddPartToOrderInput) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(input.OrderID)
	if err != nil {
		return nil, err
	}

	part, err := u.partRepo.FindByID(input.PartID)
	if err != nil {
		return nil, err
	}

	if !part.HasSufficientStock(input.Quantity) {
		return nil, entity.ErrPartInsufficientStock
	}

	if err := order.AddPart(part, input.Quantity); err != nil {
		return nil, err
	}

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
