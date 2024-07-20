package entities

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model

	ScheduleID      uint       `json:"schedule_id" gorm:"idx_schedule_patient_doctor_date_time"`
	PatientID       uint       `json:"patient_id" gorm:"idx_schedule_patient_doctor_date_time"`
	DoctorID        uint       `json:"doctor_id" gorm:"idx_schedule_patient_doctor_date_time"`
	DateTime        time.Time  `json:"date_time" gorm:"idx_schedule_patient_doctor_date_time"`
	StartedAt       *time.Time `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	ConfirmedAt     *time.Time `json:"confirmed_at"`
	CancelledBy     *uint      `json:"cancelled_by"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	CancelledReason *string    `json:"cancelled_reason"`
}
