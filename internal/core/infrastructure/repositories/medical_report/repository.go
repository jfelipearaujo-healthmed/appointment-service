package medical_report_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
	medical_report_repository_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/repositories/medical_report"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) medical_report_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, userID, appointmentID, medicalReportID uint, roleName role.Role) (*entities.MedicalReport, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	medicalReport := new(entities.MedicalReport)

	query := "medical_reports.id = ? AND appointments.patient_id = ? AND appointments.id = ?"

	if roleName == role.Doctor {
		query = "medical_reports.id = ? AND appointments.doctor_id = ? AND appointments.id = ?"
	}

	if err := tx.Preload("Appointment").
		Order("medical_reports.created_at DESC").
		Joins("JOIN appointments ON appointments.id = medical_reports.appointment_id").
		Where(query, medicalReportID, userID, appointmentID).
		Find(&medicalReport).Error; err != nil {
		return nil, err
	}

	return medicalReport, nil
}

func (rp *repository) GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.MedicalReport, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	medicalReport := new(entities.MedicalReport)

	if err := tx.Where("appointment_id = ?", appointmentID).First(medicalReport).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("medicalReport with appointment id %d not found", appointmentID))
		}

		return nil, err
	}

	return medicalReport, nil
}

func (rp *repository) List(ctx context.Context, userID, appointmentID uint, roleName role.Role) ([]entities.MedicalReport, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	medicalReports := new([]entities.MedicalReport)

	query := "appointments.patient_id = ? AND appointments.id = ?"

	if roleName == role.Doctor {
		query = "appointments.doctor_id = ? AND appointments.id = ?"
	}

	if err := tx.Preload("Appointment").
		Order("medical_reports.created_at DESC").
		Joins("JOIN appointments ON appointments.id = medical_reports.appointment_id").
		Where(query, userID, appointmentID).
		Find(&medicalReports).Error; err != nil {
		return nil, err
	}

	return *medicalReports, nil
}

func (rp *repository) Create(ctx context.Context, medicalReport *entities.MedicalReport) (*entities.MedicalReport, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(medicalReport).Error; err != nil {
		return nil, err
	}

	return medicalReport, nil
}

func (rp *repository) Update(ctx context.Context, medicalReport *entities.MedicalReport) (*entities.MedicalReport, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(medicalReport).Where("id = ?", medicalReport.ID).Updates(medicalReport).Error; err != nil {
		return nil, err
	}

	return medicalReport, nil
}

func (rp *repository) Delete(ctx context.Context, medicalReportID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.MedicalReport{}, "id = ?", medicalReportID).Error; err != nil {
		return err
	}

	return nil
}
