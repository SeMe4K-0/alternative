package repository

import (
	"backend/models"
	"fmt"
)

func (r *Repository) CreateComet(comet *models.Comet) error {
	if err := r.db.Create(comet).Error; err != nil {
		return fmt.Errorf("repository.CreateComet: %w", err)
	}
	return nil
}

func (r *Repository) GetCometByID(id uint64) (*models.Comet, error) {
	var comet models.Comet
	if err := r.db.First(&comet, id).Error; err != nil {
		return nil, fmt.Errorf("repository.GetCometByID: %w", err)
	}
	return &comet, nil
}

func (r *Repository) UpdateComet(comet *models.Comet) error {
	if err := r.db.Save(comet).Error; err != nil {
		return fmt.Errorf("repository.UpdateComet: %w", err)
	}
	return nil
}

func (r *Repository) DeleteComet(id uint64) error {
	if err := r.db.Delete(&models.Comet{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteComet: %w", err)
	}
	return nil
}

func (r *Repository) ListCometsByUserID(userID uint64) ([]models.Comet, error) {
	var comets []models.Comet
	if err := r.db.Where("user_id = ?", userID).Find(&comets).Error; err != nil {
		return nil, fmt.Errorf("repository.ListCometsByUserID: %w", err)
	}
	return comets, nil
}