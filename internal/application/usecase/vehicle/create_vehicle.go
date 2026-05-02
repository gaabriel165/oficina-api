package vehicle

import (
	"errors"

	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
)

type CreateVehicleInput struct {
	CustomerID string
	Plate      string
	Brand      string
	Model      string
	Year       int
}

type CreateVehicleUseCase struct {
	vehicleRepo  repository.VehicleRepository
	customerRepo repository.CustomerRepository
}

func NewCreateVehicleUseCase(
	vehicleRepo repository.VehicleRepository,
	customerRepo repository.CustomerRepository,
) *CreateVehicleUseCase {
	return &CreateVehicleUseCase{vehicleRepo: vehicleRepo, customerRepo: customerRepo}
}

func (u *CreateVehicleUseCase) Execute(input CreateVehicleInput) (*entity.Vehicle, error) {
	if _, err := u.customerRepo.FindByID(input.CustomerID); err != nil {
		return nil, err
	}

	plateVO, err := valueobject.NewPlate(input.Plate)
	if err != nil {
		return nil, err
	}

	_, err = u.vehicleRepo.FindByPlate(plateVO.Value())
	if err == nil {
		return nil, usecase.ErrPlateAlreadyRegistered
	}
	if !errors.Is(err, repository.ErrVehicleNotFound) {
		return nil, err
	}

	vehicle, err := entity.NewVehicle(input.CustomerID, input.Plate, input.Brand, input.Model, input.Year)
	if err != nil {
		return nil, err
	}

	if err := u.vehicleRepo.Create(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}
