package events

import "github.com/jfelipearaujo-healthmed/appointment-service/internal/external/topic"

const (
	CreateAppointment topic.EventType = "create_appointment"
	UpdateAppointment topic.EventType = "update_appointment"

	CreateFeedback topic.EventType = "create_feedback"
)
