package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PlaceID   *uuid.UUID `gorm:"type:uuid" json:"place_id,omitempty"`
	Email     string     `gorm:"unique" json:"email"`
	Username  string     `json:"username"`
	Phone     int        `json:"phone,omitempty"`
	Password  []byte     `json:"-"`
	Role      string     `json:"role"`
	Status	  string	 `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
