package feedback_dto

type FeedbackCommonResponse struct {
	Message string `json:"message"`
}

func NewFeedbackSent() *FeedbackCommonResponse {
	return &FeedbackCommonResponse{
		Message: "your feedback has been sent, please wait for the feedback to be processed",
	}
}
