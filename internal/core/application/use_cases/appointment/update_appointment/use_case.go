package update_appointment_uc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/event"
	update_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/update_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

const (
	dateTimeLayout = "2006-01-02 15:04"
)

type useCase struct {
	topicService          topic.TopicService
	appointmentRepository appointment_repository_contract.Repository
	eventRepository       event_repository_contract.Repository
	location              *time.Location
}

func NewUseCase(
	topicService topic.TopicService,
	appointmentRepository appointment_repository_contract.Repository,
	eventRepository event_repository_contract.Repository,
	location *time.Location,
) update_appointment_contract.UseCase {
	return &useCase{
		topicService:          topicService,
		appointmentRepository: appointmentRepository,
		location:              location,
	}
}

func (uc *useCase) Execute(ctx context.Context, patientID, appointmentID uint, request *appointment_dto.UpdateAppointmentRequest) error {
	parsedTime, err := time.ParseInLocation(dateTimeLayout, request.DateTime, uc.location)
	if err != nil {
		return app_error.New(http.StatusBadRequest, "unable to parse the date and time provided")
	}

	year, month, day := parsedTime.Date()
	hour, minute, _ := parsedTime.Clock()

	finalTime := time.Date(year, month, day, hour, minute, 0, 0, uc.location)

	if finalTime.Before(time.Now()) {
		return app_error.New(http.StatusBadRequest, "date and time must be in the future")
	}

	appointment, err := uc.appointmentRepository.GetByID(ctx, patientID, appointmentID, role.Patient)
	if err != nil {
		return err
	}

	if appointment.Status == entities.Concluded || appointment.Status == entities.Cancelled {
		return app_error.New(http.StatusBadRequest, "appointment concluded or cancelled and cannot be re-scheduled")
	}

	if appointment.Status == entities.ReScheduleInAnalysis {
		return app_error.New(http.StatusBadRequest, "appointment already in re-scheduling state, wait for the appointment to conclude")
	}

	appointment.ScheduleID = request.ScheduleID
	appointment.PatientID = patientID
	appointment.DoctorID = request.DoctorID
	appointment.DateTime = finalTime
	appointment.Status = entities.ReScheduleInAnalysis

	appointmentData, err := json.Marshal(appointment)
	if err != nil {
		return app_error.New(http.StatusInternalServerError, "unable to marshal the appointment data")
	}

	event := &entities.Event{
		UserID:    patientID,
		Data:      string(appointmentData),
		EventType: events.UpdateAppointment,
	}

	existingEvent, err := uc.eventRepository.GetByIDsAndDateTime(ctx, event)
	if err != nil && !app_error.IsAppError(err) {
		return err
	}

	if existingEvent != nil {
		return app_error.New(http.StatusBadRequest, "re-schedule already requested")
	}

	if _, err := uc.appointmentRepository.Update(ctx, patientID, appointment); err != nil {
		return err
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
