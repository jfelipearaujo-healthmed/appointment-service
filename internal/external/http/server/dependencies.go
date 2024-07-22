package server

import (
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/feedback"
	confirm_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/confirm_appointment"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	get_appointment_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_by_id"
	list_appointments_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/list_appointments"
	update_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/update_appointment"
	create_feedback_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/create_feedback"
	list_feedbacks_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/list_feedbacks"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

type Dependencies struct {
	Cache     cache.Cache
	DbService *persistence.DbService

	AppointmentTopic topic.TopicService
	FeedbackTopic    topic.TopicService

	AppointmentRepository appointment_repository_contract.Repository
	EventRepository       event_repository_contract.Repository
	FeedbackRepository    feedback_repository_contract.Repository

	CreateAppointmentUseCase  create_appointment_contract.UseCase
	GetAppointmentByIdUseCase get_appointment_by_id_contract.UseCase
	ListAppointmentsUseCase   list_appointments_contract.UseCase
	UpdateAppointmentUseCase  update_appointment_contract.UseCase
	ConfirmAppointmentUseCase confirm_appointment_contract.UseCase

	CreateFeedbackUseCase create_feedback_contract.UseCase
	ListFeedbacksUseCase  list_feedbacks_contract.UseCase
}
