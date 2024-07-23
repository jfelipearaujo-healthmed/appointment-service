package get_appointment_files

import (
	"strconv"

	file_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file"
	get_appointment_files_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/get_appointment_files"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase get_appointment_files_contract.UseCase
}

func NewHandler(useCase get_appointment_files_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)

	appointmentId := c.Param("appointmentId")
	parsedAppointmentId, err := strconv.ParseUint(appointmentId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid appointment id", err)
	}

	files, err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId))
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, file_dto.MapFromDomainSlice(files))
}
