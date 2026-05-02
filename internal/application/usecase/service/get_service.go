package service

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

func NewGetServiceUseCase(serviceRepo repository.ServiceRepository) *GetServiceUseCase {
	return &GetServiceUseCase{serviceRepo: serviceRepo}
}

func (u *GetServiceUseCase) Execute(id string) (*entity.Service, error) {
	return u.serviceRepo.FindByID(id)
}
