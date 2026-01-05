package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user"`

	ResourceID uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	Resource   Resource  `gorm:"foreignKey:ResourceID;references:ID" json:"resource"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  string    `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
