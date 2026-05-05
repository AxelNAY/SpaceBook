package models

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CompanyID   uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	Company     Company   `gorm:"foreignKey:CompanyID;references:ID;constraint:OnDelete:CASCADE" json:"company,omitempty"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
