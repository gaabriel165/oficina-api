package model

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
)

type UserModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;column:id"`
	Email        string    `gorm:"type:varchar(255);not null;uniqueIndex;column:email"`
	PasswordHash string    `gorm:"type:varchar(255);not null;column:password_hash"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string { return "users" }

func (m *UserModel) ToDomain() *entity.User {
	return entity.RestoreUser(m.ID, m.Email, m.PasswordHash, m.CreatedAt, m.UpdatedAt)
}

func UserModelFromDomain(u *entity.User) *UserModel {
	return &UserModel{
		ID:           u.ID(),
		Email:        u.Email(),
		PasswordHash: u.PasswordHash(),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
	}
}
