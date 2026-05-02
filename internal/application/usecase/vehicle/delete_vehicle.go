package vehicle

import "github.com/gabrielcamargo/oficina-api/internal/domain/repository"

type DeleteVehicleUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewDeleteVehicleUseCase(vehicleRepo repository.VehicleRepository) *DeleteVehicleUseCase {
	return &DeleteVehicleUseCase{vehicleRepo: vehicleRepo}
}

func (u *DeleteVehicleUseCase) Execute(id string) error {
	if _, err := u.vehicleRepo.FindByID(id); err != nil {
		return err
	}
	return u.vehicleRepo.Delete(id)
}
