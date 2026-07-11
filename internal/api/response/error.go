package response

import (
	"errors"
	"net/http"

	"github.com/gabrielcamargo/oficina-api/internal/application/usecase"
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	"github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/domain/valueobject"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(c *gin.Context, err error) {
	status := statusFromError(err)
	if status == http.StatusInternalServerError {
		c.JSON(status, ErrorResponse{Error: "internal server error"})
		return
	}
	c.JSON(status, ErrorResponse{Error: err.Error()})
}

func statusFromError(err error) int {
	notFound := []error{
		repository.ErrCustomerNotFound,
		repository.ErrVehicleNotFound,
		repository.ErrPartNotFound,
		repository.ErrServiceNotFound,
		repository.ErrServiceOrderNotFound,
		repository.ErrUserNotFound,
	}
	for _, e := range notFound {
		if errors.Is(err, e) {
			return http.StatusNotFound
		}
	}

	conflict := []error{
		usecase.ErrDocumentAlreadyRegistered,
		usecase.ErrPlateAlreadyRegistered,
		usecase.ErrEmailAlreadyRegistered,
	}
	for _, e := range conflict {
		if errors.Is(err, e) {
			return http.StatusConflict
		}
	}

	unprocessable := []error{
		entity.ErrCustomerNameRequired,
		entity.ErrCustomerPhoneRequired,
		entity.ErrCustomerEmailRequired,
		entity.ErrInvalidDocument,
		entity.ErrVehicleBrandRequired,
		entity.ErrVehicleModelRequired,
		entity.ErrVehicleInvalidYear,
		entity.ErrVehicleCustomerRequired,
		entity.ErrPartNameRequired,
		entity.ErrPartInvalidPrice,
		entity.ErrPartInvalidStock,
		entity.ErrPartInsufficientStock,
		entity.ErrServiceNameRequired,
		entity.ErrServiceInvalidLaborPrice,
		entity.ErrServiceInvalidEstimatedMin,
		entity.ErrServiceOrderCustomerRequired,
		entity.ErrServiceOrderVehicleRequired,
		entity.ErrServiceOrderVehicleNotOwnedByCustomer,
		entity.ErrServiceOrderInvalidTransition,
		entity.ErrServiceOrderNoItems,
		entity.ErrServiceOrderItemInvalidQty,
		usecase.ErrCNPJNotActive,
		valueobject.ErrInvalidCPF,
		valueobject.ErrInvalidCNPJ,
		valueobject.ErrInvalidPlate,
	}
	for _, e := range unprocessable {
		if errors.Is(err, e) {
			return http.StatusUnprocessableEntity
		}
	}

	if errors.Is(err, usecase.ErrInvalidCredentials) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
