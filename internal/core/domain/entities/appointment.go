package entities

import (
	"time"

	"gorm.io/gorm"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	ScheduleInAnalysis   Status = "schedule_in_analysis"
	ReScheduleInAnalysis Status = "re_schedule_in_analysis"
	Confirmed            Status = "confirmed"
	InProgress           Status = "in_progress"
	Concluded            Status = "concluded"
	Cancelled            Status = "cancelled"
)

type Appointment struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ScheduleID      uint       `json:"schedule_id"`
	PatientID       uint       `json:"patient_id"`
	DoctorID        uint       `json:"doctor_id"`
	DateTime        time.Time  `json:"date_time"`
	Status          Status     `json:"status"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	ConfirmedAt     *time.Time `json:"confirmed_at,omitempty"`
	CancelledBy     *uint      `json:"cancelled_by,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CancelledReason *string    `json:"cancelled_reason,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	EventID uint `json:"event_id"`
}
