package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetServiceOrderUseCase struct {
	orderRepo repository.ServiceOrderRepository
}

func NewGetServiceOrderUseCase(orderRepo repository.ServiceOrderRepository) *GetServiceOrderUseCase {
	return &GetServiceOrderUseCase{orderRepo: orderRepo}
}

func (u *GetServiceOrderUseCase) Execute(orderID string) (*entity.ServiceOrder, error) {
	return u.orderRepo.FindByID(orderID)
}

func (u *GetServiceOrderUseCase) ExecuteForCustomer(orderID, customerID string) (*entity.ServiceOrder, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.CustomerID() != customerID {
		return nil, usecase.ErrForbidden
	}
	return order, nil
}
