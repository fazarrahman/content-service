package entity

import "github.com/fazarrahman/content-service/domain"

type Image struct {
	domain.BaseModel
	FileName       string `json:"file_name" gorm:"type:varchar(255);not null"`
	Path           string `json:"path" gorm:"type:varchar(255);not null"`
	Group          string `json:"group" gorm:"type:varchar(255);not null"`
	ImageBase64Str string `json:"image_base64_str" gorm:"-"`
}
