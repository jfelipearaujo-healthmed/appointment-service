package list_medical_reports

import (
	"strconv"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/medical_report_dto"
	list_medical_reports_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/medical_report/list_medical_reports"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase list_medical_reports_contract.UseCase
}

func NewHandler(useCase list_medical_reports_contract.UseCase) *handler {
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

	medicalReports, err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId), role)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, medical_report_dto.MapFromDomainSlice(medicalReports))
}
