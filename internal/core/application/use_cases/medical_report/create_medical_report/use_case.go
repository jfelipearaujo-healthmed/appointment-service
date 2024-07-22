package create_medical_report_uc

import (
	"context"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/medical_report_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/appointment"
	medical_report_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/medical_report"
	create_medical_report_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/create_medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	appointmentRepository   appointment_repository_contract.Repository
	medicalReportRepository medical_report_repository_contract.Repository
}

func NewUseCase(
	appointmentRepository appointment_repository_contract.Repository,
	medicalReportRepository medical_report_repository_contract.Repository,
) create_medical_report_contract.UseCase {
	return &useCase{
		appointmentRepository:   appointmentRepository,
		medicalReportRepository: medicalReportRepository,
	}
}

func (uc *useCase) Execute(ctx context.Context, patientID, appointmentID uint, request *medical_report_dto.CreateMedicalReportRequest) (*entities.MedicalReport, error) {
	appointment, err := uc.appointmentRepository.GetByID(ctx, patientID, appointmentID, role.Doctor)
	if err != nil {
		return nil, err
	}

	if appointment.Status == entities.ScheduleInAnalysis || appointment.Status == entities.ReScheduleInAnalysis {
		return nil, app_error.New(http.StatusBadRequest, "the appointment must be confirmed or in progress before adding a medical report")
	}

	medicalReport := &entities.MedicalReport{
		AppointmentID: appointmentID,
		Comment:       request.Comment,
	}

	return uc.medicalReportRepository.Create(ctx, medicalReport)
}
