package create_feedback_uc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/feedback_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	create_feedback_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/create_feedback"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

type useCase struct {
	topicService          topic.TopicService
	appointmentRepository appointment_repository_contract.Repository
	eventRepository       event_repository_contract.Repository
}

func NewUseCase(
	topicService topic.TopicService,
	appointmentRepository appointment_repository_contract.Repository,
	eventRepository event_repository_contract.Repository,
) create_feedback_contract.UseCase {
	return &useCase{
		topicService:          topicService,
		appointmentRepository: appointmentRepository,
		eventRepository:       eventRepository,
	}
}

func (uc *useCase) Execute(ctx context.Context, patientID uint, appointmentID uint, request *feedback_dto.CreateFeedbackRequest) error {
	appointment, err := uc.appointmentRepository.GetByID(ctx, patientID, appointmentID, role.Patient)
	if err != nil {
		return err
	}

	if appointment.Status != entities.Concluded {
		return app_error.New(http.StatusBadRequest, "the appointment must be concluded before adding a feedback")
	}

	feedback := &entities.Feedback{
		AppointmentID: appointmentID,
		Rating:        request.Rating,
		Comment:       request.Comment,
	}

	feedbackData, err := json.Marshal(feedback)
	if err != nil {
		return app_error.New(http.StatusInternalServerError, "unable to marshal the feedback data")
	}

	event := &entities.Event{
		UserID:    patientID,
		Data:      string(feedbackData),
		EventType: events.CreateFeedback,
	}

	existingEvent, err := uc.eventRepository.GetByIDsAndDateTime(ctx, event)
	if err != nil && !app_error.IsAppError(err) {
		return err
	}

	if existingEvent != nil {
		return app_error.New(http.StatusBadRequest, "feedback already sent, please wait for the feedback to be processed")
	}

	messageId, err := uc.topicService.Publish(ctx, topic.NewMessage(event.EventType, event))
	if err != nil {
		return err
	}

	event.MessageID = *messageId

	if _, err := uc.eventRepository.Create(ctx, event); err != nil {
		return err
	}

	return nil
}
