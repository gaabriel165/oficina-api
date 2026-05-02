package repository

import (
	"errors"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *entity.User) error {
	return r.db.Create(model.UserModelFromDomain(user)).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var m model.UserModel
	err := r.db.First(&m, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainrepo.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}
