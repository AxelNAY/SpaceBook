package models

import "github.com/google/uuid"

type ReservationResource struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ResourceID    uuid.UUID `gorm:"type:uuid;not null" json:"resource_id"`
	Resource      Resource  `gorm:"foreignKey:ResourceID;references:ID;constraint:OnDelete:CASCADE" json:"resource,omitempty"`
	ReservationID uuid.UUID `gorm:"type:uuid;not null" json:"reservation_id"`
	Status        string    `gorm:"default:'pending'" json:"status"`
}
