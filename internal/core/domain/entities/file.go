package entities

import "gorm.io/gorm"

type File struct {
	gorm.Model

	UserId       uint   `json:"user_id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
	Url          string `json:"url"`
}
