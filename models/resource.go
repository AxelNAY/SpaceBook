package models

import "time"

type Resource struct {
	ID          string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Category    string    `gorm:"default:none"`
	Capacity    int
	Status      string    `gorm:"default:available"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
