package models

import "time"

type Observation struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CometID    uint64    `json:"comet_id" gorm:"index"`
	ObservedAt time.Time `json:"observed_at" gorm:"index"`
	RA         float64   `json:"right_ascension"`
	Dec        float64   `json:"declination"`
	ImageURL   string    `json:"image_url,omitempty"`
	Notes      string    `json:"notes,omitempty"`
}
