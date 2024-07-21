package create_appointment_uc

import (
	"context"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
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

func NewUseCase(topicService topic.TopicService,
	repository appointment_repository_contract.Repository,
	location *time.Location,
) create_appointment_contract.UseCase {
	return &useCase{
		topicService: topicService,
		repository:   repository,
		location:     location,
	}
}

func (uc *useCase) Execute(ctx context.Context, patientID uint, request *appointment_dto.CreateAppointmentRequest) (*entities.Appointment, error) {
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
		Status:     entities.ScheduleInAnalysis,
	}

	existingAppointment, err := uc.repository.GetByIDsAndDateTime(ctx, request.ScheduleID, patientID, request.DoctorID, finalTime)
	if err != nil && !app_error.IsAppError(err) {
		return nil, err
	}

	if existingAppointment != nil {
		return nil, app_error.New(http.StatusBadRequest, "appointment already exists")
	}

	appointment, err = uc.repository.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	if _, err := uc.topicService.Publish(ctx, topic.NewMessage(events.CreateAppointment, appointment)); err != nil {
		return nil, err
	}

	return appointment, nil
}
