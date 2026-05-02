package vehicle

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type UpdateVehicleInput struct {
	ID    string
	Brand string
	Model string
	Year  int
}

type UpdateVehicleUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewUpdateVehicleUseCase(vehicleRepo repository.VehicleRepository) *UpdateVehicleUseCase {
	return &UpdateVehicleUseCase{vehicleRepo: vehicleRepo}
}

func (u *UpdateVehicleUseCase) Execute(input UpdateVehicleInput) (*entity.Vehicle, error) {
	vehicle, err := u.vehicleRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := vehicle.Update(input.Brand, input.Model, input.Year); err != nil {
		return nil, err
	}

	if err := u.vehicleRepo.Update(vehicle); err != nil {
		return nil, err
	}

	return vehicle, nil
}
