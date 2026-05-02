package auth

import (
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Email    string
	Password string
}

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepo: userRepo}
}

func (u *CreateUserUseCase) Execute(input CreateUserInput) (*entity.User, error) {
	_, err := u.userRepo.FindByEmail(input.Email)
	if err == nil {
		return nil, usecase.ErrEmailAlreadyRegistered
	}
	if err != repository.ErrUserNotFound {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := entity.NewUser(input.Email, string(hash))
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
