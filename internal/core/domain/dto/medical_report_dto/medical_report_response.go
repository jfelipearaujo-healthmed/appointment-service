package medical_report_dto

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type MedicalReportResponse struct {
	ID uint `json:"id"`

	AppointmentID uint   `json:"appointment_id"`
	Comment       string `json:"comment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapFromDomain(medicalReport *entities.MedicalReport) *MedicalReportResponse {
	return &MedicalReportResponse{
		ID: medicalReport.ID,

		AppointmentID: medicalReport.AppointmentID,
		Comment:       medicalReport.Comment,

		CreatedAt: medicalReport.CreatedAt,
		UpdatedAt: medicalReport.UpdatedAt,
	}
}

func MapFromDomainSlice(medicalReports []entities.MedicalReport) []*MedicalReportResponse {
	mapped := make([]*MedicalReportResponse, len(medicalReports))

	for i := range medicalReports {
		mapped[i] = MapFromDomain(&medicalReports[i])
	}

	return mapped
}
