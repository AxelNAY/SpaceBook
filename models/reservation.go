package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	RessourceID uuid.UUID `gorm:"type:uuid;not null;column:ressource_id" json:"ressource_id"`
	Ressource   Resource  `gorm:"foreignKey:RessourceID;references:ID;constraint:OnDelete:CASCADE" json:"ressource,omitempty"`

	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Status        string    `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
