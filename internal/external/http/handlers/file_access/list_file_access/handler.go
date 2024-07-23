package list_file_access

import (
	file_access_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file_access"
	list_file_access_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file_access/list_file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase list_file_access_contract.UseCase
}

func NewHandler(useCase list_file_access_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)

	fileAccess, err := h.useCase.Execute(ctx, userId)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.Created(c, file_access_dto.MapFromDomainSlice(fileAccess))
}
