package repository

import "github.com/gabrielcamargo/oficina-api/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
