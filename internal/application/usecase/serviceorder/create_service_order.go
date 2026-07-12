package serviceorder

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type CreateServiceOrderPartInput struct {
	PartID   string
	Quantity int
}

type CreateServiceOrderInput struct {
	CustomerID string
	VehicleID  string
	Notes      string
	Services   []string
	Parts      []CreateServiceOrderPartInput
}

type CreateServiceOrderUseCase struct {
	orderRepo    repository.ServiceOrderRepository
	customerRepo repository.CustomerRepository
	vehicleRepo  repository.VehicleRepository
	serviceRepo  repository.ServiceRepository
	partRepo     repository.PartRepository
}

func NewCreateServiceOrderUseCase(
	orderRepo repository.ServiceOrderRepository,
	customerRepo repository.CustomerRepository,
	vehicleRepo repository.VehicleRepository,
	serviceRepo repository.ServiceRepository,
	partRepo repository.PartRepository,
) *CreateServiceOrderUseCase {
	return &CreateServiceOrderUseCase{
		orderRepo:    orderRepo,
		customerRepo: customerRepo,
		vehicleRepo:  vehicleRepo,
		serviceRepo:  serviceRepo,
		partRepo:     partRepo,
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

	for _, serviceID := range input.Services {
		service, err := u.serviceRepo.FindByID(serviceID)
		if err != nil {
			return nil, err
		}
		order.AddItem(service)
	}

	for _, partInput := range input.Parts {
		part, err := u.partRepo.FindByID(partInput.PartID)
		if err != nil {
			return nil, err
		}

		if !part.HasSufficientStock(partInput.Quantity) {
			return nil, entity.ErrPartInsufficientStock
		}

		if err := order.AddPart(part, partInput.Quantity); err != nil {
			return nil, err
		}
	}

	if err := u.orderRepo.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}
