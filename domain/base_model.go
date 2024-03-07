package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy uuid.UUID      `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamptz;null"`
	UpdatedBy uuid.UUID      `json:"updated_by" gorm:"type:uuid;null"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz;null"`
	DeletedBy uuid.UUID      `json:"deleted_by" gorm:"type:uuid;null"`
}
