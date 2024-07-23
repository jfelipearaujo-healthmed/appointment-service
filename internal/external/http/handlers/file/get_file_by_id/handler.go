package get_file_by_id

import (
	"strconv"

	file_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file"
	get_file_by_id_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/get_file_by_id"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase get_file_by_id_contract.UseCase
}

func NewHandler(useCase get_file_by_id_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)
	fileId := c.Param("fileId")
	parsedFileId, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid file id", err)
	}

	file, err := h.useCase.Execute(ctx, userId, uint(parsedFileId))
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.Created(c, file_dto.MapFromDomain(file))
}
