package server

import (
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/feedback"
	medical_report_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/medical_report"
	cancel_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/cancel_appointment"
	confirm_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/confirm_appointment"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	get_appointment_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_by_id"
	list_appointments_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/list_appointments"
	update_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/update_appointment"
	create_feedback_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/create_feedback"
	get_feedback_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/get_feedback_by_id"
	list_feedbacks_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/list_feedbacks"
	create_medical_report_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/create_medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/cache"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

type Dependencies struct {
	Cache     cache.Cache
	DbService *persistence.DbService

	AppointmentTopic topic.TopicService
	FeedbackTopic    topic.TopicService

	AppointmentRepository   appointment_repository_contract.Repository
	EventRepository         event_repository_contract.Repository
	FeedbackRepository      feedback_repository_contract.Repository
	MedicalReportRepository medical_report_repository_contract.Repository

	CreateAppointmentUseCase  create_appointment_contract.UseCase
	GetAppointmentByIdUseCase get_appointment_by_id_contract.UseCase
	ListAppointmentsUseCase   list_appointments_contract.UseCase
	UpdateAppointmentUseCase  update_appointment_contract.UseCase
	ConfirmAppointmentUseCase confirm_appointment_contract.UseCase
	CancelAppointmentUseCase  cancel_appointment_contract.UseCase

	CreateFeedbackUseCase  create_feedback_contract.UseCase
	GetFeedbackByIdUseCase get_feedback_by_id_contract.UseCase
	ListFeedbacksUseCase   list_feedbacks_contract.UseCase

	CreateMedicalReportUseCase create_medical_report_contract.UseCase
}
