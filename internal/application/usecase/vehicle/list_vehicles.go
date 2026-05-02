package vehicle

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ListVehiclesUseCase struct {
	vehicleRepo repository.VehicleRepository
}

func NewListVehiclesUseCase(vehicleRepo repository.VehicleRepository) *ListVehiclesUseCase {
	return &ListVehiclesUseCase{vehicleRepo: vehicleRepo}
}

func (u *ListVehiclesUseCase) Execute() ([]*entity.Vehicle, error) {
	return u.vehicleRepo.FindAll()
}
