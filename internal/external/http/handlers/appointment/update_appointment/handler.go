package update_appointment

import (
	"strconv"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	update_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/update_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/validator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase update_appointment_contract.UseCase
}

func NewHandler(useCase update_appointment_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(appointment_dto.UpdateAppointmentRequest)
	if err := c.Bind(req); err != nil {
		return http_response.BadRequest(c, "unable to parse the request body", err)
	}

	if err := validator.Validate(req); err != nil {
		return http_response.UnprocessableEntity(c, "invalid request body", err)
	}

	userId := c.Get("userId").(uint)
	appointmentId := c.Param("appointmentId")
	parsedAppointmentId, err := strconv.ParseUint(appointmentId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid appointment id", err)
	}

	if err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId), req); err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, appointment_dto.NewAppointmentUpdateRequested())
}
