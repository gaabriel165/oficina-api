package api

import (
	"fmt"

	_ "github.com/gabrielcamargo/oficina-api/docs"
	"github.com/gabrielcamargo/oficina-api/configs"
	"github.com/gabrielcamargo/oficina-api/internal/api/handler"
	"github.com/gabrielcamargo/oficina-api/internal/api/middleware"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/auth"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/customer"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/part"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/service"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/serviceorder"
	"github.com/gabrielcamargo/oficina-api/internal/application/usecase/vehicle"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/external"
	infrarepo "github.com/gabrielcamargo/oficina-api/internal/infrastructure/repository"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	config *configs.Config
	db     *gorm.DB
}

func NewServer(cfg *configs.Config, db *gorm.DB) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		config: cfg,
		db:     db,
	}

	server.registerRoutes()

	return server
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.ServerPort))
}

func (s *Server) registerRoutes() {
	userRepo := infrarepo.NewGormUserRepository(s.db)
	customerRepo := infrarepo.NewGormCustomerRepository(s.db)
	vehicleRepo := infrarepo.NewGormVehicleRepository(s.db)
	partRepo := infrarepo.NewGormPartRepository(s.db)
	serviceRepo := infrarepo.NewGormServiceRepository(s.db)
	orderRepo := infrarepo.NewGormServiceOrderRepository(s.db)
	cnpjSvc := external.NewBrasilAPIClient()

	authHandler := handler.NewAuthHandler(
		auth.NewCreateUserUseCase(userRepo),
		auth.NewLoginUseCase(userRepo, s.config.JWTSecret),
	)

	customerHandler := handler.NewCustomerHandler(
		customer.NewCreateCustomerUseCase(customerRepo, cnpjSvc),
		customer.NewUpdateCustomerUseCase(customerRepo),
		customer.NewDeleteCustomerUseCase(customerRepo),
		customer.NewGetCustomerUseCase(customerRepo),
		customer.NewListCustomersUseCase(customerRepo),
	)

	vehicleHandler := handler.NewVehicleHandler(
		vehicle.NewCreateVehicleUseCase(vehicleRepo, customerRepo),
		vehicle.NewUpdateVehicleUseCase(vehicleRepo),
		vehicle.NewDeleteVehicleUseCase(vehicleRepo),
		vehicle.NewGetVehicleUseCase(vehicleRepo),
		vehicle.NewListVehiclesUseCase(vehicleRepo),
	)

	partHandler := handler.NewPartHandler(
		part.NewCreatePartUseCase(partRepo),
		part.NewUpdatePartUseCase(partRepo),
		part.NewDeletePartUseCase(partRepo),
		part.NewGetPartUseCase(partRepo),
		part.NewListPartsUseCase(partRepo),
		part.NewUpdateStockUseCase(partRepo),
	)

	serviceHandler := handler.NewServiceHandler(
		service.NewCreateServiceUseCase(serviceRepo),
		service.NewUpdateServiceUseCase(serviceRepo),
		service.NewDeleteServiceUseCase(serviceRepo),
		service.NewGetServiceUseCase(serviceRepo),
		service.NewListServicesUseCase(serviceRepo),
	)

	orderHandler := handler.NewServiceOrderHandler(
		serviceorder.NewCreateServiceOrderUseCase(orderRepo, customerRepo, vehicleRepo),
		serviceorder.NewStartDiagnosisUseCase(orderRepo),
		serviceorder.NewAddServiceToOrderUseCase(orderRepo, serviceRepo),
		serviceorder.NewAddPartToOrderUseCase(orderRepo, partRepo),
		serviceorder.NewSendBudgetUseCase(orderRepo),
		serviceorder.NewApproveBudgetUseCase(orderRepo, partRepo),
		serviceorder.NewRejectBudgetUseCase(orderRepo),
		serviceorder.NewFinishExecutionUseCase(orderRepo),
		serviceorder.NewDeliverVehicleUseCase(orderRepo),
		serviceorder.NewGetServiceOrderUseCase(orderRepo),
		serviceorder.NewListServiceOrdersUseCase(orderRepo),
		serviceorder.NewGetExecutionMetricsUseCase(orderRepo),
	)

	s.router.GET("/health", func(ctx *gin.Context) { ctx.Status(200) })
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := s.router.Group("/api/v1")
	protected := v1.Group("", middleware.Auth(s.config.JWTSecret))

	authHandler.RegisterRoutes(v1)
	orderHandler.RegisterPublicRoutes(v1)
	customerHandler.RegisterRoutes(protected)
	vehicleHandler.RegisterRoutes(protected)
	partHandler.RegisterRoutes(protected)
	serviceHandler.RegisterRoutes(protected)
	orderHandler.RegisterRoutes(protected)
}
