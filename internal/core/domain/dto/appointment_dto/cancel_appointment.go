package appointment_dto

type CancelAppointmentRequest struct {
	Reason string `json:"reason" validate:"required,min=1,max=255"`
}
