package upload_file

import (
	"fmt"

	upload_file_contract "github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/use_cases/file/upload_file"
	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/infrastructure/shared/http_response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	useCase upload_file_contract.UseCase
}

func NewHandler(useCase upload_file_contract.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userId").(uint)
	file, err := c.FormFile("file")
	if err != nil {
		return http_response.BadRequest(c, "unable to get file", err)
	}

	mimeType := file.Header.Get("Content-Type")

	if mimeType != "application/pdf" {
		return http_response.BadRequest(c, "invalid file type, only PDF files are allowed", nil)
	}

	fileData, err := file.Open()
	if err != nil {
		return http_response.BadRequest(c, "unable to open file", err)
	}
	defer fileData.Close()

	if err := h.useCase.Execute(ctx, userId, file.Filename, mimeType, file.Size, fileData); err != nil {
		return http_response.HandleErr(c, err)
	}

	return http_response.Created(c, map[string]interface{}{
		"status": fmt.Sprintf("file '%s' uploaded", file.Filename),
	})
}
