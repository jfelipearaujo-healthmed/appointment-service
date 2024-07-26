package get_appointment_by_id

import (
	"log/slog"
	"strconv"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	get_appointment_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase get_appointment_by_id_contract.UseCase
}

func NewHandler(useCase get_appointment_by_id_contract.UseCase) *handler {
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

	slog.InfoContext(ctx, "getting appointment by id", "userId", userId, "role", roleName, "appointmentId", appointmentId)

	appointment, err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId), role)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	slog.InfoContext(ctx, "appointment by id retrieved", "appointment", appointment)

	return http_response.OK(c, appointment_dto.MapFromDomain(appointment))
}
