package vehicle

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetVehicleUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewGetVehicleUseCase(vehicleRepo repository.VehicleRepository) *GetVehicleUseCase {
	return &GetVehicleUseCase{vehicleRepo: vehicleRepo}
}

func (u *GetVehicleUseCase) Execute(id string) (*entity.Vehicle, error) {
	return u.vehicleRepo.FindByID(id)
}
