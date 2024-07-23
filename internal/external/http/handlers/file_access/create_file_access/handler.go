package create_file_access

import (
	"strconv"

	file_access_dto "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/dto/file_access"
	create_file_access_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file_access/create_file_access"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/validator"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase create_file_access_contract.UseCase
}

func NewHandler(useCase create_file_access_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)

	req := new(file_access_dto.CreateFileAccess)
	if err := c.Bind(req); err != nil {
		return http_response.BadRequest(c, "unable to parse the request body", err)
	}

	fileId := c.Param("fileId")
	parsedFileId, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		return http_response.BadRequest(c, "invalid file id", err)
	}

	req.FileID = uint(parsedFileId)

	if err := validator.Validate(req); err != nil {
		return http_response.UnprocessableEntity(c, "invalid request body", err)
	}

	fileAccess, err := h.useCase.Execute(ctx, userId, req)
	if err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.Created(c, file_access_dto.MapFromDomain(fileAccess))
}
