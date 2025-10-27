package usecase

import (
	"backend/models"
	"backend/named_errors"
	"backend/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CometUsecase struct {
	cometRepo *repository.Repository
}

func NewCometUsecase(cometRepo *repository.Repository) *CometUsecase {
	return &CometUsecase{cometRepo: cometRepo}
}

func (u *CometUsecase) CreateComet(userID uint64, name, description string) (*models.Comet, error) {
	if name == "" {
		return nil, named_errors.ErrInvalidInput
	}

	comet := &models.Comet{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := u.cometRepo.CreateComet(comet); err != nil {
		return nil, fmt.Errorf("usecase.CreateComet: %w", err)
	}

	return comet, nil
}

func (u *CometUsecase) GetCometByID(userID, cometID uint64) (*models.Comet, error) {
	comet, err := u.cometRepo.GetCometByID(cometID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, named_errors.ErrNotFound
		}
		return nil, fmt.Errorf("usecase.GetCometByID: %w", err)
	}

	if comet.UserID != userID {
		return nil, named_errors.ErrAccessDenied
	}

	return comet, nil
}

func (u *CometUsecase) ListCometsByUserID(userID uint64) ([]models.Comet, error) {
	comets, err := u.cometRepo.ListCometsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListCometsByUserID: %w", err)
	}
	return comets, nil
}

func (u *CometUsecase) UpdateComet(userID, cometID uint64, name, description string) (*models.Comet, error) {
	comet, err := u.GetCometByID(userID, cometID)
	if err != nil {
		return nil, fmt.Errorf("usecase.UpdateComet: %w", err)
	}

	if name != "" {
		comet.Name = name
	}
	comet.Description = description

	if err := u.cometRepo.UpdateComet(comet); err != nil {
		return nil, fmt.Errorf("usecase.UpdateComet: %w", err)
	}

	return comet, nil
}

func (u *CometUsecase) DeleteComet(userID, cometID uint64) error {
	if _, err := u.GetCometByID(userID, cometID); err != nil {
		return fmt.Errorf("usecase.DeleteComet: %w", err)
	}

	if err := u.cometRepo.DeleteComet(cometID); err != nil {
		return fmt.Errorf("usecase.DeleteComet: %w", err)
	}
	return nil
}