package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/service"
	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	createUC *service.CreateServiceUseCase
	updateUC *service.UpdateServiceUseCase
	deleteUC *service.DeleteServiceUseCase
	getUC    *service.GetServiceUseCase
	listUC   *service.ListServicesUseCase
}

func NewServiceHandler(
	createUC *service.CreateServiceUseCase,
	updateUC *service.UpdateServiceUseCase,
	deleteUC *service.DeleteServiceUseCase,
	getUC *service.GetServiceUseCase,
	listUC *service.ListServicesUseCase,
) *ServiceHandler {
	return &ServiceHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		getUC:    getUC,
		listUC:   listUC,
	}
}

func (h *ServiceHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/services", h.Create)
	router.GET("/services", h.List)
	router.GET("/services/:id", h.Get)
	router.PUT("/services/:id", h.Update)
	router.DELETE("/services/:id", h.Delete)
}

// Create godoc
// @Summary      Create a new service
// @Tags         services
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateServiceRequest  true  "Service data"
// @Success      201   {object}  dto.ServiceResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      422   {object}  response.ErrorResponse
// @Router       /services [post]
func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUC.Execute(service.CreateServiceInput{
		Name:             req.Name,
		Description:      req.Description,
		LaborPrice:       req.LaborPrice,
		EstimatedMinutes: req.EstimatedMinutes,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToServiceResponse(result))
}

// Update godoc
// @Summary      Update a service
// @Tags         services
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "Service ID"
// @Param        body  body      dto.UpdateServiceRequest  true  "Service data"
// @Success      200   {object}  dto.ServiceResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /services/{id} [put]
func (h *ServiceHandler) Update(c *gin.Context) {
	var req dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUC.Execute(service.UpdateServiceInput{
		ID:               c.Param("id"),
		Name:             req.Name,
		Description:      req.Description,
		LaborPrice:       req.LaborPrice,
		EstimatedMinutes: req.EstimatedMinutes,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToServiceResponse(result))
}

// Delete godoc
// @Summary      Delete a service
// @Tags         services
// @Security     BearerAuth
// @Param        id  path  string  true  "Service ID"
// @Success      204
// @Failure      404  {object}  response.ErrorResponse
// @Router       /services/{id} [delete]
func (h *ServiceHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary      Get a service by ID
// @Tags         services
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Service ID"
// @Success      200  {object}  dto.ServiceResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /services/{id} [get]
func (h *ServiceHandler) Get(c *gin.Context) {
	result, err := h.getUC.Execute(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToServiceResponse(result))
}

// List godoc
// @Summary      List all services
// @Tags         services
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.ServiceResponse
// @Router       /services [get]
func (h *ServiceHandler) List(c *gin.Context) {
	results, err := h.listUC.Execute()
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.ServiceResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.ToServiceResponse(r))
	}
	c.JSON(http.StatusOK, responses)
}
