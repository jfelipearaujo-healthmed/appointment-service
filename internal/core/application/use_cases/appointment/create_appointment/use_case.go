package create_appointment_uc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/events"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

const (
	dateTimeLayout = "2006-01-02 15:04"
)

type useCase struct {
	topicService topic.TopicService
	repository   event_repository_contract.Repository
	location     *time.Location
}

func NewUseCase(
	topicService topic.TopicService,
	repository event_repository_contract.Repository,
	location *time.Location,
) create_appointment_contract.UseCase {
	return &useCase{
		topicService: topicService,
		repository:   repository,
		location:     location,
	}
}

func (uc *useCase) Execute(ctx context.Context, patientID uint, request *appointment_dto.CreateAppointmentRequest) (*entities.Event, error) {
	parsedTime, err := time.ParseInLocation(dateTimeLayout, request.DateTime, uc.location)
	if err != nil {
		return nil, app_error.New(http.StatusBadRequest, "unable to parse the date and time provided")
	}

	year, month, day := parsedTime.Date()
	hour, minute, _ := parsedTime.Clock()

	finalTime := time.Date(year, month, day, hour, minute, 0, 0, uc.location)

	if finalTime.Before(time.Now()) {
		return nil, app_error.New(http.StatusBadRequest, "date and time must be in the future")
	}

	appointment := &entities.Appointment{
		ScheduleID: request.ScheduleID,
		PatientID:  patientID,
		DoctorID:   request.DoctorID,
		DateTime:   finalTime,
	}

	appointmentData, err := json.Marshal(appointment)
	if err != nil {
		return nil, app_error.New(http.StatusInternalServerError, "unable to marshal the appointment data")
	}

	event := &entities.Event{
		UserID:    patientID,
		Data:      string(appointmentData),
		EventType: events.CreateAppointment,
	}

	existingEvent, err := uc.repository.GetByIDsAndDateTime(ctx, event)
	if err != nil && !app_error.IsAppError(err) {
		return nil, err
	}

	// if the event was not processed yet, we can resend it
	tolerance := time.Minute * 2

	if existingEvent != nil && existingEvent.CreatedAt.Add(tolerance).Before(time.Now()) {
		return nil, app_error.New(http.StatusBadRequest, "schedule already requested")
	}

	messageId, err := uc.topicService.Publish(ctx, topic.NewMessage(event.EventType, event))
	if err != nil {
		return nil, err
	}

	event.MessageID = *messageId

	if existingEvent != nil {
		event, err = uc.repository.Update(ctx, event)
		if err != nil {
			return nil, err
		}

		return event, nil
	}

	event, err = uc.repository.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
