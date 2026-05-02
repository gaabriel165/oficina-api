package part

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
)

type ListPartsUseCase struct {
	partRepo repository.PartRepository
}

func NewListPartsUseCase(partRepo repository.PartRepository) *ListPartsUseCase {
	return &ListPartsUseCase{partRepo: partRepo}
}

func (u *ListPartsUseCase) Execute() ([]*entity.Part, error) {
	return u.partRepo.FindAll()
}
