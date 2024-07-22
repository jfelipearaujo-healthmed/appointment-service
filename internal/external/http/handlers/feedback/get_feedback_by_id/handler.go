package get_feedback_by_id

import (
	"strconv"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/feedback_dto"
	get_feedback_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/feedback/get_feedback_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/external/http/middlewares/role"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase get_feedback_by_id_contract.UseCase
}

func NewHandler(useCase get_feedback_by_id_contract.UseCase) *handler {
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

	feedbackId := c.Param("feedbackId")
	parsedFeedbackId, err := strconv.ParseUint(feedbackId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid feedback id", err)
	}

	feedback, err := h.useCase.Execute(ctx, userId, uint(parsedAppointmentId), uint(parsedFeedbackId), role)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, feedback_dto.MapFromDomain(feedback))
}
