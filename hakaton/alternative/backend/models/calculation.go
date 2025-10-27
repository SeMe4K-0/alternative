package models

import "time"

type OrbitalCalculation struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CometID          uint64    `json:"comet_id" gorm:"index"`
	CalculatedAt     time.Time `json:"calculated_at"`
	SemiMajorAxis    float64   `json:"semi_major_axis"`
	Eccentricity     float64   `json:"eccentricity"`
	Inclination      float64   `json:"inclination"`
	LonAscendingNode float64   `json:"lon_ascending_node"`
	ArgPeriapsis     float64   `json:"arg_periapsis"`
	TimePerihelion   time.Time `json:"time_perihelion"`
	IsLatest         bool      `json:"is_latest" gorm:"index"`

	ApproachDate time.Time `json:"approach_date"`
	DistanceAU   float64   `json:"distance_au"`
	DistanceKM   float64   `json:"distance_km"`

	Observations []*Observation `json:"observations,omitempty" gorm:"many2many:calculation_observations;"`
}
