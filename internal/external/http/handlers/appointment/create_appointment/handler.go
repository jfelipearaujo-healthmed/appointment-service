package create_appointment

import (
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	create_appointment_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/create_appointment"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/validator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase create_appointment_contract.UseCase
}

func NewHandler(useCase create_appointment_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(appointment_dto.CreateAppointmentRequest)
	if err := c.Bind(req); err != nil {
		return http_response.BadRequest(c, "unable to parse the request body", err)
	}

	if err := validator.Validate(req); err != nil {
		return http_response.UnprocessableEntity(c, "invalid request body", err)
	}

	userId := c.Get("userId").(uint)

	if _, err := h.useCase.Execute(ctx, userId, req); err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, appointment_dto.NewAppointmentCreateRequested())
}
