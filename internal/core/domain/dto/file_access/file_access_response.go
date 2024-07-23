package file_access_dto

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type FileAccessResponse struct {
	ID uint `json:"id"`

	FileID        uint `json:"file_id"`
	DoctorID      uint `json:"doctor_id"`
	AppointmentID uint `json:"appointment_id"`

	ExpiresAt time.Time `json:"expires_at"`
}

func MapFromDomain(file *entities.FileAccess) *FileAccessResponse {
	return &FileAccessResponse{
		ID: file.ID,

		FileID:        file.FileID,
		DoctorID:      file.DoctorID,
		AppointmentID: file.AppointmentID,

		ExpiresAt: file.ExpiresAt,
	}
}
