package usecase

import (
	"backend/models"
	"backend/named_errors"
	"backend/repository"
	"backend/store"
	"errors"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"
)

type ObservationUsecase struct {
	repo       *repository.Repository
	minioStore *store.MinIOStore
}

func NewObservationUsecase(repo *repository.Repository, minioStore *store.MinIOStore) *ObservationUsecase {
	return &ObservationUsecase{repo: repo, minioStore: minioStore}
}

func (u *ObservationUsecase) checkCometOwnership(userID, cometID uint64) error {
	comet, err := u.repo.GetCometByID(cometID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return named_errors.ErrNotFound
		}
		return err
	}

	if comet.UserID != userID {
		return named_errors.ErrAccessDenied
	}

	return nil
}

func (u *ObservationUsecase) CreateObservation(userID, cometID uint64, observedAt time.Time, ra, dec float64, notes string) (*models.Observation, error) {
	if err := u.checkCometOwnership(userID, cometID); err != nil {
		return nil, fmt.Errorf("usecase.CreateObservation: ownership check failed: %w", err)
	}

	observation := &models.Observation{
		CometID:    cometID,
		ObservedAt: observedAt,
		RA:         ra,
		Dec:        dec,
		Notes:      notes,
	}

	if err := u.repo.CreateObservation(observation); err != nil {
		return nil, fmt.Errorf("usecase.CreateObservation: %w", err)
	}

	return observation, nil
}

func (u *ObservationUsecase) GetObservationByID(userID, observationID uint64) (*models.Observation, error) {
	observation, err := u.repo.GetObservationByID(observationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, named_errors.ErrNotFound
		}
		return nil, fmt.Errorf("usecase.GetObservationByID: %w", err)
	}

	if err := u.checkCometOwnership(userID, observation.CometID); err != nil {
		return nil, fmt.Errorf("usecase.GetObservationByID: ownership check failed: %w", err)
	}

	return observation, nil
}

func (u *ObservationUsecase) ListObservationsByComet(userID, cometID uint64) ([]models.Observation, error) {
	if err := u.checkCometOwnership(userID, cometID); err != nil {
		return nil, fmt.Errorf("usecase.ListObservationsByComet: ownership check failed: %w", err)
	}

	observations, err := u.repo.ListObservationsByCometID(cometID)
	if err != nil {
		return nil, fmt.Errorf("usecase.ListObservationsByComet: %w", err)
	}

	return observations, nil
}

func (u *ObservationUsecase) UpdateObservation(userID, observationID uint64, newObsData models.Observation) (*models.Observation, error) {
	existingObs, err := u.GetObservationByID(userID, observationID)
	if err != nil {
		return nil, fmt.Errorf("usecase.UpdateObservation: %w", err)
	}

	existingObs.ObservedAt = newObsData.ObservedAt
	existingObs.RA = newObsData.RA
	existingObs.Dec = newObsData.Dec
	existingObs.Notes = newObsData.Notes
	existingObs.ImageURL = newObsData.ImageURL

	if err := u.repo.UpdateObservation(existingObs); err != nil {
		return nil, fmt.Errorf("usecase.UpdateObservation: %w", err)
	}

	return existingObs, nil
}

func (u *ObservationUsecase) DeleteObservation(userID, observationID uint64) error {
	if _, err := u.GetObservationByID(userID, observationID); err != nil {
		return fmt.Errorf("usecase.DeleteObservation: %w", err)
	}

	if err := u.repo.DeleteObservation(observationID); err != nil {
		return fmt.Errorf("usecase.DeleteObservation: %w", err)
	}

	return nil
}

func (u *ObservationUsecase) UploadObservationImage(userID, observationID uint64, file io.Reader, filename, contentType string) (*models.Observation, error) {
	observation, err := u.GetObservationByID(userID, observationID)
	if err != nil {
		return nil, fmt.Errorf("usecase.UploadObservationImage: %w", err)
	}

	if !u.minioStore.ValidateImageType(contentType) {
		return nil, named_errors.ErrInvalidInput
	}

	if observation.ImageURL != "" {
		oldObjectName := observation.ImageURL[len("/api/files/"):]
		if err := u.minioStore.DeleteObservationImage(oldObjectName); err != nil {
			fmt.Printf("Warning: failed to delete old image from MinIO: %v\n", err)
		}
	}

	objectName, err := u.minioStore.UploadObservationImage(observationID, file, filename, contentType)
	if err != nil {
		return nil, fmt.Errorf("usecase.UploadObservationImage: failed to upload image: %w", err)
	}

	imageURL := u.minioStore.GetObservationImageURL(objectName)
	observation.ImageURL = imageURL

	if err := u.repo.UpdateObservation(observation); err != nil {
		u.minioStore.DeleteObservationImage(objectName)
		return nil, fmt.Errorf("usecase.UploadObservationImage: failed to update observation: %w", err)
	}

	return observation, nil
}

func (u *ObservationUsecase) DeleteObservationImage(userID, observationID uint64) (*models.Observation, error) {
	observation, err := u.GetObservationByID(userID, observationID)
	if err != nil {
		return nil, fmt.Errorf("usecase.DeleteObservationImage: %w", err)
	}

	if observation.ImageURL != "" {
		objectName := observation.ImageURL[len("/api/files/"):]
		if err := u.minioStore.DeleteObservationImage(objectName); err != nil {
			fmt.Printf("Warning: failed to delete image from MinIO: %v\n", err)
		}
	}

	observation.ImageURL = ""
	if err := u.repo.UpdateObservation(observation); err != nil {
		return nil, fmt.Errorf("usecase.DeleteObservationImage: failed to update observation: %w", err)
	}

	return observation, nil
}
