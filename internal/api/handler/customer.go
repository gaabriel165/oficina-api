package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/customer"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	createUC *customer.CreateCustomerUseCase
	updateUC *customer.UpdateCustomerUseCase
	deleteUC *customer.DeleteCustomerUseCase
	getUC    *customer.GetCustomerUseCase
	listUC   *customer.ListCustomersUseCase
}

func NewCustomerHandler(
	createUC *customer.CreateCustomerUseCase,
	updateUC *customer.UpdateCustomerUseCase,
	deleteUC *customer.DeleteCustomerUseCase,
	getUC *customer.GetCustomerUseCase,
	listUC *customer.ListCustomersUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		getUC:    getUC,
		listUC:   listUC,
	}
}

func (h *CustomerHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/customers", h.Create)
	router.GET("/customers", h.List)
	router.GET("/customers/:id", h.Get)
	router.PUT("/customers/:id", h.Update)
	router.DELETE("/customers/:id", h.Delete)
}

// Create godoc
// @Summary      Create a new customer
// @Tags         customers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateCustomerRequest  true  "Customer data"
// @Success      201   {object}  dto.CustomerResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      409   {object}  response.ErrorResponse
// @Failure      422   {object}  response.ErrorResponse
// @Router       /customers [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUC.Execute(customer.CreateCustomerInput{
		Name:     req.Name,
		Document: req.Document,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToCustomerResponse(result))
}

// Update godoc
// @Summary      Update a customer
// @Tags         customers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "Customer ID"
// @Param        body  body      dto.UpdateCustomerRequest  true  "Customer data"
// @Success      200   {object}  dto.CustomerResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /customers/{id} [put]
func (h *CustomerHandler) Update(c *gin.Context) {
	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUC.Execute(customer.UpdateCustomerInput{
		ID:    c.Param("id"),
		Name:  req.Name,
		Phone: req.Phone,
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToCustomerResponse(result))
}

// Delete godoc
// @Summary      Delete a customer
// @Tags         customers
// @Security     BearerAuth
// @Param        id  path  string  true  "Customer ID"
// @Success      204
// @Failure      404  {object}  response.ErrorResponse
// @Router       /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary      Get a customer by ID
// @Tags         customers
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Customer ID"
// @Success      200  {object}  dto.CustomerResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /customers/{id} [get]
func (h *CustomerHandler) Get(c *gin.Context) {
	result, err := h.getUC.Execute(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToCustomerResponse(result))
}

// List godoc
// @Summary      List all customers
// @Tags         customers
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.CustomerResponse
// @Router       /customers [get]
func (h *CustomerHandler) List(c *gin.Context) {
	results, err := h.listUC.Execute()
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.CustomerResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.ToCustomerResponse(r))
	}
	c.JSON(http.StatusOK, responses)
}
