package feedback_dto

type CreateFeedbackRequest struct {
	Rating  float64 `json:"rating" validate:"required,gte=0,lte=5"`
	Comment string  `json:"comment" validate:"required,min=5,max=255"`
}
