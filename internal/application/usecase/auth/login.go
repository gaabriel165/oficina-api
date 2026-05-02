package auth

import (
	"time"

	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type LoginUseCase struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewLoginUseCase(userRepo repository.UserRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (u *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	user, err := u.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, usecase.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash()), []byte(input.Password)); err != nil {
		return nil, usecase.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{Token: signed}, nil
}
