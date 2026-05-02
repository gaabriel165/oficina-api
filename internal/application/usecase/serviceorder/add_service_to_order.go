package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type AddServiceToOrderInput struct {
	OrderID   string
	ServiceID string
}

type AddServiceToOrderUseCase struct {
	orderRepo   repository.ServiceOrderRepository
	serviceRepo repository.ServiceRepository
}

func NewAddServiceToOrderUseCase(
	orderRepo repository.ServiceOrderRepository,
	serviceRepo repository.ServiceRepository,
) *AddServiceToOrderUseCase {
	return &AddServiceToOrderUseCase{orderRepo: orderRepo, serviceRepo: serviceRepo}
}

func (u *AddServiceToOrderUseCase) Execute(input AddServiceToOrderInput) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(input.OrderID)
	if err != nil {
		return nil, err
	}

	svc, err := u.serviceRepo.FindByID(input.ServiceID)
	if err != nil {
		return nil, err
	}

	order.AddItem(svc)

	if err := u.orderRepo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}
