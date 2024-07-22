package medical_report_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
)

type Repository interface {
	GetByID(ctx context.Context, userID, appointmentID, medicalReportID uint, roleName role.Role) (*entities.MedicalReport, error)
	GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.MedicalReport, error)
	List(ctx context.Context, userID, appointmentID uint, roleName role.Role) ([]entities.MedicalReport, error)
	Create(ctx context.Context, medicalReport *entities.MedicalReport) (*entities.MedicalReport, error)
	Update(ctx context.Context, medicalReport *entities.MedicalReport) (*entities.MedicalReport, error)
	Delete(ctx context.Context, medicalReportID uint) error
}
