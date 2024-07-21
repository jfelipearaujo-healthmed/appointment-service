package appointment_dto

type ConfirmAppointmentRequest struct {
	Confirmed bool `json:"confirmed" validate:"required"`
}
