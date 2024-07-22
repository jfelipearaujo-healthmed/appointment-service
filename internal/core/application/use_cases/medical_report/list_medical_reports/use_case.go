package list_medical_reports_uc

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	medical_report_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/medical_report"
	list_medical_reports_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/list_medical_reports"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type useCase struct {
	repository medical_report_repository_contract.Repository
}

func NewUseCase(repository medical_report_repository_contract.Repository) list_medical_reports_contract.UseCase {
	return &useCase{
		repository: repository,
	}
}

func (uc *useCase) Execute(ctx context.Context, userID, appointmentID uint, roleName role.Role) ([]entities.MedicalReport, error) {
	return uc.repository.List(ctx, userID, appointmentID, roleName)
}
