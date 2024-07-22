package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	cancel_appointment_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/cancel_appointment"
	confirm_appointment_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/confirm_appointment"
	create_appointment_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/create_appointment"
	get_appointment_by_id_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/get_appointment_by_id"
	list_appointments_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/list_appointments"
	update_appointment_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/appointment/update_appointment"
	create_feedback_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/feedback/create_feedback"
	get_feedback_by_id_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/feedback/get_feedback_by_id"
	list_feedbacks_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/feedback/list_feedbacks"
	create_medical_report_uc "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/application/use_cases/medical_report/create_medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/config"
	appointment_repository "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/repositories/appointment"
	event_repository "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/repositories/event"
	feedback_repository "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/repositories/feedback"
	medical_report_repository "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/repositories/medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/cancel_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/confirm_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/create_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/get_appointment_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/list_appointments"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/appointment/update_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/feedback/create_feedback"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/feedback/get_feedback_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/feedback/list_feedbacks"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/health"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/handlers/medical_report/create_medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/logger"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/token"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/secret"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
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

	appointmentTopic := topic.NewService(config.CloudConfig.AppointmentTopicName, cloudConfig)

	if err := appointmentTopic.UpdateTopicArn(ctx); err != nil {
		slog.ErrorContext(ctx, "error updating topic arn", "error", err)
		return nil, err
	}

	feedbackTopic := topic.NewService(config.CloudConfig.FeedbackTopicName, cloudConfig)

	if err := feedbackTopic.UpdateTopicArn(ctx); err != nil {
		slog.ErrorContext(ctx, "error updating topic arn", "error", err)
		return nil, err
	}

	cache := cache.NewRedisCache(ctx, config)

	appointmentRepository := appointment_repository.NewRepository(cache, dbService)
	eventRepository := event_repository.NewRepository(dbService)
	feedbackRepository := feedback_repository.NewRepository(dbService)
	medicalReportRepository := medical_report_repository.NewRepository(dbService)

	return &Server{
		Config: config,
		Dependencies: Dependencies{
			Cache:     cache,
			DbService: dbService,

			AppointmentTopic: appointmentTopic,
			FeedbackTopic:    feedbackTopic,

			AppointmentRepository:   appointmentRepository,
			EventRepository:         eventRepository,
			FeedbackRepository:      feedbackRepository,
			MedicalReportRepository: medicalReportRepository,

			CreateAppointmentUseCase:  create_appointment_uc.NewUseCase(appointmentTopic, eventRepository, config.ApiConfig.Location),
			GetAppointmentByIdUseCase: get_appointment_by_id_uc.NewUseCase(appointmentRepository),
			ListAppointmentsUseCase:   list_appointments_uc.NewUseCase(appointmentRepository),
			UpdateAppointmentUseCase: update_appointment_uc.NewUseCase(appointmentTopic,
				appointmentRepository,
				eventRepository,
				config.ApiConfig.Location),
			ConfirmAppointmentUseCase: confirm_appointment_uc.NewUseCase(appointmentRepository),
			CancelAppointmentUseCase:  cancel_appointment_uc.NewUseCase(appointmentRepository),

			CreateFeedbackUseCase:  create_feedback_uc.NewUseCase(feedbackTopic, appointmentRepository, eventRepository),
			GetFeedbackByIdUseCase: get_feedback_by_id_uc.NewUseCase(feedbackRepository),
			ListFeedbacksUseCase:   list_feedbacks_uc.NewUseCase(feedbackRepository),

			CreateMedicalReportUseCase: create_medical_report_uc.NewUseCase(appointmentRepository, medicalReportRepository),
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
	s.addFeedbackRoutes(api)
	s.addMedicalReportRoutes(api)

	return e
}

func (s *Server) addHealthCheckRoutes(e *echo.Echo) {
	healthHandler := health.NewHandler(s.DbService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) addAppointmentRoutes(g *echo.Group) {
	createAppointmentHandler := create_appointment.NewHandler(s.CreateAppointmentUseCase)
	getAppointmentByIdHandler := get_appointment_by_id.NewHandler(s.GetAppointmentByIdUseCase)
	listAppointmentsHandler := list_appointments.NewHandler(s.ListAppointmentsUseCase)
	updateAppointmentHandler := update_appointment.NewHandler(s.UpdateAppointmentUseCase)
	confirmAppointmentHandler := confirm_appointment.NewHandler(s.ConfirmAppointmentUseCase)
	cancelAppointmentHandler := cancel_appointment.NewHandler(s.CancelAppointmentUseCase)

	g.POST("/appointments", createAppointmentHandler.Handle, role.Middleware(role.Patient))
	g.GET("/appointments", listAppointmentsHandler.Handle, role.Middleware(role.Any))
	g.GET("/appointments/:appointmentId", getAppointmentByIdHandler.Handle, role.Middleware(role.Any))
	g.PUT("/appointments/:appointmentId", updateAppointmentHandler.Handle, role.Middleware(role.Patient))
	g.POST("/appointments/:appointmentId/confirm", confirmAppointmentHandler.Handle, role.Middleware(role.Doctor))
	g.POST("/appointments/:appointmentId/cancel", cancelAppointmentHandler.Handle, role.Middleware(role.Any))
}

func (s *Server) addFeedbackRoutes(g *echo.Group) {
	createFeedbackHandler := create_feedback.NewHandler(s.CreateFeedbackUseCase)
	getFeedbackByIdHandler := get_feedback_by_id.NewHandler(s.GetFeedbackByIdUseCase)
	listFeedbacksHandler := list_feedbacks.NewHandler(s.ListFeedbacksUseCase)

	g.POST("/appointments/:appointmentId/feedbacks", createFeedbackHandler.Handle, role.Middleware(role.Patient))
	g.GET("/appointments/:appointmentId/feedbacks", listFeedbacksHandler.Handle, role.Middleware(role.Any))
	g.GET("/appointments/:appointmentId/feedbacks/:feedbackId", getFeedbackByIdHandler.Handle, role.Middleware(role.Any))
}

func (s *Server) addMedicalReportRoutes(g *echo.Group) {
	createMedicalReportHandler := create_medical_report.NewHandler(s.CreateMedicalReportUseCase)

	g.POST("/appointments/:appointmentId/medical-reports", createMedicalReportHandler.Handle, role.Middleware(role.Doctor))
}
