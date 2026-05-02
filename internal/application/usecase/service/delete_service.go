package service

import "github.com/gabrielcamargo/oficina-api/internal/domain/repository"

type DeleteServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

func NewDeleteServiceUseCase(serviceRepo repository.ServiceRepository) *DeleteServiceUseCase {
	return &DeleteServiceUseCase{serviceRepo: serviceRepo}
}

func (u *DeleteServiceUseCase) Execute(id string) error {
	if _, err := u.serviceRepo.FindByID(id); err != nil {
		return err
	}
	return u.serviceRepo.Delete(id)
}
