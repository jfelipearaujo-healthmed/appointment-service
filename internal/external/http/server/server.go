package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	create_appointment_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/create_appointment"
	list_appointments_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/list_appointments"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/config"
	appointment_repository "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/repositories/appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/create_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/list_appointments"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/health"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/logger"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/token"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/secret"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config *config.Config

	Dependencies
}

func NewServer(ctx context.Context, config *config.Config) (*Server, error) {
	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error getting aws config", "error", err)
		return nil, err
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	secretService := secret.NewService(cloudConfig)

	dbUrl, err := secretService.GetSecret(ctx, config.DbConfig.UrlSecretName)
	if err != nil {
		slog.ErrorContext(ctx, "error getting secret", "secret_name", config.DbConfig.UrlSecretName, "error", err)
		return nil, err
	}

	cacheUrl, err := secretService.GetSecret(ctx, config.CacheConfig.HostSecretName)
	if err != nil {
		slog.ErrorContext(ctx, "error getting secret", "secret_name", config.CacheConfig.HostSecretName, "error", err)
		return nil, err
	}

	config.DbConfig.Url = dbUrl
	config.CacheConfig.Host = cacheUrl

	dbService := persistence.NewDbService()

	if err := dbService.Connect(config); err != nil {
		slog.ErrorContext(ctx, "error connecting to database", "error", err)
		return nil, err
	}

	cache := cache.NewRedisCache(ctx, config)

	appointmentRepository := appointment_repository.NewRepository(cache, dbService)

	return &Server{
		Config: config,
		Dependencies: Dependencies{
			Cache:     cache,
			DbService: dbService,

			AppointmentRepository: appointmentRepository,

			CreateAppointmentUseCase: create_appointment_uc.NewUseCase(appointmentRepository, config.ApiConfig.Location),
			ListAppointmentsUseCase:  list_appointments_uc.NewUseCase(appointmentRepository),
		},
	}, nil
}

func (s *Server) GetServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(logger.Middleware())
	e.Use(middleware.Recover())

	s.addHealthCheckRoutes(e)

	api := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	api.Use(token.Middleware())
	s.addAppointmentRoutes(api)

	return e
}

func (s *Server) addHealthCheckRoutes(e *echo.Echo) {
	healthHandler := health.NewHandler(s.DbService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) addAppointmentRoutes(g *echo.Group) {
	createAppointmentHandler := create_appointment.NewHandler(s.CreateAppointmentUseCase)
	listAppointmentsHandler := list_appointments.NewHandler(s.ListAppointmentsUseCase)

	g.POST("/appointments", createAppointmentHandler.Handle, role.Middleware(role.Patient))
	g.GET("/appointments", listAppointmentsHandler.Handle, role.Middleware(role.Any))
}
