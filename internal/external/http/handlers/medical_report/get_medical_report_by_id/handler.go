package get_medical_report_by_id

import (
	"strconv"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/medical_report_dto"
	get_medical_report_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/get_medical_report_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase get_medical_report_by_id_contract.UseCase
}

func NewHandler(useCase get_medical_report_by_id_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)
	roleName := c.Get("role").(string)
	role := role.GetRoleByName(roleName)

	appointmentId := c.Param("appointmentId")
	parsedAppointmentId, err := strconv.ParseUint(appointmentId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid appointment id", err)
	}

	medicalReportId := c.Param("medicalReportId")
	parsedMedicalReportId, err := strconv.ParseUint(medicalReportId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid feedback id", err)
	}

	medicalReport, err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId), uint(parsedMedicalReportId), role)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, medical_report_dto.MapFromDomain(medicalReport))
}
