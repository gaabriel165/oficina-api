package service

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type UpdateServiceInput struct {
	ID               string
	Name             string
	Description      string
	LaborPrice       float64
	EstimatedMinutes int
}

type UpdateServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

func NewUpdateServiceUseCase(serviceRepo repository.ServiceRepository) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{serviceRepo: serviceRepo}
}

func (u *UpdateServiceUseCase) Execute(input UpdateServiceInput) (*entity.Service, error) {
	svc, err := u.serviceRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := svc.Update(input.Name, input.Description, input.LaborPrice, input.EstimatedMinutes); err != nil {
		return nil, err
	}

	if err := u.serviceRepo.Update(svc); err != nil {
		return nil, err
	}

	return svc, nil
}
