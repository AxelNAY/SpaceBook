package models

import "github.com/google/uuid"

type PlaceUser struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PlaceID uuid.UUID `gorm:"type:uuid;not null" json:"place_id"`
	Place   Place     `gorm:"foreignKey:PlaceID;references:ID;constraint:OnDelete:CASCADE" json:"place,omitempty"`
	UserID  uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User    User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}
