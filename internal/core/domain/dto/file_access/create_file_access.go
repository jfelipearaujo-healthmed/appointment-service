package file_access_dto

type CreateFileAccess struct {
	FileID        uint
	DoctorID      uint   `json:"doctor_id" validate:"required"`
	AppointmentID uint   `json:"appointment_id" validate:"required"`
	ExpiresAt     string `json:"expires_at" validate:"required"`
}
