package get_medical_report_by_id_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	medical_report_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/medical_report"
	get_medical_report_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/get_medical_report_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository medical_report_repository_contract.Repository
}

func NewUseCase(repository medical_report_repository_contract.Repository) get_medical_report_by_id_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID, medicalReportID uint, roleName role.Role) (*entities.MedicalReport, error) {
	return uc.repository.GetByID(ctx, userID, appointmentID, medicalReportID, roleName)
}
