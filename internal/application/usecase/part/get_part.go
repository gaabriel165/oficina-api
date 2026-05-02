package part

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type GetPartUseCase struct {
	partRepo repository.PartRepository
}

func NewGetPartUseCase(partRepo repository.PartRepository) *GetPartUseCase {
	return &GetPartUseCase{partRepo: partRepo}
}

func (u *GetPartUseCase) Execute(id string) (*entity.Part, error) {
	return u.partRepo.FindByID(id)
}
