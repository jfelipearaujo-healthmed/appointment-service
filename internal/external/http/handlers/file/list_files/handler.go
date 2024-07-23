package list_files

import (
	file_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file"
	list_files_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/list_files"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase list_files_contract.UseCase
}

func New(useCase list_files_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)

	files, err := h.useCase.Execute(ctx, userId)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.OK(c, file_dto.MapFromDomainSlice(files))
}
