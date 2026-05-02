package part

import "github.com/gabrielcamargo/oficina-api/internal/domain/repository"

type DeletePartUseCase struct {
	partRepo repository.PartRepository
}

func NewDeletePartUseCase(partRepo repository.PartRepository) *DeletePartUseCase {
	return &DeletePartUseCase{partRepo: partRepo}
}

func (u *DeletePartUseCase) Execute(id string) error {
	if _, err := u.partRepo.FindByID(id); err != nil {
		return err
	}
	return u.partRepo.Delete(id)
}
