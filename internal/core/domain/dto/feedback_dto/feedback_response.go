package feedback_dto

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type FeedbackResponse struct {
	ID uint `json:"id"`

	AppointmentID uint    `json:"appointment_id"`
	Rating        float64 `json:"rating"`
	Comment       string  `json:"comment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapFromDomain(feedback *entities.Feedback) *FeedbackResponse {
	return &FeedbackResponse{
		ID: feedback.ID,

		AppointmentID: feedback.AppointmentID,
		Rating:        feedback.Rating,
		Comment:       feedback.Comment,
		CreatedAt:     feedback.CreatedAt,
		UpdatedAt:     feedback.UpdatedAt,
	}
}

func MapFromDomainSlice(feedbacks []entities.Feedback) []*FeedbackResponse {
	mapped := make([]*FeedbackResponse, len(feedbacks))

	for i := range feedbacks {
		mapped[i] = MapFromDomain(&feedbacks[i])
	}

	return mapped
}
