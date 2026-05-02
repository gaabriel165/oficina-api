package auth_test

import (
	"testing"

	"github.com/gabrielcamargo/oficina-api/internal/application/mocks"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/auth"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "test-secret"

func makeUser(t *testing.T) *entity.User {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	user, _ := entity.NewUser("admin@workshop.com", string(hash))
	return user
}

func TestCreateUser_ShouldCreateSuccessfully(t *testing.T) {
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.On("FindByEmail", "admin@workshop.com").Return(nil, repository.ErrUserNotFound)
	userRepo.On("Create", mock.Anything).Return(nil)

	uc := auth.NewCreateUserUseCase(userRepo)
	result, err := uc.Execute(auth.CreateUserInput{
		Email:    "admin@workshop.com",
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin@workshop.com", result.Email())
	userRepo.AssertExpectations(t)
}

func TestCreateUser_ShouldReturnErrorWhenEmailAlreadyRegistered(t *testing.T) {
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.On("FindByEmail", "admin@workshop.com").Return(makeUser(t), nil)

	uc := auth.NewCreateUserUseCase(userRepo)
	_, err := uc.Execute(auth.CreateUserInput{
		Email:    "admin@workshop.com",
		Password: "secret123",
	})

	assert.ErrorIs(t, err, usecase.ErrEmailAlreadyRegistered)
}

func TestLogin_ShouldReturnTokenOnValidCredentials(t *testing.T) {
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.On("FindByEmail", "admin@workshop.com").Return(makeUser(t), nil)

	uc := auth.NewLoginUseCase(userRepo, jwtSecret)
	result, err := uc.Execute(auth.LoginInput{
		Email:    "admin@workshop.com",
		Password: "secret123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
}

func TestLogin_ShouldReturnErrorWhenUserNotFound(t *testing.T) {
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.On("FindByEmail", "unknown@workshop.com").Return(nil, repository.ErrUserNotFound)

	uc := auth.NewLoginUseCase(userRepo, jwtSecret)
	_, err := uc.Execute(auth.LoginInput{
		Email:    "unknown@workshop.com",
		Password: "secret123",
	})

	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}

func TestLogin_ShouldReturnErrorWhenPasswordIsWrong(t *testing.T) {
	userRepo := mocks.NewMockUserRepository(t)

	userRepo.On("FindByEmail", "admin@workshop.com").Return(makeUser(t), nil)

	uc := auth.NewLoginUseCase(userRepo, jwtSecret)
	_, err := uc.Execute(auth.LoginInput{
		Email:    "admin@workshop.com",
		Password: "wrong-password",
	})

	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}
