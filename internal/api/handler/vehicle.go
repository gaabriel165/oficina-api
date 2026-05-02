package handler

import (
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/api/dto"
	"github.com/gabrielcamargo/oficina-api/internal/api/response"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/vehicle"
	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	createUC *vehicle.CreateVehicleUseCase
	updateUC *vehicle.UpdateVehicleUseCase
	deleteUC *vehicle.DeleteVehicleUseCase
	getUC    *vehicle.GetVehicleUseCase
	listUC   *vehicle.ListVehiclesUseCase
}

func NewVehicleHandler(
	createUC *vehicle.CreateVehicleUseCase,
	updateUC *vehicle.UpdateVehicleUseCase,
	deleteUC *vehicle.DeleteVehicleUseCase,
	getUC *vehicle.GetVehicleUseCase,
	listUC *vehicle.ListVehiclesUseCase,
) *VehicleHandler {
	return &VehicleHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		getUC:    getUC,
		listUC:   listUC,
	}
}

func (h *VehicleHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/vehicles", h.Create)
	router.GET("/vehicles", h.List)
	router.GET("/vehicles/:id", h.Get)
	router.PUT("/vehicles/:id", h.Update)
	router.DELETE("/vehicles/:id", h.Delete)
}

// Create godoc
// @Summary      Create a new vehicle
// @Tags         vehicles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateVehicleRequest  true  "Vehicle data"
// @Success      201   {object}  dto.VehicleResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      409   {object}  response.ErrorResponse
// @Router       /vehicles [post]
func (h *VehicleHandler) Create(c *gin.Context) {
	var req dto.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createUC.Execute(vehicle.CreateVehicleInput{
		CustomerID: req.CustomerID,
		Plate:      req.Plate,
		Brand:      req.Brand,
		Model:      req.Model,
		Year:       req.Year,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ToVehicleResponse(result))
}

// Update godoc
// @Summary      Update a vehicle
// @Tags         vehicles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "Vehicle ID"
// @Param        body  body      dto.UpdateVehicleRequest  true  "Vehicle data"
// @Success      200   {object}  dto.VehicleResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /vehicles/{id} [put]
func (h *VehicleHandler) Update(c *gin.Context) {
	var req dto.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateUC.Execute(vehicle.UpdateVehicleInput{
		ID:    c.Param("id"),
		Brand: req.Brand,
		Model: req.Model,
		Year:  req.Year,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToVehicleResponse(result))
}

// Delete godoc
// @Summary      Delete a vehicle
// @Tags         vehicles
// @Security     BearerAuth
// @Param        id  path  string  true  "Vehicle ID"
// @Success      204
// @Failure      404  {object}  response.ErrorResponse
// @Router       /vehicles/{id} [delete]
func (h *VehicleHandler) Delete(c *gin.Context) {
	if err := h.deleteUC.Execute(c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Get godoc
// @Summary      Get a vehicle by ID
// @Tags         vehicles
// @Security     BearerAuth
// @Produce      json
// @Param        id  path      string  true  "Vehicle ID"
// @Success      200  {object}  dto.VehicleResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /vehicles/{id} [get]
func (h *VehicleHandler) Get(c *gin.Context) {
	result, err := h.getUC.Execute(c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToVehicleResponse(result))
}

// List godoc
// @Summary      List all vehicles
// @Tags         vehicles
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   dto.VehicleResponse
// @Router       /vehicles [get]
func (h *VehicleHandler) List(c *gin.Context) {
	results, err := h.listUC.Execute()
	if err != nil {
		response.Error(c, err)
		return
	}

	responses := make([]dto.VehicleResponse, 0, len(results))
	for _, r := range results {
		responses = append(responses, dto.ToVehicleResponse(r))
	}
	c.JSON(http.StatusOK, responses)
}
