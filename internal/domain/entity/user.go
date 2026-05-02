package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserEmailRequired    = errors.New("user email is required")
	ErrUserPasswordRequired = errors.New("user password hash is required")
)

type User struct {
	id           string
	email        string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(email, passwordHash string) (*User, error) {
	if email == "" {
		return nil, ErrUserEmailRequired
	}
	if passwordHash == "" {
		return nil, ErrUserPasswordRequired
	}

	now := time.Now()
	return &User{
		id:           uuid.NewString(),
		email:        email,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func RestoreUser(id, email, passwordHash string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (u *User) ID() string           { return u.id }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }
