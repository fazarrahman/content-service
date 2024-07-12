package entity

import (
	"time"

	"github.com/fazarrahman/content-service/domain"
	"github.com/google/uuid"
)

type Image struct {
	domain.BaseModel
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	FileName    string    `json:"file_name" gorm:"type:varchar(255);not null"`
	Path        string    `json:"path" gorm:"type:varchar(255);not null"`
	Group       string    `json:"group" gorm:"type:varchar(255);not null"`
	File        string    `json:"file" gorm:"-"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text;not null"`
	TagStr      []string  `json:"tags" gorm:"-"`
	Tags        []Tag     `json:"-" gorm:"many2many:image_tags"`
}

type Tag struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	Images    []Image   `gorm:"many2many:image_tags"`
}
