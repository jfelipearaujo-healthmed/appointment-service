package entities

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"
	"gorm.io/gorm"
)

type Event struct {
	ID uint `json:"id" gorm:"primaryKey"`

	ScheduleID uint            `json:"schedule_id"`
	PatientID  uint            `json:"patient_id"`
	DoctorID   uint            `json:"doctor_id"`
	DateTime   time.Time       `json:"date_time"`
	MessageID  string          `json:"message_id"`
	EventType  topic.EventType `json:"event_type"`
	Outcome    *string         `json:"outcome"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Appointment *Appointment `gorm:"foreignKey:EventID"`
}
