package entities

import (
	"time"

	"gorm.io/gorm"
)

type FileAccess struct {
	gorm.Model

	UserID        uint `json:"user_id"`
	FileID        uint `json:"file_id"`
	DoctorID      uint `json:"doctor_id"`
	AppointmentID uint `json:"appointment_id"`

	ExpiresAt time.Time `json:"expires_at"`
}
