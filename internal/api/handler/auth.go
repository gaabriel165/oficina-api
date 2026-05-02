package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	createUserUC *auth.CreateUserUseCase
	loginUC      *auth.LoginUseCase
}

func NewAuthHandler(createUserUC *auth.CreateUserUseCase, loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{createUserUC: createUserUC, loginUC: loginUC}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/register", h.Register)
	router.POST("/auth/login", h.Login)
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserRequest  true  "User credentials"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  response.ErrorResponse
// @Failure      409   {object}  response.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.createUserUC.Execute(auth.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// Login godoc
// @Summary      Authenticate user and get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login credentials"
// @Success      200   {object}  dto.TokenResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.loginUC.Execute(auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{Token: output.Token})
}
