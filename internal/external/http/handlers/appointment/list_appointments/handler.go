package list_appointments

import (
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/appointment_dto"
	list_appointments_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/appointment/list_appointments"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase list_appointments_contract.UseCase
}

func NewHandler(useCase list_appointments_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)
	roleName := c.Get("role").(string)
	role := role.GetRoleByName(roleName)

	appointment, err := h.useCase.Execute(ctx, userId, role)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, appointment_dto.MapFromDomainSlice(appointment))
}
