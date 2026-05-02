package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type CreateServiceOrderInput struct {
	CustomerID string
	VehicleID  string
	Notes      string
}

type CreateServiceOrderUseCase struct {
	orderRepo    repository.ServiceOrderRepository
	customerRepo repository.CustomerRepository
	vehicleRepo  repository.VehicleRepository
}

func NewCreateServiceOrderUseCase(
	orderRepo repository.ServiceOrderRepository,
	customerRepo repository.CustomerRepository,
	vehicleRepo repository.VehicleRepository,
) *CreateServiceOrderUseCase {
	return &CreateServiceOrderUseCase{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
	}
}

func (u *CreateServiceOrderUseCase) Execute(input CreateServiceOrderInput) (*entity.ServiceOrder, error) {
	if _, err := u.customerRepo.FindByID(input.CustomerID); err != nil {
		return nil, err
	}

	vehicle, err := u.vehicleRepo.FindByID(input.VehicleID)
	if err != nil {
		return nil, err
	}

	if vehicle.CustomerID() != input.CustomerID {
		return nil, entity.ErrServiceOrderVehicleNotOwnedByCustomer
	}

	order, err := entity.NewServiceOrder(input.CustomerID, input.VehicleID, input.Notes)
	if err != nil {
		return nil, err
	}

	if err := u.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}
