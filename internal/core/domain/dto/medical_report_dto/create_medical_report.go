package medical_report_dto

type CreateMedicalReportRequest struct {
	Comment string `json:"comment" validate:"required,min=5,max=255"`
}
