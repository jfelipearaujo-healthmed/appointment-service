package appointment_dto

type AppointmentRequested struct {
	Message string `json:"message"`
}

func NewCreateAppointmentRequested() *AppointmentRequested {
	return &AppointmentRequested{
		Message: "appointment schedule created, access your appointment later to see the status",
	}
}

func NewUpdateAppointmentRequested() *AppointmentRequested {
	return &AppointmentRequested{
		Message: "appointment update requested, access your appointment later to see the status",
	}
}
