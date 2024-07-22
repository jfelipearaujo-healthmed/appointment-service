package create_medical_report_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/medical_report_dto"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type UseCase interface {
	Execute(ctx context.Context, patientID, appointmentID uint, request *medical_report_dto.CreateMedicalReportRequest) (*entities.MedicalReport, error)
}
