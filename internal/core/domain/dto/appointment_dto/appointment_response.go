package appointment_dto

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type AppointmentResponse struct {
	ID uint `json:"id"`

	ScheduleID      uint       `json:"schedule_id"`
	PatientID       uint       `json:"patient_id"`
	DoctorID        uint       `json:"doctor_id"`
	DateTime        time.Time  `json:"date_time"`
	StartedAt       *time.Time `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	CancelledBy     *uint      `json:"cancelled_by"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	CancelledReason *string    `json:"cancelled_reason"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapFromDomain(appointment *entities.Appointment) *AppointmentResponse {
	return &AppointmentResponse{
		ID:              appointment.ID,
		ScheduleID:      appointment.ScheduleID,
		PatientID:       appointment.PatientID,
		DoctorID:        appointment.DoctorID,
		DateTime:        appointment.DateTime,
		StartedAt:       appointment.StartedAt,
		EndedAt:         appointment.EndedAt,
		ConfirmedAt:     appointment.ConfirmedAt,
		CancelledBy:     appointment.CancelledBy,
		CancelledAt:     appointment.CancelledAt,
		CancelledReason: appointment.CancelledReason,
		CreatedAt:       appointment.CreatedAt,
		UpdatedAt:       appointment.UpdatedAt,
	}
}
