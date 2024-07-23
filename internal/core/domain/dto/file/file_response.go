package file_dto

import (
	"time"

	"github.com/jfelipearaujo-healthmed/appointment-service/internal/core/domain/entities"
)

type FileResponse struct {
	ID uint `json:"id"`

	FileName string `json:"file_name"`
	Path     string `json:"path"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapFromDomain(file *entities.File) *FileResponse {
	return &FileResponse{
		ID: file.ID,

		FileName: file.OriginalName,
		Path:     file.Url,

		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

func MapFromDomainSlice(files []entities.File) []*FileResponse {
	mapped := make([]*FileResponse, len(files))

	for i := range files {
		mapped[i] = MapFromDomain(&files[i])
	}

	return mapped
}
