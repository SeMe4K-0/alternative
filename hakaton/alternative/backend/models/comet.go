package models

import (
	"time"
)


type Comet struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID      uint64 `json:"user_id" gorm:"index"`
	Name        string `json:"name" gorm:"size:255"`
	Description string `json:"description"`

	Observations        []Observation        `json:"observations,omitempty"`
	OrbitalCalculations []OrbitalCalculation `json:"orbital_calculations,omitempty"`
}
