package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type PartRepository interface {
	Create(part *entity.Part) error
	Update(part *entity.Part) error
	Delete(id string) error
	FindByID(id string) (*entity.Part, error)
	FindAll() ([]*entity.Part, error)
}
