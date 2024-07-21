package appointment_dto

type AppointmentCommonResponse struct {
	Message string `json:"message"`
}

func NewCreateAppointmentRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment schedule created, access your appointment later to see the status",
	}
}

func NewUpdateAppointmentRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment update requested, access your appointment later to see the status",
	}
}

func NewConfirmedAppointmentRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment confirmed, the patient will be notified",
	}
}

func NewCancelledAppointmentRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment cancelled, the patient will be notified",
	}
}
