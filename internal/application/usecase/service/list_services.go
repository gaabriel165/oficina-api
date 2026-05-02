package service

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ListServicesUseCase struct {
	serviceRepo repository.ServiceRepository
}

func NewListServicesUseCase(serviceRepo repository.ServiceRepository) *ListServicesUseCase {
	return &ListServicesUseCase{serviceRepo: serviceRepo}
}

func (u *ListServicesUseCase) Execute() ([]*entity.Service, error) {
	return u.serviceRepo.FindAll()
}
