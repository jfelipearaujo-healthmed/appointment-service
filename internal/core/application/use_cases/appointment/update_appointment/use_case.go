package update_appointment_uc

import (
	"context"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	update_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/update_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
)

const (
	dateTimeLayout = "2006-01-02 15:04"
)

type useCase struct {
	topicService topic.TopicService
	repository   appointment_repository_contract.Repository
	location     *time.Location
}

func NewUseCase(
	topicService topic.TopicService,
	repository appointment_repository_contract.Repository,
	location *time.Location,
) update_appointment_contract.UseCase {
	return &useCase{
		topicService: topicService,
		repository:   repository,
		location:     location,
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

	appointment, err := uc.repository.GetByID(ctx, patientID, appointmentID, role.Patient)
	if err != nil {
		return err
	}

	if appointment.Status == entities.Concluded || appointment.Status == entities.Cancelled {
		return app_error.New(http.StatusBadRequest, "appointment concluded or cancelled and cannot be re-scheduled")
	}

	if appointment.Status == entities.ReScheduleInAnalysis {
		return app_error.New(http.StatusBadRequest, "appointment already in re-scheduling state, wait for the appointment to conclude")
	}

	appointment.Status = entities.ReScheduleInAnalysis

	if _, err := uc.repository.Update(ctx, patientID, appointment); err != nil {
		return err
	}

	reSchedule := &entities.Appointment{
		ScheduleID: request.ScheduleID,
		PatientID:  patientID,
		DoctorID:   request.DoctorID,
		DateTime:   finalTime,
	}

	if _, err := uc.topicService.Publish(ctx, topic.NewMessage(events.UpdateAppointment, reSchedule)); err != nil {
		return err
	}

	return nil
}
