package appointment_dto

type AppointmentCommonResponse struct {
	Message string `json:"message"`
}

func NewAppointmentCreateRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment schedule created, it will be processed soon",
	}
}

func NewAppointmentUpdateRequested() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment update requested, it will be processed soon",
	}
}

func NewAppointmentConfirmed() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment confirmed",
	}
}

func NewAppointmentCancelled() *AppointmentCommonResponse {
	return &AppointmentCommonResponse{
		Message: "appointment cancelled",
	}
}
