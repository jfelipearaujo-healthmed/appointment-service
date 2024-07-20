package appointment_dto

type CreateAppointmentRequest struct {
	ScheduleID uint   `json:"schedule_id" validate:"required,gt=0"`
	DoctorID   uint   `json:"doctor_id" validate:"required,gt=0"`
	DateTime   string `json:"date_time" validate:"required"`
}
