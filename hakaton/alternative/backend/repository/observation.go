package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateObservation(observation *models.Observation) error {
	if err := r.db.Create(observation).Error; err != nil {
		return fmt.Errorf("repository.CreateObservation: %w", err)
	}
	return nil
}

func (r *Repository) GetObservationByID(id uint64) (*models.Observation, error) {
	var observation models.Observation
	if err := r.db.First(&observation, id).Error; err != nil {
		return nil, fmt.Errorf("repository.GetObservationByID: %w", err)
	}
	return &observation, nil
}

func (r *Repository) UpdateObservation(observation *models.Observation) error {
	if err := r.db.Save(observation).Error; err != nil {
		return fmt.Errorf("repository.UpdateObservation: %w", err)
	}
	return nil
}

func (r *Repository) DeleteObservation(id uint64) error {
	if err := r.db.Delete(&models.Observation{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteObservation: %w", err)
	}
	return nil
}

func (r *Repository) ListObservationsByCometID(cometID uint64) ([]models.Observation, error) {
	var observations []models.Observation
	if err := r.db.Where("comet_id = ?", cometID).Order("observed_at DESC").Find(&observations).Error; err != nil {
		return nil, fmt.Errorf("repository.ListObservationsByCometID: %w", err)
	}
	return observations, nil
}