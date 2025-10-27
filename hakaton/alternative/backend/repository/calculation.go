package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateOrbitalCalculation(calc *models.OrbitalCalculation) error {
	if err := r.db.Create(calc).Error; err != nil {
		return fmt.Errorf("repository.CreateOrbitalCalculation: %w", err)
	}
	return nil
}

func (r *Repository) GetOrbitalCalculationByID(id uint64) (*models.OrbitalCalculation, error) {
	var calc models.OrbitalCalculation
	if err := r.db.First(&calc, id).Error; err != nil {
		return nil, fmt.Errorf("repository.GetOrbitalCalculationByID: %w", err)
	}
	return &calc, nil
}

func (r *Repository) UpdateOrbitalCalculation(calc *models.OrbitalCalculation) error {
	if err := r.db.Save(calc).Error; err != nil {
		return fmt.Errorf("repository.UpdateOrbitalCalculation: %w", err)
	}
	return nil
}

func (r *Repository) DeleteOrbitalCalculation(id uint64) error {
	if err := r.db.Delete(&models.OrbitalCalculation{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteOrbitalCalculation: %w", err)
	}
	return nil
}

func (r *Repository) ListCalculationsByCometID(cometID uint64) ([]models.OrbitalCalculation, error) {
	var calcs []models.OrbitalCalculation
	if err := r.db.Where("comet_id = ?", cometID).Order("calculated_at DESC").Find(&calcs).Error; err != nil {
		return nil, fmt.Errorf("repository.ListCalculationsByCometID: %w", err)
	}
	return calcs, nil
}