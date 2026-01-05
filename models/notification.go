package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`

	Type    string `json:"type"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`

	CreatedAt time.Time `json:"created_at"`
}
