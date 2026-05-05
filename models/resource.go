package models

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PlaceID    *uuid.UUID `gorm:"type:uuid" json:"place_id,omitempty"`
	Place      *Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnDelete:SET NULL" json:"place,omitempty"`
	CategoryID *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	Category   *Category  `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name"`
	Status     string     `gorm:"default:available" json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
